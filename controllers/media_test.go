package controllers

import (
	"bytes"
	"context"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/mich31/scoreplay-media-api/models"
	"github.com/mich31/scoreplay-media-api/repositories"
	"github.com/mich31/scoreplay-media-api/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockMediaRepository struct {
	mock.Mock
}

type mockStorageService struct {
	mock.Mock
}

func (r *mockMediaRepository) Create(media *models.Media, tagIDs []uint) (uint, error) {
	args := r.Called(media, tagIDs)
	return args.Get(0).(uint), args.Error(1)
}

func (r *mockMediaRepository) FindByTag(tag string) ([]models.MediaWithTagNames, error) {
	args := r.Called(tag)
	return args.Get(0).([]models.MediaWithTagNames), args.Error(1)
}

func (s *mockStorageService) CreateBucket(ctx context.Context, bucketName string) error {
	args := s.Called(ctx)
	return args.Error(1)
}

func (s *mockStorageService) UploadObject(ctx context.Context, fileHeader *multipart.FileHeader) (string, error) {
	args := s.Called(ctx, fileHeader)
	return args.Get(0).(string), args.Error(1)
}

func TestGetMedias(t *testing.T) {
	tests := []struct {
		description          string
		tag                  string
		mockReturn           []models.MediaWithTagNames
		mockError            error
		expectedStatusCode   int
		expectedBodyResponse string
	}{
		{
			description: "Get medias should return a list of medias associated to a given tag and HTTP status code 200",
			tag:         "hernandez",
			mockReturn: []models.MediaWithTagNames{
				{
					ID:          1,
					Name:        "lucas_hernandez",
					Description: "Lucas Hernandez",
					FileUrl:     "http://localhost:9000/medias/611e175c-c0bc-488e-b4b7-f5d005e4fa5b.png",
					TagNames:    []string{"hernandez", "football", "france"},
				},
			},
			mockError:          nil,
			expectedStatusCode: 200,
			expectedBodyResponse: `{
				"success":true,
				"message":"",
				"data":[
					{"id":1,"name":"lucas_hernandez", "description":"Lucas Hernandez", "fileUrl":"http://localhost:9000/medias/611e175c-c0bc-488e-b4b7-f5d005e4fa5b.png", "tagNames": ["hernandez", "football", "france"] }
				]}`,
		},
		{
			description:        "Get medias should return an empty list for an unexisting tag and HTTP status code 404",
			tag:                "unexisting_tag",
			mockReturn:         []models.MediaWithTagNames{},
			mockError:          nil,
			expectedStatusCode: 404,
			expectedBodyResponse: `{
				"success":true,
				"message":"",
				"data":[]}`,
		},
		{
			description:        "Get medias should return HTTP status code 500 if an unexpected error occurs",
			tag:                "unexisting_tag",
			mockReturn:         nil,
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

			mockMediaRepository := new(mockMediaRepository)
			mockMediaRepository.On("FindByTag", tt.tag).Return(tt.mockReturn, tt.mockError)
			mockTagRepository := new(mockTagRepository)
			mockStorageService := new(mockStorageService)
			mediaService := services.NewMediaService(mockMediaRepository, mockTagRepository, mockStorageService)
			mediaController := NewMediaController(*mediaService)

			// routes
			api.Route("medias", func(router fiber.Router) {
				router.Get("/", mediaController.GetMedias)
			})

			req := httptest.NewRequest("GET", "/api/medias?tag="+tt.tag, nil)
			resp, _ := app.Test(req)

			assert.Equal(t, tt.expectedStatusCode, resp.StatusCode)
			body, _ := io.ReadAll(resp.Body)
			assert.JSONEq(t, tt.expectedBodyResponse, string(body))
			mockMediaRepository.AssertExpectations(t)
		})
	}
}

