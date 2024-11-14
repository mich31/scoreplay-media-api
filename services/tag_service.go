package services

import (
	"github.com/mich31/scoreplay-media-api/models"
	"github.com/mich31/scoreplay-media-api/repositories"
)

type TagService struct {
	repository repositories.ITagRepository
}

func NewTagService(tagRepository repositories.ITagRepository) *TagService {
	return &TagService{
		repository: tagRepository,
	}
}

func (service *TagService) GetTags(name string) ([]*models.Tag, error) {
	var tags []*models.Tag
	var err error
	if name != "" {
		tags, err = service.repository.FindByName(name)
	} else {
		tags, err = service.repository.Find()
	}

	if err != nil {
		return nil, err
	}
	return tags, nil
}

func (service *TagService) CreateTag(tag *models.Tag) (uint, error) {
	id, err := service.repository.Create(tag)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (service *TagService) DeleteTag(id string) error {
	err := service.repository.Delete(id)
	if err != nil {
		return err
	}
	return nil
}
