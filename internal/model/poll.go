package model

type Poll struct {
	Id        string `json:"id"`
	MessageId int    `json:"message_id"`
	Poll_type string `json:"poll_type"`
}

type PollResult struct {
	UserId   int    `json:"user_id"`
	PollId   string `json:"poll_id"`
	AnswerId int    `json:"answer_id"`
}
