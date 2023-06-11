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
	app.Use(cors.New(cors.Config{
		AllowHeaders: "*",
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PATCH,PUT,DELETE",
	}))

	//app.Use(func(ctx *fiber.Ctx) error {
	//	fmt.Printf("\nRequest %s %s%s", ctx.Method(), ctx.BaseURL(), ctx.OriginalURL())
	//	n := ctx.Next()
	//	fmt.Printf("\nResponse %d %s", ctx.Response().StatusCode(), ctx.Response().Body())
	//	return n
	//})
	// Просмотр списка услуг
	app.Get("/services", handlers.HandlerGetServices)
	// Получение информации об услуге по id
	app.Get("/services/:id", handlers.HandlerGetService)

	// Оставить заявку по услуге
	app.Post("/requests", handlers.HandlerAddServiceRequest)

	// Получение новостей
	app.Get("/news", handlers.HandlerGetNewsSlice)
	app.Get("/news/:id", handlers.HandlerGetNews)

	// Получение файлов
	app.Get("/uploads/:fileuid", handlers.HandlerGetFile)

	// Авторизация
	app.Post("/login", handlers.HandlerLogin)

	//app.Use(handlers.AuthMiddlewareAdmin)
	services := app.Group("/services", handlers.AuthMiddlewareAdmin)
	// Добавление новой услуги
	services.Post("", handlers.HandlerAddService)
	// Изменение услуги по id
	services.Patch("/:id", handlers.HandlerUpdateService)
	// Удаление услуги по id
	services.Delete("/:id", handlers.HandlerDeleteService)

	requests := app.Group("/requests", handlers.AuthMiddlewareAdmin)
	// Просмотр списка заявок
	requests.Get("", handlers.HandlerGetServiceRequests)
	// Изменение статуса заявки
	requests.Patch("/:id", handlers.HandlerUpdateRequestStatus)

	news := app.Group("/news", handlers.AuthMiddlewareNewsMaker)
	// Добавление новой новости
	news.Post("", handlers.HandlerAddNewsRequest)
	news.Patch("/:id", handlers.HandlerUpdateNews)
	news.Delete("/:id", handlers.HandlerDeleteNews)

	var addr string
	utilities.LookupEnv(&addr, "LISTEN_ADDR")
	app.Listen(addr)
}
