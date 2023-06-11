package handlers

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/savsgio/gotils/uuid"
	"os"
	"pahomov_frolovsky_cson/handlers/response"
	"pahomov_frolovsky_cson/postgres"
	"pahomov_frolovsky_cson/utilities"
	"time"
)

type News struct {
	Id          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	FileUrl     string `json:"fileUrl"`
	FileName    string `json:"file_name"`
	CreateAt    string `json:"createAt"`
	UpdateAt    string `json:"updateAt"`
}

func HandlerAddNewsRequest(ctx *fiber.Ctx) error {
	// Request
	var request News
	//var request struct {
	//	Id          int       `json:"id,omitempty"`
	//	Title       string    `json:"title"`
	//	Description string    `json:"description"`
	//	FileUrl     string    `json:"fileUrl"`
	//	CreateAt    time.Time `json:"createAt,omitempty"`
	//}
	var uploadPath string
	utilities.LookupEnv(&uploadPath, "UPLOAD_PATH", "")
	fmt.Printf("\nHandlerAddNewsRequest %s", ctx.Request().Body())
	err := ctx.BodyParser(&request)
	if err != nil {
		return response.ErrBadRequest.Send(ctx)
	}

	form, _ := ctx.MultipartForm()

	//Loop through files:
	var fileId int
	for _, file := range form.File["File"] {
		fmt.Printf("\nFile: %s %d %s", file.Filename, file.Size, file.Header["Content-Type"][0])

		fileUid := uuid.V4()
		fileUrl := fmt.Sprintf("%s/%s", uploadPath, fileUid)

		// Save the files to disk:
		if err = ctx.SaveFile(file, fmt.Sprintf("./%s", fileUrl)); err != nil {
			return response.ErrInternal.AddMessage(err).Send(ctx)
		}

		err = postgres.Conn.QueryRow(ctx.Context(),
			"insert into public.files (file_path, file_name, file_size, file_uid) values ($1, $2, $3, $4) returning id ",
			fileUrl, file.Filename, file.Size, fileUid).Scan(&fileId)
		if err != nil {
			return response.OK.SetData(request).Send(ctx)
		}
		request.FileUrl = fileUrl
		request.FileName = file.Filename
		break
	}

	// Insert postgres
	createAt := time.Now()
	request.CreateAt = createAt.Format(time.DateTime)
	err = postgres.Conn.QueryRow(ctx.Context(), `
		INSERT INTO public.news 
    	(title, description,file_url, create_at,update_at, file_name, file_id)
		VALUES ($1, $2, $3, $4, $4, $5, $6)
		RETURNING id`,
		request.Title, request.Description, request.FileUrl, request.CreateAt, request.FileName, fileId,
	).Scan(&request.Id)
	if err != nil {
		return response.ErrInternal.AddMessage(err).Send(ctx)
	}

	return response.OK.SetData(request).Send(ctx)
}

func HandlerUpdateNews(ctx *fiber.Ctx) error {
	// Request
	var request struct {
		Id          int
		Title       string    `json:"title"`
		Description string    `json:"description"`
		FileUrl     string    `json:"fileUrl"`
		CreateAt    time.Time `json:"createAt"`
		UpdateAt    time.Time `json:"updateAt"`
	}
	var uploadPath string
	utilities.LookupEnv(&uploadPath, "UPLOAD_PATH", "")

	err := ctx.BodyParser(&request)
	if err != nil {
		return response.ErrBadRequest.Send(ctx)
	}

	request.Id, err = ctx.ParamsInt("id")
	if err != nil {
		return response.ErrBadRequest.AddMessage(err).Send(ctx)
	}

	form, _ := ctx.MultipartForm()

	//Loop through files:
	for _, file := range form.File["File"] {
		fileUrl := fmt.Sprintf("%s%s", uploadPath, file.Filename)
		fmt.Println(file.Filename, file.Size, file.Header["Content-Type"][0])
		// => "tutorial.pdf" 360641 "application/pdf"

		// Save the files to disk:
		if err := ctx.SaveFile(file, fmt.Sprintf("./%s", fileUrl)); err != nil {
			return response.ErrInternal.AddMessage(err).Send(ctx)
		}
		request.FileUrl = fileUrl
		break
	}

	request.UpdateAt = time.Now()
	_, err = postgres.Conn.Exec(ctx.Context(), "UPDATE news SET title=$1, description=$2, update_at=$3 WHERE id=$4 ",
		request.Title, request.Description, request.UpdateAt, request.Id)
	if err != nil {
		return response.ErrInternal.AddMessage(err).Send(ctx)
	}

	if request.FileUrl != "" {
		_, err = postgres.Conn.Exec(ctx.Context(), "UPDATE news SET file_url=$1 WHERE id=$2",
			request.FileUrl, request.Id)
		if err != nil {
			return response.ErrInternal.AddMessage(err).Send(ctx)
		}
	}

	return response.OK.SetData(request).Send(ctx)
}

