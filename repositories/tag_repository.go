package repositories

import (
	"github.com/mich31/scoreplay-media-api/models"
	"gorm.io/gorm"
)

type ITagRepository interface {
	Create(tag *models.Tag) (uint, error)
	Delete(id string) error
	Find() ([]*models.Tag, error)
	FindByName(name string) ([]*models.Tag, error)
}

type TagRepository struct {
	db *gorm.DB
}

func NewTagRepository(db *gorm.DB) *TagRepository {
	return &TagRepository{db: db}
}

func (repository *TagRepository) Create(tag *models.Tag) (uint, error) {
	result := repository.db.Where(models.Tag{Name: tag.Name}).FirstOrCreate(tag)
	return tag.ID, result.Error
}

func (repository *TagRepository) Find() ([]*models.Tag, error) {
	var tags []*models.Tag
	if err := repository.db.Find(&tags).Error; err != nil {
		return nil, err
	}
	return tags, nil
}

func (repository *TagRepository) FindByName(name string) ([]*models.Tag, error) {
	var tags []*models.Tag
	if err := repository.db.Where("name ILIKE ?", "%"+name+"%").Find(&tags).Error; err != nil {
		return nil, err
	}
	return tags, nil
}

func (repository *TagRepository) Delete(id string) error {
	if err := repository.db.Delete(&models.Tag{}, id).Error; err != nil {
		return err
	}
	return nil
}
