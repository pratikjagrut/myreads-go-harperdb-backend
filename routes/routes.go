package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/pratikjagrut/go-jwt-auth-harperDB/controllers"
)

func Setup(app *fiber.App) {
	app.Post("/api/register", controllers.Register)
	app.Post("/api/login", controllers.Login)
	app.Get("/api/user", controllers.User)
	app.Post("/api/logout", controllers.Logout)

	app.Post("/api/addbook", controllers.AddBook)
	app.Get("/api/getallbooks", controllers.GetAllBooks)
	app.Get("/api/getfinshedbooks", controllers.GetFinshedBooks)
}
