package handler

import (
	"auth-comparison/internal/domain"
	"auth-comparison/internal/middleware"
	"auth-comparison/internal/usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

func RegisterRoutes(app *fiber.App) {
	authService := usecase.NewJWTService("mysecretkey")
	pasetoService := usecase.NewPasetoService([]byte("supersecretkey1234567890123456"))
	sessionStore := session.New()
	sessionService := usecase.NewSessionService(sessionStore)

	app.Post("/login", func(c *fiber.Ctx) error {
		var creds struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}
		if err := c.BodyParser(&creds); err != nil {
			return err
		}
		token, err := authService.Login(creds.Username, creds.Password)
		if err != nil {
			return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
		}
		return c.JSON(fiber.Map{"token": token})
	})

	app.Get("/me", middleware.JWTMiddleware(authService), func(c *fiber.Ctx) error {
		user := c.Locals("user").(domain.User)
		return c.JSON(user)
	})

	app.Post("/paseto-login", func(c *fiber.Ctx) error {
		var creds struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}
		if err := c.BodyParser(&creds); err != nil {
			return err
		}
		token, err := pasetoService.Login(creds.Username, creds.Password)
		if err != nil {
			return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
		}
		return c.JSON(fiber.Map{"token": token})
	})

	app.Get("/me-paseto", middleware.JWTMiddleware(pasetoService), func(c *fiber.Ctx) error {
		user := c.Locals("user").(domain.User)
		return c.JSON(user)
	})

	app.Post("/session-login", func(c *fiber.Ctx) error {
		var creds struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}
		if err := c.BodyParser(&creds); err != nil {
			return err
		}
		err := sessionService.Login(c, creds.Username, creds.Password)
		if err != nil {
			return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
		}
		return c.JSON(fiber.Map{"message": "session created"})
	})

	app.Get("/me-session", middleware.SessionMiddleware(sessionService), func(c *fiber.Ctx) error {
		user := c.Locals("user").(domain.User)
		return c.JSON(user)
	})
}
