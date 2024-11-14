package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mich31/scoreplay-media-api/models"
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
	name := c.Query("name")
	results, err := ctrl.service.GetTags(name)
	if err != nil {
		return c.Status(500).JSON(&fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}
	return c.Status(200).JSON(&fiber.Map{
		"success": true,
		"data":    results,
	})
}

func (ctrl TagController) CreateTag(c *fiber.Ctx) error {
	input := new(models.Tag)
	if err := c.BodyParser(input); err != nil {
		return c.Status(400).JSON(&fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	id, err := ctrl.service.CreateTag(input)
	if err != nil {
		return c.Status(500).JSON(&fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}
	return c.Status(201).JSON(&fiber.Map{
		"success": true,
		"id":      id,
	})
}

func (ctrl TagController) UpdateTag(c *fiber.Ctx) error {
	return c.Status(501).JSON(&fiber.Map{
		"success": false,
		"message": "Not implemented yet",
	})
}

func (ctrl TagController) DeleteTag(c *fiber.Ctx) error {
	id := c.Params("id")
	err := ctrl.service.DeleteTag(id)
	if err != nil {
		return c.Status(500).JSON(&fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}
	return c.Status(200).JSON(&fiber.Map{
		"success": true,
	})
}