func HandlerGetNewsSlice(ctx *fiber.Ctx) error {
	var uploadPath string
	utilities.LookupEnv(&uploadPath, "UPLOAD_PATH", "")

	rows, err := postgres.Conn.Query(ctx.Context(), "select id, title, description, file_url, create_at, update_at from news where not is_delete")
	if err != nil {
		return response.ErrInternal.AddMessage(err).Send(ctx)
	}

	var resp []News
	for rows.Next() {
		var item News
		var createAt, UpdateAt *time.Time
		err = rows.Scan(&item.Id, &item.Title, &item.Description, &item.FileUrl, &createAt, &UpdateAt)
		if createAt != nil {
			item.CreateAt = createAt.Format("2006-01-02 15:04")
		}
		if UpdateAt != nil {
			item.UpdateAt = UpdateAt.Format("2006-01-02 15:04")
		}
		if err != nil {
			return response.ErrInternal.AddMessage(err).Send(ctx)
		}
		resp = append(resp, item)
	}
	return response.OK.SetData(resp).Send(ctx)
}

func HandlerGetNews(ctx *fiber.Ctx) error {
	var uploadPath string
	utilities.LookupEnv(&uploadPath, "UPLOAD_PATH", "")
	newsId, err := ctx.ParamsInt("id")
	if err != nil {
		return response.ErrBadRequest.AddMessage(err).Send(ctx)
	}

	var item News
	var createAt, updateAt *time.Time
	err = postgres.Conn.QueryRow(ctx.Context(),
		"select id, title, description, file_url, create_at, update_at from news where id = $1 and not is_delete", newsId).
		Scan(&item.Id, &item.Title, &item.Description, &item.FileUrl, &createAt, &updateAt)
	if createAt != nil {
		item.CreateAt = createAt.Format("2006-01-02 15:04")
	}
	if updateAt != nil {
		item.UpdateAt = updateAt.Format("2006-01-02 15:04")
	}
	if err != nil {
		return response.ErrInternal.AddMessage(err).Send(ctx)
	}

	return response.OK.SetData(item).Send(ctx)
}

func HandlerDeleteNews(ctx *fiber.Ctx) error {
	newsId, err := ctx.ParamsInt("id")
	if err != nil {
		return response.ErrBadRequest.AddMessage(err).Send(ctx)
	}

	_, err = postgres.Conn.Exec(ctx.Context(),
		"update public.news set is_delete = true where id = $1", newsId)
	if err != nil {
		return response.ErrInternal.AddMessage(err).Send(ctx)
	}

	return response.OK.Send(ctx)
}

func HandlerGetFile(c *fiber.Ctx) error {
	fileUid := c.Params("fileuid")

	var filePath, fileName string
	err := postgres.Conn.QueryRow(c.Context(),
		"SELECT file_path, file_name FROM public.files WHERE file_uid = $1",
		fileUid).Scan(&filePath, &fileName)
	if err != nil {
		return response.ErrInternal.AddMessage(err).Send(c)
	}

	// Проверяем существование файла
	_, err = os.Stat(filePath)
	if os.IsNotExist(err) {
		return response.ErrNotFound.AddMessage("Файл не найден").Send(c)
	}

	c.Set("Content-Disposition", "attachment; filename="+fileName)
	return c.SendFile(filePath)
}
