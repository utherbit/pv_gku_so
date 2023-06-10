package handlers

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"pahomov_frolovsky_cson/handlers/response"
	"pahomov_frolovsky_cson/postgres"
	"pahomov_frolovsky_cson/utilities"
	"time"
)

func HandlerAddNewsRequest(ctx *fiber.Ctx) error {
	// Request
	var request struct {
		Id          int
		Title       string    `json:"title"`
		Description string    `json:"description"`
		FileUrl     string    `json:"fileUrl"`
		CreateAt    time.Time `json:"createAt"`
	}
	var uploadPath string
	utilities.LookupEnv(&uploadPath, "UPLOAD_PATH", "")

	err := ctx.BodyParser(&request)
	if err != nil {
		return response.ErrBadRequest.Send(ctx)
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

	// Insert postgres
	request.CreateAt = time.Now()
	err = postgres.Conn.QueryRow(ctx.Context(), `
		INSERT INTO public.news 
    	(title, description,file_url, create_at,update_at )
		VALUES ($1, $2, $3, $4)
		RETURNING id`,
		request.Title, request.Description, request.FileUrl, request.CreateAt, request.CreateAt,
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
	_, err = postgres.Conn.Exec(ctx.Context(), "UPDATE news SET title=$1, description=$2, update_at=$3 WHERE id=$4",
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
