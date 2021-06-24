package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/pratikjagrut/myreads-go-backend/controllers"
)

func Setup(app *fiber.App) {
	// User apis
	app.Post("/api/register", controllers.Register)
	app.Post("/api/login", controllers.Login)
	app.Get("/api/user", controllers.User)
	app.Post("/api/logout", controllers.Logout)

	// Book apis
	app.Post("/api/books/add", controllers.BookEntry)
	app.Get("/api/books/all", func(c *fiber.Ctx) error {
		return controllers.GetBoooks(c, nil)
	})
	app.Get("/api/books/reading", func(c *fiber.Ctx) error {
		return controllers.GetBoooks(c, &controllers.Reading)
	})
	app.Get("/api/books/finished", func(c *fiber.Ctx) error {
		return controllers.GetBoooks(c, &controllers.Finished)
	})
	app.Get("/api/books/wishlist", func(c *fiber.Ctx) error {
		return controllers.GetBoooks(c, &controllers.Wishlist)
	})
	app.Post("/api/books/updatestatus", controllers.UpdateStatus)

	app.Post("/api/books/deletebook", controllers.DeleteBook)
}
