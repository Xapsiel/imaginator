package service

import (
	"imageBot/internal/model"
	"imageBot/internal/repository"
)

type pollService struct {
	repo repository.Poll
}

func NewPollService(repo repository.Poll) *pollService {
	return &pollService{repo: repo}
}
func (s *pollService) SavePoll(poll_id string, message_id int, poll_type string) error {
	return s.repo.SavePoll(poll_id, message_id, poll_type)
}

func (s *pollService) GetPoll(poll_type string, delta int) (poll *model.Poll, err error) {
	return s.repo.GetPoll(poll_type, delta)
}
func (s *pollService) Vote(user_id int64, poll_id string, answer_id int) error {
	return s.repo.Vote(user_id, poll_id, answer_id)
}
func (s *pollService) GetPollResults(poll *model.Poll) (answer_id, count int, err error) {
	res, err := s.repo.GetPollResults(poll)
	if err != nil {
		return 0, 0, err
	}
	var m map[string]int = make(map[string]int)
	maxAnswer := 0
	maxAnswerVote := 0
	for _, e := range *res {
		m[e.PollId]++
		if m[e.PollId] > maxAnswerVote {
			maxAnswer = e.AnswerId
			maxAnswerVote = m[e.PollId]
		}
	}
	return maxAnswer, maxAnswerVote, nil

}
