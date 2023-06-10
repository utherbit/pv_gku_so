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

	//app.Use(handlers.AuthMiddlewareAdmin)
	// Добавление новой услуги
	app.Post("/services", handlers.AuthMiddlewareAdmin, handlers.HandlerAddService)
	// Получение информации об услуге по id
	app.Get("/services/:id", handlers.AuthMiddlewareAdmin, handlers.HandlerGetService)
	// Изменение услуги по id
	app.Patch("/services/:id", handlers.AuthMiddlewareAdmin, handlers.HandlerUpdateService)
	// Удаление услуги по id
	app.Delete("/services/:id", handlers.AuthMiddlewareAdmin, handlers.HandlerDeleteService)

	// Просмотр списка заявок
	app.Get("/requests", handlers.AuthMiddlewareAdmin, handlers.HandlerGetServiceRequests)
	// Изменение статуса заявки
	app.Patch("/requests/:id", handlers.AuthMiddlewareAdmin, handlers.HandlerUpdateRequestStatus)

	//app.Use(handlers.AuthMiddlewareNewsMaker)
	// Добавление новой новости
	app.Post("/news", handlers.AuthMiddlewareNewsMaker, handlers.HandlerAddNewsRequest)
	app.Patch("/news/:id", handlers.AuthMiddlewareNewsMaker, handlers.HandlerUpdateNews)

	var addr string
	utilities.LookupEnv(&addr, "LISTEN_ADDR")
	app.Listen(addr)
}
