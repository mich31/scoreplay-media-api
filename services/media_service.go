package services

import (
	"context"
	"fmt"
	"log"
	"mime/multipart"

	"github.com/mich31/scoreplay-media-api/config"
	"github.com/mich31/scoreplay-media-api/models"
	"github.com/mich31/scoreplay-media-api/repositories"
)

type MediaService struct {
	mediaRepository repositories.IMediaRepository
	tagRepository   repositories.ITagRepository
	storage         StorageService
}

func NewMediaService(mediaRepository repositories.IMediaRepository, tagRepository repositories.ITagRepository) *MediaService {
	ctx := context.Background()
	storage, err := NewStorageService()
	if err != nil {
		log.Fatalf("Unable to initialize storage service: %s", err)
	}
	storage.CreateBucket(ctx, config.Config("STORAGE_BUCKET_NAME"))
	return &MediaService{
		mediaRepository: mediaRepository,
		tagRepository:   tagRepository,
		storage:         *storage,
	}
}

func (service *MediaService) CreateMedia(ctx context.Context, name string, tags []string, file *multipart.FileHeader) (uint, error) {
	// TODO find media by name
	fileUrl, err := service.storage.UploadObject(ctx, file)
	if err != nil {
		return 0, err
	}
	fmt.Printf("File uploaded at: %s\n", fileUrl)
	media := &models.Media{
		Name:     name,
		FileUrl:  fileUrl,
		FileSize: file.Size,
	}
	id, err := service.mediaRepository.Create(media, tags)
	if err != nil {
		fmt.Printf("Unable to create media %s: %s\n", name, err.Error())
		return 0, err
	} else if id == 0 {
		fmt.Errorf("Media %s already exists", media.Name)
		return 0, err
	}
	fmt.Printf("Media %s created\n", name)
	return id, nil
}
