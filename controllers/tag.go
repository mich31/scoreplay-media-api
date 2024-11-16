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

// GetTags godoc
//
//	@Summary		GET tags
//	@Description	Get tags (optional: by name)
//	@Tags			Tag
//	@Accept			json
//	@Produce		json
//	@Param  name  query     string  false "search by tag name"
//	@Success		200	{object}	controllers.GetTags.response "Returns success true and a list of tags found"
//	@Failure		500	{object}	controllers.GetTags.response "Returns error for internal server error"
//	@Router			/api/tags   [GET]
func (ctrl TagController) GetTags(c *fiber.Ctx) error {
	type response struct {
		Success bool          `json:"success"`
		Data    []*models.Tag `json:"data"`
		Message string        `json:"message"`
	}
	name := c.Query("name")
	results, err := ctrl.service.GetTags(name)
	if err != nil {
		return c.Status(500).JSON(response{
			Success: false,
			Message: err.Error(),
		})
	}
	return c.Status(200).JSON(response{
		Success: true,
		Data:    results,
	})
}

// CreateTag godoc
//
//	@Summary		Create a new tag
//	@Description	Creates a new tag
//	@Tags			Tag
//	@Accept			json
//	@Produce		json
//	@Param			tag	body		models.Tag	true	"tag object to be created"
//	@Success		201	{object}	controllers.CreateTag.response	"Returns success true and created tag ID"
//	@Failure		400	{object}	controllers.CreateTag.response	"Returns error for invalid input"
//	@Failure		500	{object}	controllers.CreateTag.response	"Returns error for internal server error"
//	@Router			/api/tags [POST]
func (ctrl TagController) CreateTag(c *fiber.Ctx) error {
	type response struct {
		Success bool   `json:"success"`
		Id      uint   `json:"id"`
		Message string `json:"message"`
	}
	input := new(models.Tag)
	if err := c.BodyParser(input); err != nil {
		return c.Status(400).JSON(response{
			Success: false,
			Message: err.Error(),
		})
	}

	id, err := ctrl.service.CreateTag(input)
	if err != nil {
		return c.Status(500).JSON(response{
			Success: false,
			Message: err.Error(),
		})
	}
	return c.Status(201).JSON(response{
		Success: true,
		Id:      id,
	})
}

// DeleteTag godoc
//
//	@Summary		Delete a tag
//	@Description	Deletes a tag by its id
//	@Tags			Tag
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"Tag id"
//	@Success		200	{object}	controllers.DeleteTag.response	"Returns success true"
//	@Failure		500	{object}	controllers.DeleteTag.response	"Returns error for internal server error"
//	@Router			/api/tags/{id} [DELETE]
func (ctrl TagController) DeleteTag(c *fiber.Ctx) error {
	type response struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
	}
	id := c.Params("id")
	err := ctrl.service.DeleteTag(id)
	if err != nil {
		return c.Status(500).JSON(response{
			Success: false,
			Message: err.Error(),
		})
	}
	return c.Status(200).JSON(response{
		Success: true,
	})
}
