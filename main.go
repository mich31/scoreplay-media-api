package main

import (
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New(fiber.Config{
		AppName: "ScorePlay Media API v0.1",
	})

	app.Listen(":3000")
}
