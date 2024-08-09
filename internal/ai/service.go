package ai

import "context"

type Message struct {
	Role       string `json:"role"`
	Content    string `json:"content"`
	ToolCallId string `json:"tool_call_id"`
	Name       string `json:"name"`
	ToolChoice string `json:"tool_choice"`
}

type AiHandler interface {
	TextHandler(ctx context.Context, msgs []Message) (msg Message, err error)
}
