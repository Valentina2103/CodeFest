package routes

import (
	"github.com/gofiber/fiber/v2"
)

func videoLive(c *fiber.Ctx) error {
	err := verifyLogin(c)
	if err != nil {
		return nil
	}

	

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Video Live",
	})
}
