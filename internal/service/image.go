package service

import (
	"imageBot/internal/api"
	"imageBot/internal/model"
	"imageBot/internal/repository"
)

type imageService struct {
	api  *api.Text2ImageAPI
	repo repository.Image
}

func NewImageService(repo repository.Image, api *api.Text2ImageAPI) *imageService {
	return &imageService{repo: repo, api: api}
}

func (s *imageService) SaveImage(image *model.Image) error {
	return s.repo.SaveImage(image)
}

func (s *imageService) GenerateImage(prompt string, w, h int) (model.Image, error) {
	image, err := s.api.Draw(prompt, w, h)
	if err != nil {
		return model.Image{}, err
	}
	img := model.Image{
		Content: image,
		Prompt:  prompt,
	}
	return img, nil
}
func (s *imageService) GetImage(delta int) ([]model.Image, error) {
	return s.repo.GetImage(delta)
}
func (s *imageService) SaveImageMessage(message_id int, photo_id string) error {
	return s.repo.SaveImageMessage(message_id, photo_id)
}
