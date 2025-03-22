package middleware

import (
	"auth-comparison/internal/domain"

	"github.com/gofiber/fiber/v2"
)

func JWTMiddleware(authService domain.AuthService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		header := c.Get("Authorization")
		if len(header) < 8 || header[:7] != "Bearer " {
			return c.Status(401).JSON(fiber.Map{"error": "invalid header"})
		}

		token := header[7:]
		user, err := authService.ValidateToken(token)
		if err != nil {
			return c.Status(401).JSON(fiber.Map{"error": "invalid token"})
		}

		c.Locals("user", user)
		return c.Next()
	}
}
