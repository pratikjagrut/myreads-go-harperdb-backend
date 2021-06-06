package main

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/pratikjagrut/go-jwt-auth-harperDB/database"
	"github.com/pratikjagrut/go-jwt-auth-harperDB/routes"
)

func main() {

	database.Init(os.Getenv("HARPERDB_HOST"), os.Getenv("HARPERDB_UNAME"), os.Getenv("HARPERDB_PSWD"), "myreads")

	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowCredentials: true,
	}))

	routes.Setup(app)

	app.Listen(":8000")

}
