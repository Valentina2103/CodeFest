package main

import (
	"log"

	"github.com/Valentina2103/CodeFest/back/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	"github.com/joho/godotenv"
)

func init() {
	// Load env variables
	err := godotenv.Load("back/.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	engine := html.New("./back/video", ".tpl")
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	// Define a route for the home page
	app.Get("/", func(c *fiber.Ctx) error {
		return c.Render("index", nil)
	})

	routes.SetupRoutes(app)

	app.Listen(":8080")
}
