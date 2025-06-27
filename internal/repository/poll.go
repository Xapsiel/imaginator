package repository

import (
	"context"
	"fmt"

	"imageBot/internal/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

type pollRepository struct {
	db *pgxpool.Pool
}

func NewPollRepository(db *pgxpool.Pool) *pollRepository {
	return &pollRepository{db: db}

}
func (r *pollRepository) SavePoll(poll_id string, message_id int, poll_type string) error {
	query := `
			INSERT INTO polls(id,message_id,poll_type)
			VALUES ($1,$2)
			`
	_, err := r.db.Exec(context.Background(), query, poll_id, message_id, poll_type)
	if err != nil {
		return err
	}
	return nil
}
func (r *pollRepository) GetPoll(pollType string, delta int) (*model.Poll, error) {
	query := `
			SELECT id,message_id,poll_type FROM polls
			WHERE poll_type = $1 AND created_at > NOW() - ($2* interval 'hours' )
			`
	rows, err := r.db.Query(context.Background(), query, pollType, delta)
	if err != nil {
		return nil, err

	}
	defer rows.Close()
	polls := make([]model.Poll, 0)
	for rows.Next() {
		poll := model.Poll{}
		if err := rows.Scan(&poll.Id, &poll.MessageId, &poll.Poll_type); err != nil {
			return nil, err
		}
		polls = append(polls, poll)
	}
	if len(polls) != 1 {
		return nil, fmt.Errorf("expected 1 poll, got %d", len(polls))
	}
	return &polls[0], nil

}
func (r *pollRepository) Vote(user_id int64, poll_id string, answer_id int) error {
	query := `
			INSERT INTO poll_votes(user_id,poll_id,answer_id)
			VALUES ($1,$2,$3)
			ON CONFLICT DO NOTHING
			`
	_, err := r.db.Exec(context.Background(), query, user_id, poll_id, answer_id)
	if err != nil {
		return err
	}
	return nil
}

func (r *pollRepository) GetPollResults(poll *model.Poll) (*[]model.PollResult, error) {
	query := `
			SELECT user_id,poll_id,answer_id 
			FROM poll_votes
			WHERE poll_id = $1
			`
	rows, err := r.db.Query(context.Background(), query, poll.Id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	results := make([]model.PollResult, 0)
	for rows.Next() {
		result := model.PollResult{}
		if err := rows.Scan(&result.UserId, &result.PollId, &result.AnswerId); err != nil {
			return nil, err
		}
		results = append(results, result)
	}
	return &results, nil
}
