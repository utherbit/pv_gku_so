package response

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"strings"
)

type Response struct {
	// set on fiber code
	StatusCode int `json:"-"`

	Status string `json:"status"`
	// В случае ошибки будет расшифровывать статус код ошибки
	ErrorCode string `json:"error_code,omitempty"`
	// Данные ответа, как пояснение к ошибке так и к удовлетворительному ответу
	Message string `json:"message,omitempty"`
	// Данные, результат выполнения
	Data interface{} `json:"data,omitempty"`
}

func NewResponseOk(statusCode int, message ...string) Response {
	response := Response{
		StatusCode: statusCode,
		Status:     "OK",

		Message: strings.Join(message, "; "),
	}
	return response
}
func NewResponseError(statusCode int, errorCode string, message ...string) Response {
	response := Response{
		StatusCode: statusCode,
		Status:     "ERROR",
		ErrorCode:  errorCode,
		Message:    strings.Join(message, "; "),
	}
	return response
}
func (r Response) AddMessage(s interface{}) Response {
	if r.Message != "" {
		r.Message += "; " + fmt.Sprintf("%s", s)
	} else {
		r.Message = fmt.Sprintf("%s", s)
	}
	return r
}
func (r Response) SetData(data interface{}) Response {
	r.Data = data
	return r
}

func (r Response) Send(c *fiber.Ctx) error {
	c.Status(r.StatusCode)
	return c.JSON(r)
}
