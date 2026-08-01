package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	app.Get("/api/v1/ping", func (c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Yeah!, Fiber is running.",
		})
	})

	log.Println("Server is running on port:8080")
	app.Listen(":8080")
}