func TestCreateMedia(t *testing.T) {
	tests := []struct {
		description          string
		setupRequest         func() (*http.Request, error)
		mockTagIDs           []uint
		mockFileUrl          string
		mockId               uint
		mockRepositoryError  error
		mockStorageError     error
		expectedStatusCode   int
		expectedBodyResponse string
	}{
		{
			description: "Create media should return HTTP status code 200",
			setupRequest: func() (*http.Request, error) {
				body := &bytes.Buffer{}
				writer := multipart.NewWriter(body)
				part, _ := writer.CreateFormFile("file", "baseball.png")
				part.Write([]byte("baseball game"))
				writer.WriteField("name", "baseball")
				writer.WriteField("tags", "[1,2]")
				writer.Close()

				req := httptest.NewRequest("POST", "/api/medias", body)
				req.Header.Set("Content-Type", writer.FormDataContentType())
				return req, nil
			},
			mockFileUrl:          "http://localhost:9000/medias/611e175c-c0bc-488e-b4b7-f5d005e4fa5b.png",
			mockTagIDs:           []uint{1, 2},
			mockId:               1,
			mockRepositoryError:  nil,
			mockStorageError:     nil,
			expectedStatusCode:   201,
			expectedBodyResponse: `{"success":true,"message":"File uploaded"}`,
		},
		{
			description: "Create media should return HTTP status code 400 if tags parameter is invalid",
			setupRequest: func() (*http.Request, error) {
				body := &bytes.Buffer{}
				writer := multipart.NewWriter(body)
				part, _ := writer.CreateFormFile("file", "baseball.png")
				part.Write([]byte("baseball game"))
				writer.WriteField("name", "test video")
				writer.WriteField("tags", "18")
				writer.Close()

				req := httptest.NewRequest("POST", "/api/medias", body)
				req.Header.Set("Content-Type", writer.FormDataContentType())
				return req, nil
			},
			mockFileUrl:          "",
			mockTagIDs:           nil,
			mockId:               0,
			mockRepositoryError:  nil,
			mockStorageError:     nil,
			expectedStatusCode:   400,
			expectedBodyResponse: `{"success":false,"message":"Invalid tags format: json: cannot unmarshal number into Go value of type []uint"}`,
		},
		{
			description: "Create media should return HTTP status code 400 if tags parameter is invalid",
			setupRequest: func() (*http.Request, error) {
				body := &bytes.Buffer{}
				writer := multipart.NewWriter(body)
				writer.WriteField("name", "rugby")
				writer.WriteField("tags", "[1,2]")
				writer.Close()

				req := httptest.NewRequest("POST", "/api/medias", body)
				req.Header.Set("Content-Type", writer.FormDataContentType())
				return req, nil
			},
			mockFileUrl:          "",
			mockTagIDs:           nil,
			mockId:               0,
			mockRepositoryError:  nil,
			mockStorageError:     nil,
			expectedStatusCode:   400,
			expectedBodyResponse: `{"success":false,"message":"Missing file to upload"}`,
		},
		{
			description: "Create media should return HTTP status code 400 when the media already exists",
			setupRequest: func() (*http.Request, error) {
				body := &bytes.Buffer{}
				writer := multipart.NewWriter(body)
				part, _ := writer.CreateFormFile("file", "baseball.png")
				part.Write([]byte("baseball game"))
				writer.WriteField("name", "baseball")
				writer.WriteField("tags", "[1,2]")
				writer.Close()

				req := httptest.NewRequest("POST", "/api/medias", body)
				req.Header.Set("Content-Type", writer.FormDataContentType())
				return req, nil
			},
			mockFileUrl:          "http://localhost:9000/medias/611e175c-c0bc-488e-b4b7-f5d005e4fa5b.png",
			mockTagIDs:           []uint{1, 2},
			mockId:               1,
			mockRepositoryError:  repositories.ErrMediaExists,
			mockStorageError:     nil,
			expectedStatusCode:   400,
			expectedBodyResponse: `{"success":false,"message":"Failed to create media: a media with the same name already exists"}`,
		},
		{
			description: "Create media should return HTTP status code 500 if an unexpected error occurs",
			setupRequest: func() (*http.Request, error) {
				body := &bytes.Buffer{}
				writer := multipart.NewWriter(body)
				part, _ := writer.CreateFormFile("file", "baseball.png")
				part.Write([]byte("baseball game"))
				writer.WriteField("name", "baseball")
				writer.WriteField("tags", "[1,2]")
				writer.Close()

				req := httptest.NewRequest("POST", "/api/medias", body)
				req.Header.Set("Content-Type", writer.FormDataContentType())
				return req, nil
			},
			mockFileUrl:          "http://localhost:9000/medias/611e175c-c0bc-488e-b4b7-f5d005e4fa5b.png",
			mockTagIDs:           []uint{1, 2},
			mockId:               1,
			mockRepositoryError:  repositories.ErrMediaCreation,
			mockStorageError:     nil,
			expectedStatusCode:   500,
			expectedBodyResponse: `{"success":false,"message":"Failed to create media: failed to create media record"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			app := fiber.New()
			api := app.Group("/api")

			mockMediaRepository := new(mockMediaRepository)
			mockMediaRepository.On("Create", mock.AnythingOfType("*models.Media"), tt.mockTagIDs).Return(tt.mockId, tt.mockRepositoryError)
			mockTagRepository := new(mockTagRepository)
			mockStorageService := new(mockStorageService)
			mockStorageService.On(
				"UploadObject",
				mock.Anything,
				mock.AnythingOfType("*multipart.FileHeader")).
				Return(tt.mockFileUrl, tt.mockStorageError)
			mediaService := services.NewMediaService(mockMediaRepository, mockTagRepository, mockStorageService)
			mediaController := NewMediaController(*mediaService)

			// routes
			api.Route("medias", func(router fiber.Router) {
				router.Post("/", mediaController.CreateMedia)
			})

			req, err := tt.setupRequest()
			assert.NoError(t, err)
			resp, _ := app.Test(req)

			assert.Equal(t, tt.expectedStatusCode, resp.StatusCode)
			body, _ := io.ReadAll(resp.Body)
			assert.JSONEq(t, tt.expectedBodyResponse, string(body))
		})
	}
}
