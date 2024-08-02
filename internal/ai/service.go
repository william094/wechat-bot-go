package ai

import "context"

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type AiHandler interface {
	TextHandler(ctx context.Context, msgs []Message) (msg Message, err error)
}
