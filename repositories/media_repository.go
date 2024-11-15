package repositories

import (
	"errors"
	"fmt"

	"github.com/mich31/scoreplay-media-api/models"
	"gorm.io/gorm"
)

var (
	ErrMediaCreation    = errors.New("failed to create media record")
	ErrMediaExists      = errors.New("a media with the same name already exists")
	ErrMediaDBOperation = errors.New("database operation failed")
)

type IMediaRepository interface {
	Create(media *models.Media, tags []string) (uint, error)
	Delete(id string) error
	FindByTag(tag string) ([]models.MediaWithTagNames, error)
}

type MediaRepository struct {
	db *gorm.DB
}

func NewMediaRepository(db *gorm.DB) *MediaRepository {
	return &MediaRepository{db: db}
}

func (repository *MediaRepository) Create(media *models.Media, tagNames []string) (uint, error) {
	result := repository.db.Where(models.Media{Name: media.Name}).FirstOrCreate(media)
	if result.Error != nil {
		return 0, fmt.Errorf("%w: %v", ErrMediaCreation, result.Error)
	}
	//If the media already exists
	if result.RowsAffected == 0 {
		return 0, fmt.Errorf("%w: media with name '%s'", ErrMediaExists, media.Name)
	}

	err := repository.db.Transaction(func(tx *gorm.DB) error {
		// Verify all tags exist
		var tags []*models.Tag
		if err := tx.Model(&models.Tag{}).Where("name IN ?", tagNames).Find(&tags).Error; err != nil {
			return fmt.Errorf("unable to check tags: %w", err)
		}
		if len(tags) != len(tagNames) {
			return fmt.Errorf("some tags do not exist")
		}

		for _, tag := range tags {
			mediaTag := models.MediaTag{
				MediaID: media.ID,
				TagID:   tag.ID,
			}

			if err := tx.Create(&mediaTag).Error; err != nil {
				return fmt.Errorf("an error occured creating media-tag association: %w", err)
			}
		}

		return nil
	})
	if err != nil {
		return 0, fmt.Errorf("%w: %w", ErrMediaDBOperation, err)
	}
	return media.ID, err
}

func (repository *MediaRepository) FindByTag(tag string) ([]models.MediaWithTagNames, error) {
	var mediaIDs []uint
	err := repository.db.Model(&models.MediaTag{}).
		Select("media_id").
		Where("tag_id = ?", tag).
		Find(&mediaIDs).Error

	if err != nil { // TODO
		return nil, err
	}

	medias := []models.MediaWithTagNames{}
	for _, mediaId := range mediaIDs {
		media := models.MediaWithTagNames{}
		err = repository.db.Model(&models.Media{}).
			Select("DISTINCT media.id, media.name, media.description, media.file_url, array_agg(tags.name) as tag_names").
			Joins("JOIN media_tags ON media_tags.media_id = media.id").
			Joins("JOIN tags ON tags.id = media_tags.tag_id").
			Where("media_tags.media_id = ?", mediaId).
			Group("media.id").
			First(&media).Error
		if err != nil {
			fmt.Errorf("unable to fetch media %s", mediaId)
		}
		medias = append(medias, media)
	}

	return medias, err
}

func (repository *MediaRepository) Delete(id string) error {
	return nil
}
