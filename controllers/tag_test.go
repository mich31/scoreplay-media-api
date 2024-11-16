package controllers

import (
	"bytes"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/mich31/scoreplay-media-api/models"
	"github.com/mich31/scoreplay-media-api/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockTagRepository struct {
	mock.Mock
}

func (m *mockTagRepository) Create(tag *models.Tag) (uint, error) {
	return 1, nil
}

func (m *mockTagRepository) Find() ([]*models.Tag, error) {
	return []*models.Tag{}, nil
}

func (m *mockTagRepository) FindByName(name string) ([]*models.Tag, error) {
	return []*models.Tag{}, nil
}

func (m *mockTagRepository) Delete(id string) error {
	return nil
}

func TestTagRoutes(t *testing.T) {
	tests := []struct {
		description  string
		route        string
		method       string
		body         io.Reader
		expectedCode int
	}{
		{
			description:  "Get all tags should return HTTP status 200",
			route:        "/api/tags",
			method:       "GET",
			body:         nil,
			expectedCode: 200,
		},
		{
			description:  "Get all tags with a query parameter should return HTTP status 200",
			route:        "/api/tags?name=zidane",
			method:       "GET",
			body:         nil,
			expectedCode: 200,
		},
		{
			description:  "Create tag should return HTTP status 201",
			route:        "/api/tags",
			method:       "POST",
			body:         bytes.NewBuffer([]byte(`{"name":"Zidane"}`)),
			expectedCode: 201,
		},
		{
			description:  "Delete tag should return HTTP status 501",
			route:        "/api/tags/1",
			method:       "DELETE",
			body:         nil,
			expectedCode: 200,
		},
	}

	tagRepository := &mockTagRepository{}
	tagService := services.NewTagService(tagRepository)
	tagController := NewTagController(*tagService)

	app := fiber.New()
	api := app.Group("/api")

	// routes
	api.Route("tags", func(router fiber.Router) {
		router.Get("/", tagController.GetTags)
		router.Post("/", tagController.CreateTag)
		router.Delete("/:id", tagController.DeleteTag)
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
