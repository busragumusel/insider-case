package model

type Response struct {
	MessageID string `json:"messageId"`
	Message   string `json:"message"`
}

type Payload struct {
	To      string `json:"to"`
	Content string `json:"content"`
}
