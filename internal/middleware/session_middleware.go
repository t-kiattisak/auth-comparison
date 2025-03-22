package middleware

import (
	"auth-comparison/internal/domain"

	"github.com/gofiber/fiber/v2"
)

type SessionValidator interface {
	Validate(c *fiber.Ctx) (domain.User, error)
}

func SessionMiddleware(service SessionValidator) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user, err := service.Validate(c)
		if err != nil {
			return c.Status(401).JSON(fiber.Map{"error": "invalid session"})
		}
		c.Locals("user", user)
		return c.Next()
	}
}
