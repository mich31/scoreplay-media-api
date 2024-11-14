package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/mich31/scoreplay-media-api/controllers"
	"github.com/mich31/scoreplay-media-api/database"
	"github.com/mich31/scoreplay-media-api/repositories"
	"github.com/mich31/scoreplay-media-api/services"
)

func main() {
	// Connect to database
	db, err := database.Connect()
	if err != nil {
		log.Fatal(err)
	}

	tagRepository := repositories.NewTagRepository(db)
	tagService := services.NewTagService(tagRepository)
	tagController := controllers.NewTagController(*tagService)

	app := fiber.New(fiber.Config{
		AppName: "ScorePlay Media API v0.1",
	})

	app.Use(logger.New())
	app.Use(healthcheck.New())

	api := app.Group("/api", logger.New())

	// routes
	api.Route("tags", func(router fiber.Router) {
		router.Get("/", tagController.GetTags)
		router.Post("/", tagController.CreateTag)
		router.Patch("/:id", tagController.UpdateTag)
		router.Delete("/:id", tagController.DeleteTag)
	})
	api.Route("medias", func(router fiber.Router) {
		router.Get("/", controllers.GetMedias)
		router.Post("/", controllers.CreateMedia)
		router.Patch("/:id", controllers.UpdateMedia)
		router.Delete("/:id", controllers.DeleteMedia)
	})

	app.Listen(":3000")
}
