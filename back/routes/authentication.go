package routes

import (
	"os"

	"github.com/Valentina2103/CodeFest/back/auth"
	"github.com/Valentina2103/CodeFest/back/models"
	"github.com/gofiber/fiber/v2"
)

func login(c *fiber.Ctx) error {
	// get from env the variable cognitoID
	cognitoID := os.Getenv("COGNITO_ID")
	cognitoClient := auth.NewCognitoClient("sa-east-1", cognitoID)

	userCredentials := new(models.User)

	err := c.BodyParser(userCredentials)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Cannot parse JSON",
		})
	}

	token, err := cognitoClient.SignIn(userCredentials.Email, userCredentials.Password)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid credentials",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"token": token,
	})
}

func register(c *fiber.Ctx) error {

	cognitoID := os.Getenv("COGNITO_ID")
	cognitoClient := auth.NewCognitoClient("sa-east-1", cognitoID)

	userCredentials := new(models.User)

	err := c.BodyParser(userCredentials)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Cannot parse JSON",
		})
	}

	err = cognitoClient.SignUp(userCredentials.Email, userCredentials.Password)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Cannot create user",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "User created",
	})
}

func verifyLogin(c *fiber.Ctx) error {
	bearerToken := c.Get("Authorization")
	if bearerToken == "" {
		c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "No token provided",
		})
		return fiber.ErrUnauthorized
	}

	// check that the bearerToken starts with the word Bearer
	if len(bearerToken) < 7 || bearerToken[:7] != "Bearer " {
		c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid token",
		})
		return fiber.ErrUnauthorized
	}
	// take the token after the Bearer word splited by space
	token := bearerToken[7:]

	cognitoID := os.Getenv("COGNITO_ID")
	cognitoClient := auth.NewCognitoClient("sa-east-1", cognitoID)

	err := cognitoClient.VerifyToken(token)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid token",
		})
	}

	return nil
}
