package repositories

import (
	"fmt"

	"github.com/mich31/scoreplay-media-api/models"
	"gorm.io/gorm"
)

type IMediaRepository interface {
	Create(media *models.Media, tags []string) (uint, error)
	Delete(id string) error
	Find() ([]*models.Media, error)
	FindByTag(tag string) ([]*models.Media, error)
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
		return 0, result.Error
	}
	//If the media already exists
	if result.RowsAffected == 0 {
		return 0, nil // TODO
	}

	err := repository.db.Transaction(func(tx *gorm.DB) error {
		// Verify all tags exist
		var tags []*models.Tag
		if err := tx.Model(&models.Tag{}).Where("name IN ?", tagNames).Find(&tags).Error; err != nil {
			return fmt.Errorf("Unable to check tags: %w", err)
		}
		if len(tags) != len(tagNames) {
			return fmt.Errorf("Some tags do not exist")
		}

		for _, tag := range tags {
			mediaTag := models.MediaTag{
				MediaID: media.ID,
				TagID:   tag.ID,
			}

			if err := tx.Create(&mediaTag).Error; err != nil {
				return fmt.Errorf("An error occured creating media-tag association: %w", err)
			}
		}

		return nil
	})
	if err != nil {
		fmt.Errorf("Error: %s", err.Error)
	}
	return media.ID, err // TODO
}

func (repository *MediaRepository) Find() ([]*models.Media, error) {
	return nil, nil
}

func (repository *MediaRepository) FindByTag(tag string) ([]*models.Media, error) {
	return nil, nil
}

func (repository *MediaRepository) Delete(id string) error {
	return nil
}
