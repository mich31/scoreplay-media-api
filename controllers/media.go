package controllers

import "github.com/gofiber/fiber/v2"

func GetMedias(c *fiber.Ctx) error {
	return c.Status(501).JSON(&fiber.Map{
		"success": false,
		"message": "Not implemented yet",
	})
}

func CreateMedia(c *fiber.Ctx) error {
	return c.Status(501).JSON(&fiber.Map{
		"success": false,
		"message": "Not implemented yet",
	})
}

func UpdateMedia(c *fiber.Ctx) error {
	return c.Status(501).JSON(&fiber.Map{
		"success": false,
		"message": "Not implemented yet",
	})
}

func DeleteMedia(c *fiber.Ctx) error {
	return c.Status(501).JSON(&fiber.Map{
		"success": false,
		"message": "Not implemented yet",
	})
}
