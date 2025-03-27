package middleware

import (
	"auth-comparison/internal/domain"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func JWTAuth(secret string, service domain.AuthService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		header := c.Get("Authorization")
		if header == "" || !strings.HasPrefix(header, "Bearer ") {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "missing or invalid token"})
		}

		tokenStr := strings.TrimPrefix(header, "Bearer ")
		user, err := service.ValidateToken(tokenStr)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
		}

		c.Locals("user", user)
		return c.Next()
	}
}
