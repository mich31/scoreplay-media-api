package controllers

import (
	"errors"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/mich31/scoreplay-media-api/repositories"
	"github.com/mich31/scoreplay-media-api/services"
)

type MediaController struct {
	service services.MediaService
}

func NewMediaController(service services.MediaService) *MediaController {
	return &MediaController{
		service,
	}
}

func (ctrl MediaController) GetMedias(c *fiber.Ctx) error {
	tag := c.Query("tag")
	results, err := ctrl.service.GetMediasByTag(tag)
	if err != nil {
		return c.Status(500).JSON(&fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	} else if len(results) == 0 {
		return c.Status(404).JSON(&fiber.Map{
			"success": true,
			"data":    results,
		})
	}

	return c.Status(200).JSON(&fiber.Map{
		"success": true,
		"data":    results,
	})
}

func (ctrl MediaController) CreateMedia(c *fiber.Ctx) error {
	name := c.FormValue("name")
	tags := strings.Split(c.FormValue("tags"), ",")
	file, err := c.FormFile("file")
	if file == nil {
		return c.Status(400).JSON(&fiber.Map{
			"success": false,
			"message": "Missing file to upload",
		})
	}
	if err != nil {
		return c.Status(500).JSON(&fiber.Map{
			"success": false,
			"message": "Failed to process uploaded file: " + err.Error(),
		})
	}

	_, err = ctrl.service.CreateMedia(c.Context(), name, tags, file)
	if err != nil {
		switch {
		case errors.Is(err, repositories.ErrMediaExists):
			return c.Status(400).JSON(&fiber.Map{
				"success": false,
				"message": "Failed to create media: " + err.Error(),
			})
		case errors.Is(err, repositories.ErrMediaCreation), errors.Is(err, repositories.ErrMediaDBOperation):
			return c.Status(500).JSON(&fiber.Map{
				"success": false,
				"message": "Failed to create media: " + err.Error(),
			})
		default:
			return c.Status(500).JSON(&fiber.Map{
				"success": false,
				"message": "internal server error",
			})

		}
	}

	return c.Status(201).JSON(&fiber.Map{
		"success": true,
		"message": "File uploaded",
	})
}
