package main

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/pratikjagrut/myreads-go-harperdb-backend/database"
	"github.com/pratikjagrut/myreads-go-harperdb-backend/routes"
)

func main() {

	database.Init(os.Getenv("HARPERDB_HOST"), os.Getenv("HARPERDB_UNAME"), os.Getenv("HARPERDB_PSWD"), "myreads")

	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowCredentials: true,
	}))

	routes.Setup(app)
	app.Static("/api/static", "./"+os.Getenv("IMAGES_DIR"))
	app.Listen(":" + os.Getenv("PORT"))

}
