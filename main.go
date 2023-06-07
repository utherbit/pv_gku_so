package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"pahomov_frolovsky_cson/handlers"
	"pahomov_frolovsky_cson/postgres"
	"pahomov_frolovsky_cson/utilities"
)

func main() {
	utilities.CheckEnvFile()
	postgres.InitPostgresConnection()
	InitFiber()
}

func InitFiber() {
	app := fiber.New()
	app.Use(cors.New(cors.ConfigDefault))

	// Просмотр списка услуг
	app.Get("/services", handlers.HandlerGetServices)
	// Оставить заявку по услуге
	app.Post("/requests", handlers.HandlerAddServiceRequest)

	// Авторизация
	app.Post("/login", handlers.HandlerLogin)
	app.Use(handlers.AuthMiddleware)

	// Добавление новой услуги
	app.Post("/services", handlers.HandlerAddService)
	// Получение информации об услуге по id
	app.Get("/services/:id", handlers.HandlerGetService)
	// Изменение услуги по id
	app.Patch("/services/:id", handlers.HandlerUpdateService)
	// Удаление услуги по id
	app.Delete("/services/:id", handlers.HandlerDeleteService)

	// Просмотр списка заявок
	app.Get("/requests", handlers.HandlerGetServiceRequests)
	// Изменение статуса заявки
	app.Patch("/requests/:id", handlers.HandlerUpdateRequestStatus)

	var addr string
	utilities.LookupEnv(&addr, "LISTEN_ADDR")
	app.Listen(addr)
}
