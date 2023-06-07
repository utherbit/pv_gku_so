package handlers

import (
	"github.com/gofiber/fiber/v2"
	"pahomov_frolovsky_cson/handlers/response"
	"pahomov_frolovsky_cson/postgres"
	"time"
)

type Service struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	InfoHref    string    `json:"info_href"`
	DateUpdate  time.Time `json:"date_update,omitempty"`
	Actual      bool      `json:"actual,omitempty"`
}

func HandlerGetServices(ctx *fiber.Ctx) error {
	// get from postgres
	rows, err := postgres.Conn.Query(ctx.Context(), "select id, title, description, info_href, date_update from public.services where actual")
	if err != nil {
		return response.ErrInternal.AddMessage(err).Send(ctx)
	}

	var services []Service
	for rows.Next() {
		service := Service{}
		err = rows.Scan(
			&service.ID,
			&service.Title,
			&service.Description,
			&service.InfoHref,
			&service.DateUpdate,
		)
		if err != nil {
			return response.ErrInternal.AddMessage(err).Send(ctx)
		}
		services = append(services, service)
	}

	return response.OK.SetData(services).Send(ctx)
}
func HandlerGetService(ctx *fiber.Ctx) error {
	serviceID, err := ctx.ParamsInt("id")
	if err != nil {
		return response.ErrBadRequest.AddMessage(err).Send(ctx)
	}

	service := Service{}
	err = postgres.Conn.QueryRow(ctx.Context(),
		"select id, title, description, info_href, date_update, actual from public.services where id = $1",
		serviceID).Scan(
		&service.ID,
		&service.Title,
		&service.Description,
		&service.InfoHref,
		&service.DateUpdate,
		&service.Actual,
	)
	if err != nil {
		return response.ErrInternal.AddMessage(err).Send(ctx)
	}

	return response.OK.SetData(service).Send(ctx)
}

func HandlerAddService(ctx *fiber.Ctx) error {
	var request struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		InfoHref    string `json:"info_href"`
	}
	err := ctx.BodyParser(&request)
	if err != nil {
		return response.ErrBadRequest.AddMessage(err).Send(ctx)
	}

	var serviceID int
	var dateUpdate = time.Now()
	err = postgres.Conn.QueryRow(ctx.Context(),
		"INSERT INTO public.services (title, description, info_href, date_update, actual) VALUES ($1, $2, $3, $4, true) RETURNING id",
		request.Title, request.Description, request.InfoHref, dateUpdate).Scan(&serviceID)
	if err != nil {
		return response.ErrInternal.AddMessage(err).Send(ctx)
	}

	return response.OK.SetData(Service{
		ID:          serviceID,
		Title:       request.Title,
		Description: request.Description,
		InfoHref:    request.InfoHref,
		DateUpdate:  dateUpdate,
		Actual:      true,
	}).Send(ctx)
}

func HandlerUpdateService(ctx *fiber.Ctx) error {
	var request Service
	err := ctx.BodyParser(&request)
	if err != nil {
		return response.ErrBadRequest.AddMessage(err).Send(ctx)
	}

	serviceID, err := ctx.ParamsInt("id")
	if err != nil {
		return response.ErrBadRequest.AddMessage(err).Send(ctx)
	}

	dateUpdate := time.Now()
	_, err = postgres.Conn.Exec(ctx.Context(), "UPDATE public.services SET title=$1, description=$2, info_href=$3, date_update=$4, actual=true WHERE id=$5",
		request.Title, request.Description, request.InfoHref, dateUpdate, serviceID)
	if err != nil {
		return response.ErrInternal.AddMessage(err).Send(ctx)
	}

	return response.OK.SetData(Service{
		ID:          serviceID,
		Title:       request.Title,
		Description: request.Description,
		InfoHref:    request.InfoHref,
		DateUpdate:  dateUpdate,
		Actual:      true,
	}).Send(ctx)
}

func HandlerDeleteService(ctx *fiber.Ctx) error {
	serviceID, err := ctx.ParamsInt("id")
	if err != nil {
		return response.ErrBadRequest.AddMessage(err).Send(ctx)
	}

	dateUpdate := time.Now()
	_, err = postgres.Conn.Exec(ctx.Context(), "UPDATE public.services SET date_update=$1, actual=false WHERE id=$2",
		dateUpdate, serviceID)
	if err != nil {
		return response.ErrInternal.AddMessage(err).Send(ctx)
	}

	return response.OK.Send(ctx)
}
