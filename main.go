package main

import (
	"log"
	"time"

	"github.com/gofiber/contrib/swagger"
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
	mediaRepository := repositories.NewMediaRepository(db)
	tagService := services.NewTagService(tagRepository)
	storageService := services.InitStorageService()
	mediaService := services.NewMediaService(mediaRepository, tagRepository, storageService)
	tagController := controllers.NewTagController(*tagService)
	mediaController := controllers.NewMediaController(*mediaService)

	app := fiber.New(fiber.Config{
		AppName: "ScorePlay Media API v0.1",
	})

	app.Use(logger.New())
	app.Use(healthcheck.New())

	cfg := swagger.Config{
		BasePath: "/",
		FilePath: "./docs/swagger.json",
		Path:     "swagger",
		Title:    "Swagger API Docs",
	}

	app.Use(swagger.New(cfg))

	api := app.Group("/api")

	// routes
	api.Get("/health", HealthCheck)
	api.Route("tags", func(router fiber.Router) {
		router.Get("/", tagController.GetTags)
		router.Post("/", tagController.CreateTag)
		router.Delete("/:id", tagController.DeleteTag)
	})
	api.Route("medias", func(router fiber.Router) {
		router.Get("/", mediaController.GetMedias)
		router.Post("/", mediaController.CreateMedia)
	})

	if err := app.Listen(":3000"); err != nil {
		log.Fatal("Error starting server:", err)
	}
}

// Healthcheck godoc
//
//	@Summary		Healthcheck endpoint
//	@Description	Healthcheck endpoint
//	@Tags			Health
//	@Produce		json
//	@Success		200	{object}	main.HealthCheck.response
//	@Router			/api/health [get]
func HealthCheck(c *fiber.Ctx) error {
	type response struct {
		Status string `json:"status"`
		Date   string `json:"date"`
	}
	return c.Status(200).JSON(response{
		Status: "OK",
		Date:   time.Now().Format("2006-01-02 15:04:05"),
	})
}
