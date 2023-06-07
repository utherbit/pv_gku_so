package handlers

import (
	"github.com/gofiber/fiber/v2"
	"pahomov_frolovsky_cson/handlers/response"
	"pahomov_frolovsky_cson/postgres"
	"time"
)

// ServiceRequest represents the structure of a service request
type ServiceRequest struct {
	ID           int    `json:"id"`
	ServiceID    int    `json:"service_id"`
	ServiceTitle string `json:"service_title,omitempty"`

	LastName          string     `json:"last_name"`
	FirstName         string     `json:"first_name"`
	MiddleName        string     `json:"middle_name"`
	PassportSeries    string     `json:"passport_series"`
	PassportNumber    string     `json:"passport_number"`
	SNILS             string     `json:"snils"`
	RequestText       string     `json:"request_text"`
	RequestDate       time.Time  `json:"request_date"`
	ConsiderationDate *time.Time `json:"consideration_date,omitempty"`
	Status            int        `json:"status"`
}

func HandlerAddServiceRequest(ctx *fiber.Ctx) error {
	// Request
	var request ServiceRequest
	err := ctx.BodyParser(&request)
	if err != nil {
		return response.ErrBadRequest.Send(ctx)
	}

	if len(request.PassportSeries) != 4 {
		return response.ErrBadRequest.AddMessage("Серия паспорта должна быть длинной 4 символа").Send(ctx)
	}
	if len(request.PassportNumber) != 6 {
		return response.ErrBadRequest.AddMessage("Номер паспорта должен быть длинной 6 символов").Send(ctx)
	}
	if len(request.SNILS) != 11 {
		return response.ErrBadRequest.AddMessage("СНИЛС должен быть длинной 11 символов").Send(ctx)
	}

	// Insert postgres
	request.RequestDate = time.Now()
	err = postgres.Conn.QueryRow(ctx.Context(), `
		INSERT INTO public.service_requests 
    	(service_id, last_name, first_name, middle_name, passport_series, passport_number, snils, request_text, request_date)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id`,
		request.ServiceID, request.LastName, request.FirstName, request.MiddleName, request.PassportSeries, request.PassportNumber, request.SNILS,
		request.RequestText, request.RequestDate,
	).Scan(&request.ID)
	if err != nil {
		return response.ErrInternal.AddMessage(err).Send(ctx)
	}

	return response.OK.SetData(request).Send(ctx)
}

func HandlerGetServiceRequests(ctx *fiber.Ctx) error {
	// Request
	all := ctx.QueryBool("all", false)

	// Query format
	query := `
		select 
		    r.id, service_id, s.title, last_name, first_name, middle_name, 
		    passport_series, passport_number, 
		    snils, request_text, 
		    request_date, consideration_date, status
		from public.service_requests r
		inner join public.services s on s.id = r.service_id `
	if !all {
		query += " where r.status = 0 order by status, consideration_date desc , request_date desc"
	} else {
		query += " order by date_update desc"
	}

	// query to postgres
	rows, err := postgres.Conn.Query(ctx.Context(), query)
	if err != nil {
		return response.ErrInternal.AddMessage(err).Send(ctx)
	}
	defer rows.Close()

	// scan
	var requests []ServiceRequest
	for rows.Next() {
		var request ServiceRequest
		err = rows.Scan(
			&request.ID, &request.ServiceID, &request.ServiceTitle,
			&request.LastName, &request.FirstName, &request.MiddleName,
			&request.PassportSeries, &request.PassportNumber, &request.SNILS,
			&request.RequestText, &request.RequestDate, &request.ConsiderationDate, &request.Status)

		if err != nil {
			return response.ErrInternal.AddMessage(err).Send(ctx)
		}
		requests = append(requests, request)
	}

	return response.OK.SetData(requests).Send(ctx)
}

func HandlerUpdateRequestStatus(ctx *fiber.Ctx) error {
	var request struct {
		Status int `json:"status"`
	}
	err := ctx.BodyParser(&request)
	if err != nil {
		return response.ErrBadRequest.AddMessage(err).Send(ctx)
	}
	if request.Status < 0 || request.Status > 2 {
		return response.ErrBadRequest.AddMessage("Статус должен иметь значение 0, 1 или 2").Send(ctx)
	}

	// request
	requestID, err := ctx.ParamsInt("id")
	if err != nil {
		return response.ErrBadRequest.Send(ctx)
	}

	// update in postgres
	_, err = postgres.Conn.Exec(ctx.Context(),
		"UPDATE public.service_requests SET status=$2, consideration_date=now() WHERE id=$1", requestID, request.Status)
	if err != nil {
		return response.ErrInternal.AddMessage(err).Send(ctx)
	}

	return response.OK.Send(ctx)
}
