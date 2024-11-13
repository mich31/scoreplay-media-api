package services

import "github.com/mich31/scoreplay-media-api/repositories"

type TagService struct {
	repository repositories.TagRepository
}

func NewTagService(tagRepository repositories.TagRepository) *TagService {
	return &TagService{
		repository: tagRepository,
	}
}
