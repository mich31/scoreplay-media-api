package controllers

import (
	"io"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestTagRoutes(t *testing.T) {
	tests := []struct {
		description  string
		route        string
		method       string
		body         io.Reader
		expectedCode int
	}{
		{
			description:  "Get all tags should return HTTP status 501",
			route:        "/api/tags",
			method:       "GET",
			body:         nil,
			expectedCode: 501,
		},
		{
			description:  "Create tag should return HTTP status 501",
			route:        "/api/tags",
			method:       "POST",
			body:         nil,
			expectedCode: 501,
		},
		{
			description:  "Update tag should return HTTP status 501",
			route:        "/api/tags/1",
			method:       "PATCH",
			body:         nil,
			expectedCode: 501,
		},
		{
			description:  "Delete tag should return HTTP status 501",
			route:        "/api/tags/1",
			method:       "DELETE",
			body:         nil,
			expectedCode: 501,
		},
	}

	app := fiber.New()
	api := app.Group("/api")

	// routes
	api.Route("tags", func(router fiber.Router) {
		router.Get("/", GetTags)
		router.Post("/", CreateTag)
		router.Patch("/:id", UpdateTag)
		router.Delete("/:id", DeleteTag)
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
