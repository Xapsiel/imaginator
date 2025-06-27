package service

import (
	"imageBot/internal/api"
	"imageBot/internal/model"
	"imageBot/internal/repository"
)

type Service struct {
	Image
	Poll
}
type Image interface {
	SaveImage(image *model.Image) error
	GenerateImage(prompt string, w, h int) (model.Image, error)
	GetImage(delta int) ([]model.Image, error)
	SaveImageMessage(message_id int, photo_id string) error
}
type Poll interface {
	SavePoll(poll_id string, message_id int, poll_type string) error
	GetPoll(poll_type string, delta int) (poll *model.Poll, err error)
	Vote(user_id int64, poll_id string, answer_id int) error
	GetPollResults(poll *model.Poll) (answer_id int, count int, err error)
}

func NewService(repo repository.Repository, api_fb *api.Text2ImageAPI) *Service {
	return &Service{
		Image: NewImageService(repo.Image, api_fb),
		Poll:  NewPollService(repo.Poll),
	}
}
