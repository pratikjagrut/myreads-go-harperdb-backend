package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/pratikjagrut/go-jwt-auth-harperDB/database"
	"github.com/pratikjagrut/go-jwt-auth-harperDB/routes"
)

func main() {

	database.Init("https://app-1-febas.harperdbcloud.com", "HDB_ADMIN", "password", "goodreads")

	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowCredentials: true,
	}))

	routes.Setup(app)

	app.Listen(":8000")

}
