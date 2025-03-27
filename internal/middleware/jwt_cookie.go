package middleware

import (
	"auth-comparison/internal/domain"

	"github.com/gofiber/fiber/v2"
)

func JWTFromCookie(authService domain.UserService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := c.Cookies("token")
		if token == "" {
			return c.Status(401).JSON(fiber.Map{"error": "missing token cookie"})
		}
		user, err := authService.ValidateToken(token)
		if err != nil {
			return c.Status(401).JSON(fiber.Map{"error": "invalid token"})
		}
		c.Locals("user", user)
		return c.Next()
	}
}
