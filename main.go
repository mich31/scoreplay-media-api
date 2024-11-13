package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/mich31/scoreplay-media-api/controllers"
)

func main() {
	app := fiber.New(fiber.Config{
		AppName: "ScorePlay Media API v0.1",
	})

	app.Use(logger.New())
	app.Use(healthcheck.New())

	api := app.Group("/api", logger.New())

	// routes
	api.Route("tags", func(router fiber.Router) {
		router.Get("/", controllers.GetTags)
		router.Post("/", controllers.CreateTag)
		router.Patch("/:id", controllers.UpdateTag)
		router.Delete("/:id", controllers.DeleteTag)
	})
	api.Route("medias", func(router fiber.Router) {
		router.Get("/", controllers.GetMedias)
		router.Post("/", controllers.CreateMedia)
		router.Patch("/:id", controllers.UpdateMedia)
		router.Delete("/:id", controllers.DeleteMedia)
	})

	app.Listen(":3000")
}
