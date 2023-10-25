package routes

import "github.com/gofiber/fiber/v2"

func serveVideo(c *fiber.Ctx) error {
	return c.SendFile("./back/video/video.mp4", true)
}
