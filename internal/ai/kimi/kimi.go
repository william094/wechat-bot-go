package kimi

import (
	"context"
	"fmt"
	"wehcat-bot-go/internal/ai"

	"github.com/go-resty/resty/v2"
)

type Kimi struct {
	ApiKey  string
	BaseUrl string
}

type KimiRequest struct {
	Messages    []ai.Message
	Model       string
	Temperature float64
	Stream      bool
}

type KimiResp struct {
	ID      string   `json:"id"`
	Object  string   `json:"object"`
	Created int64    `json:"created"`
	Model   string   `json:"model"`
	Choices []Choice `json:"choices"`
	Usage   Usage    `json:"usage"`
}

type Choice struct {
	Index        int64      `json:"index"`
	Message      ai.Message `json:"message"`
	FinishReason string     `json:"finish_reason"`
}

type Usage struct {
	PromptTokens     int64 `json:"prompt_tokens"`
	CompletionTokens int64 `json:"completion_tokens"`
	TotalTokens      int64 `json:"total_tokens"`
}

func (k *Kimi) TextHandler(ctx context.Context, msgs []ai.Message) (msg ai.Message, err error) {
	results := &KimiResp{}
	body := &KimiRequest{Model: "moonshot-v1-8k", Temperature: 0.3, Messages: msgs}
	client := resty.New()
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", fmt.Sprintf("Bearer %s", k.ApiKey)).
		SetBody(body).SetResult(results).
		Post(fmt.Sprintf("%s/v1/chat/completions", k.BaseUrl))
	if err != nil {
		return
	}
	if resp.IsError() {

	}
	if len(results.Choices) > 0 {
		msg = results.Choices[0].Message
	}
	return
}
