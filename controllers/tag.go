package controllers

import "github.com/gofiber/fiber/v2"

func GetTags(c *fiber.Ctx) error {
	return c.Status(501).JSON(&fiber.Map{
		"success": false,
		"message": "Not implemented yet",
	})
}

func CreateTag(c *fiber.Ctx) error {
	return c.Status(501).JSON(&fiber.Map{
		"success": false,
		"message": "Not implemented yet",
	})
}

func UpdateTag(c *fiber.Ctx) error {
	return c.Status(501).JSON(&fiber.Map{
		"success": false,
		"message": "Not implemented yet",
	})
}

func DeleteTag(c *fiber.Ctx) error {
	return c.Status(501).JSON(&fiber.Map{
		"success": false,
		"message": "Not implemented yet",
	})
}
