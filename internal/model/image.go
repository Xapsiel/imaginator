package model

type Image struct {
	ID        string `json:"id"`
	MessageId int    `json:"message_id"`
	Text      string `json:"text"`
	Content   []byte `json:"content"`
	CreatedAt string `json:"created_at"`
	Status    string `json:"status"`
	Prompt    string `json:"prompt"`
}
