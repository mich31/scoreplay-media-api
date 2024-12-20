package services

import (
	"context"
	"fmt"
	"mime/multipart"

	"github.com/mich31/scoreplay-media-api/models"
	"github.com/mich31/scoreplay-media-api/repositories"
)

type MediaService struct {
	mediaRepository repositories.IMediaRepository
	tagRepository   repositories.ITagRepository // TODO
	storage         IStorageService
}

func NewMediaService(mediaRepository repositories.IMediaRepository, tagRepository repositories.ITagRepository, storageService IStorageService) *MediaService {
	return &MediaService{
		mediaRepository: mediaRepository,
		tagRepository:   tagRepository,
		storage:         storageService,
	}
}

func (service *MediaService) CreateMedia(ctx context.Context, name string, tagIDs []uint, file *multipart.FileHeader) (uint, error) {
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
	id, err := service.mediaRepository.Create(media, tagIDs)
	if err != nil {
		fmt.Printf("unable to create media %s: %s\n", name, err.Error())
		return 0, err
	}
	fmt.Printf("Media %s created\n", name)
	return id, nil
}

func (service *MediaService) GetMediasByTag(tag string) ([]models.MediaWithTagNames, error) {
	medias, err := service.mediaRepository.FindByTag(tag)
	if err != nil {
		return nil, err
	}

	return medias, nil
}
