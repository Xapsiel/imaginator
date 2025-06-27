package repository

import (
	"imageBot/internal/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	Image
	Poll
}

type Image interface {
	SaveImage(image *model.Image) error
	GetImage(delta int) ([]model.Image, error)
	SaveImageMessage(message_id int, photo_id string) error
}
type Poll interface {
	SavePoll(poll_id string, message_id int, poll_type string) error
	GetPoll(pollType string, delta int) (*model.Poll, error)
	Vote(user_id int64, poll_id string, answer_id int) error
	GetPollResults(poll *model.Poll) (*[]model.PollResult, error)
}

func NewRepository(db *pgxpool.Pool) Repository {
	return Repository{
		Image: NewImageRepository(db),
		Poll:  NewPollRepository(db),
	}
}
