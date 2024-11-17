package controllers

import (
	"encoding/json"
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/mich31/scoreplay-media-api/models"
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

// GetMedias godoc
//
//	@Summary		Get media files by tag id
//	@Description	Get medias by tag id
//	@Tags			Media
//	@Accept			json
//	@Produce		json
//	@Param			tag	query		string	false	"search by tag id"
//	@Success		200	{object}	controllers.GetMedias.response	"Returns success true and array of medias"
//	@Success		404	{object}	controllers.GetMedias.response	"Returns success true with empty data when no media found"
//	@Failure		500	{object}	controllers.GetMedias.response	"Returns error for internal server error"
//	@Router			/api/medias [GET]
func (ctrl MediaController) GetMedias(c *fiber.Ctx) error {
	type response struct {
		Success bool                       `json:"success"`
		Data    []models.MediaWithTagNames `json:"data"`
		Message string                     `json:"message"`
	}
	tag := c.Query("tag")
	results, err := ctrl.service.GetMediasByTag(tag)
	if err != nil {
		return c.Status(500).JSON(response{
			Success: false,
			Message: err.Error(),
		})
	} else if len(results) == 0 {
		return c.Status(404).JSON(response{
			Success: true,
			Data:    results,
		})
	}

	return c.Status(200).JSON(response{
		Success: true,
		Data:    results,
	})
}

// CreateMedia godoc
//
//	@Summary		Upload a new media file
//	@Description	Upload a new media file to storage and creates a new media entry with file url, name and associated tags
//	@Tags			Media
//	@Accept			multipart/form-data
//	@Produce		json
//	@Param			file	formData	file	true	"Media file to upload"
//	@Param			name	formData	string	true	"Media name"
//	@Param			tags	formData	string	true	"Array of tag IDs (example: [123, 75, 18873])"
//	@Success		201	{object}	controllers.CreateMedia.response	"Returns success true when file is uploaded and a new media is created"
//	@Failure		400	{object}	controllers.CreateMedia.response	"Returns error for missing file or existing media"
//	@Failure		500	{object}	controllers.CreateMedia.response	"Returns error for internal server error"
//	@Router			/api/medias [POST]
func (ctrl MediaController) CreateMedia(c *fiber.Ctx) error {
	type response struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
	}

	name := c.FormValue("name")
	tagsStr := c.FormValue("tags")
	var tags []uint
	if err := json.Unmarshal([]byte(tagsStr), &tags); err != nil {
		return c.Status(400).JSON(response{
			Success: false,
			Message: "Invalid tags format: " + err.Error(),
		})
	}

	file, err := c.FormFile("file")
	if file == nil {
		return c.Status(400).JSON(response{
			Success: false,
			Message: "Missing file to upload",
		})
	}
	if err != nil {
		return c.Status(500).JSON(response{
			Success: false,
			Message: "Failed to process uploaded file: " + err.Error(),
		})
	}

	_, err = ctrl.service.CreateMedia(c.Context(), name, tags, file)
	if err != nil {
		switch {
		case errors.Is(err, repositories.ErrMediaExists):
			return c.Status(400).JSON(response{
				Success: false,
				Message: "Failed to create media: " + err.Error(),
			})
		case errors.Is(err, repositories.ErrMediaCreation), errors.Is(err, repositories.ErrMediaDBOperation):
			return c.Status(500).JSON(response{
				Success: false,
				Message: "Failed to create media: " + err.Error(),
			})
		default:
			return c.Status(500).JSON(response{
				Success: false,
				Message: "internal server error",
			})

		}
	}

	return c.Status(201).JSON(response{
		Success: true,
		Message: "File uploaded",
	})
}
