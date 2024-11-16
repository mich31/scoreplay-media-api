package controllers

import (
	"io"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/mich31/scoreplay-media-api/models"
	"github.com/mich31/scoreplay-media-api/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockMediaRepository struct {
	mock.Mock
}

func (m *mockMediaRepository) Create(media *models.Media, tags []string) (uint, error) {
	return 1, nil
}

func (m *mockMediaRepository) FindByTag(tag string) ([]models.MediaWithTagNames, error) {
	return []models.MediaWithTagNames{}, nil
}

func TestMediaRoutes(t *testing.T) {
	tests := []struct {
		description  string
		route        string
		method       string
		body         io.Reader
		expectedCode int
	}{
		{
			description:  "Get all medias should return HTTP status 501",
			route:        "/api/medias",
			method:       "GET",
			body:         nil,
			expectedCode: 501,
		},
		{
			description:  "Create media should return HTTP status 501",
			route:        "/api/medias",
			method:       "POST",
			body:         nil,
			expectedCode: 501,
		},
	}

	mediaRepository := &mockMediaRepository{}
	tagRepository := &mockTagRepository{}
	mediaService := services.NewMediaService(mediaRepository, tagRepository)
	mediaController := NewMediaController(*mediaService)

	app := fiber.New()
	api := app.Group("/api")

	// routes
	api.Route("medias", func(router fiber.Router) {
		router.Get("/", mediaController.GetMedias)
		router.Post("/", mediaController.CreateMedia)
	})

	// Iterate through test single test cases
	for _, test := range tests {
		req := httptest.NewRequest(test.method, test.route, test.body)
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req, 1)

		// Verify, if the status code is as expected.
		assert.Equalf(t, test.expectedCode, resp.StatusCode, test.description)
	}
}
