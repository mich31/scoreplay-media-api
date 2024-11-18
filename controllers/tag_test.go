package controllers

import (
	"errors"
	"io"
	"net/http/httptest"
	"strings"
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
	args := m.Called(tag)
	return args.Get(0).(uint), args.Error(1)
}

func (m *mockTagRepository) Find() ([]*models.Tag, error) {
	args := m.Called()
	return args.Get(0).([]*models.Tag), args.Error(1)
}

func (m *mockTagRepository) FindByName(name string) ([]*models.Tag, error) {
	args := m.Called(name)
	return args.Get(0).([]*models.Tag), args.Error(1)
}

func (m *mockTagRepository) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestGetTags(t *testing.T) {
	tests := []struct {
		description          string
		tagName              string
		mockTags             []*models.Tag
		mockError            error
		expectedStatusCode   int
		expectedBodyResponse string
	}{
		{
			description: "Get all tags should return a list of tags and HTTP status code 200",
			tagName:     "",
			mockTags: []*models.Tag{
				{ID: 1, Name: "Zidane"},
				{ID: 2, Name: "Brady"},
				{ID: 3, Name: "Hamilton"},
			},
			mockError:          nil,
			expectedStatusCode: 200,
			expectedBodyResponse: `{
				"success":true,
				"message":"",
				"data":[
					{"id":1,"name":"Zidane", "createdAt":"0001-01-01T00:00:00Z", "updatedAt":"0001-01-01T00:00:00Z"},
					{"id":2,"name":"Brady", "createdAt":"0001-01-01T00:00:00Z", "updatedAt":"0001-01-01T00:00:00Z"},
					{"id":3,"name":"Hamilton", "createdAt":"0001-01-01T00:00:00Z", "updatedAt":"0001-01-01T00:00:00Z"}
				]}`,
		},
		{
			description: "Get all tags with a query parameter should return a list of tags and HTTP status 200",
			tagName:     "paul",
			mockTags: []*models.Tag{
				{ID: 4, Name: "Paul Scholes"},
				{ID: 7, Name: "Paul Pogba"},
			},
			mockError:          nil,
			expectedStatusCode: 200,
			expectedBodyResponse: `{
				"success":true,
				"message":"",
				"data":[
					{"id":4,"name":"Paul Scholes", "createdAt":"0001-01-01T00:00:00Z", "updatedAt":"0001-01-01T00:00:00Z"},
					{"id":7,"name":"Paul Pogba", "createdAt":"0001-01-01T00:00:00Z", "updatedAt":"0001-01-01T00:00:00Z"}
				]}`,
		},
		{
			description:        "Get all tags with an unexisting tag name as query parameter should return an empty list of tags and HTTP status 200",
			tagName:            "john",
			mockTags:           []*models.Tag{},
			mockError:          nil,
			expectedStatusCode: 200,
			expectedBodyResponse: `{
				"success":true,
				"message":"",
				"data":[]}`,
		},
		{
			description:        "Get all tags should return HTTP status 500 if an unexpected error occurs",
			tagName:            "",
			mockTags:           nil,
			mockError:          errors.New("database unreachable"),
			expectedStatusCode: 500,
			expectedBodyResponse: `{
				"success":false,
				"message":"database unreachable",
				"data":null}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			app := fiber.New()
			api := app.Group("/api")

			mockTagRepository := new(mockTagRepository)
			if tt.mockError != nil {
				mockTagRepository.On("Find").Return(tt.mockTags, tt.mockError)
			} else if tt.tagName != "" {
				mockTagRepository.On("FindByName", tt.tagName).Return(tt.mockTags, tt.mockError)
			} else {
				mockTagRepository.On("Find").Return(tt.mockTags, tt.mockError)
			}
			tagService := services.NewTagService(mockTagRepository)
			tagController := NewTagController(*tagService)

			// routes
			api.Route("tags", func(router fiber.Router) {
				router.Get("/", tagController.GetTags)
			})

			req := httptest.NewRequest("GET", "/api/tags?name="+tt.tagName, nil)
			resp, _ := app.Test(req)

			assert.Equal(t, tt.expectedStatusCode, resp.StatusCode)
			body, _ := io.ReadAll(resp.Body)
			assert.JSONEq(t, tt.expectedBodyResponse, string(body))
			mockTagRepository.AssertExpectations(t)
		})
	}
}

func TestCreateTag(t *testing.T) {
	tests := []struct {
		description          string
		body                 string
		mockId               uint
		mockError            error
		expectedStatusCode   int
		expectedBodyResponse string
	}{
		{
			description:        "Create tag should return id of the tag created and HTTP status code 201",
			body:               `{ "name":"nba" }`,
			mockId:             1,
			mockError:          nil,
			expectedStatusCode: 201,
			expectedBodyResponse: `{
				"success":true,
				"message":"",
				"id": 1
				}`,
		},
		{
			description:        "Create tag should return an error and HTTP status code 400",
			body:               `"nba"`,
			mockId:             1,
			mockError:          nil,
			expectedStatusCode: 400,
			expectedBodyResponse: `{
				"success":false,
				"message":"json: cannot unmarshal string into Go value of type models.Tag",
				"id": 0
				}`,
		},
		{
			description:        "Create tag should return HTTP status code 500 if an unexpected error occurs",
			body:               `{ "name":"nba" }`,
			mockId:             0,
			mockError:          errors.New("database unreachable"),
			expectedStatusCode: 500,
			expectedBodyResponse: `{
				"success":false,
				"message":"database unreachable",
				"id": 0
			}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			app := fiber.New()
			api := app.Group("/api")

			mockTagRepository := new(mockTagRepository)
			mockTagRepository.On("Create", mock.AnythingOfType("*models.Tag")).Return(tt.mockId, tt.mockError)
			tagService := services.NewTagService(mockTagRepository)
			tagController := NewTagController(*tagService)

			// routes
			api.Route("tags", func(router fiber.Router) {
				router.Post("/", tagController.CreateTag)
			})

			req := httptest.NewRequest("POST", "/api/tags", strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")
			resp, _ := app.Test(req)

			assert.Equal(t, tt.expectedStatusCode, resp.StatusCode)
			body, _ := io.ReadAll(resp.Body)
			assert.JSONEq(t, tt.expectedBodyResponse, string(body))
			if tt.expectedStatusCode != 400 {
				mockTagRepository.AssertExpectations(t)
			}
		})
	}
}
