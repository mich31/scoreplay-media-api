package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mich31/scoreplay-media-api/services"
)

type TagController struct {
	service services.TagService
}

func NewTagController(service services.TagService) *TagController {
	return &TagController{
		service,
	}
}

func (ctrl TagController) GetTags(c *fiber.Ctx) error {
	// var results []models.Tag
	return c.Status(501).JSON(&fiber.Map{
		"success": false,
		"message": "Not implemented yet",
	})
}

func (ctrl TagController) CreateTag(c *fiber.Ctx) error {
	return c.Status(501).JSON(&fiber.Map{
		"success": false,
		"message": "Not implemented yet",
	})
}

func (ctrl TagController) UpdateTag(c *fiber.Ctx) error {
	return c.Status(501).JSON(&fiber.Map{
		"success": false,
		"message": "Not implemented yet",
	})
}

func (ctrl TagController) DeleteTag(c *fiber.Ctx) error {
	return c.Status(501).JSON(&fiber.Map{
		"success": false,
		"message": "Not implemented yet",
	})
}
