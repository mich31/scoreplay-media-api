package repositories

import (
	"github.com/mich31/scoreplay-media-api/models"
	"gorm.io/gorm"
)

type TagRepository struct {
	db *gorm.DB
}

func NewTagRepository(db *gorm.DB) *TagRepository {
	return &TagRepository{db: db}
}

func (repository *TagRepository) CreateTag(tag *models.Tag) error {
	return repository.db.Create(tag).Error
}
