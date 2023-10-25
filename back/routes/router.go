package routes

import (
	"github.com/gofiber/fiber/v2"
)

func holamundo(c *fiber.Ctx) error {
	return c.SendString("Server is running! Send your request")
}

// SetupRoutes is a function that defines the routes of the application
func SetupRoutes(app *fiber.App) {
	// AuthRoutes
	app.Post("/login", login)
	app.Post("/register", register)

	// VideoRoutes
	app.Get("/live", videoLive)
	app.Get("/videos", holamundo)
	app.Get("/video/:id", serveVideo)

	// ScheduleRoutes
	app.Get("/schedule", holamundo)
	app.Post("schedule", holamundo)
	app.Put("/schedule/:id", holamundo)
	app.Delete("/schedule/:id", holamundo)

	// TranscriptionRoutes
	app.Get("/transcription/:id", holamundo)
}
