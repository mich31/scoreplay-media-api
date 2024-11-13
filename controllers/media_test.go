package controllers

import (
	"io"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

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
		{
			description:  "Update media should return HTTP status 501",
			route:        "/api/medias/1",
			method:       "PATCH",
			body:         nil,
			expectedCode: 501,
		},
		{
			description:  "Delete media should return HTTP status 501",
			route:        "/api/medias/1",
			method:       "DELETE",
			body:         nil,
			expectedCode: 501,
		},
	}

	app := fiber.New()
	api := app.Group("/api")

	// routes
	api.Route("medias", func(router fiber.Router) {
		router.Get("/", GetMedias)
		router.Post("/", CreateMedia)
		router.Patch("/:id", UpdateMedia)
		router.Delete("/:id", DeleteMedia)
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
