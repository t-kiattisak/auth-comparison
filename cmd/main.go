package main

import (
	"auth-comparison/internal/handler"
	"log"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	handler.RegisterRoutes(app)

	log.Fatal(app.Listen(":8080"))
}
