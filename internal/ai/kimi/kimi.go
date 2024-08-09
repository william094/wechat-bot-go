package kimi

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
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
	Tools       interface{}
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
	Index        int64   `json:"index"`
	Message      KimiMsg `json:"message"`
	FinishReason string  `json:"finish_reason"`
}

type Usage struct {
	PromptTokens     int64 `json:"prompt_tokens"`
	CompletionTokens int64 `json:"completion_tokens"`
	TotalTokens      int64 `json:"total_tokens"`
}

type KimiMsg struct {
	Role      string     `json:"role"`
	Content   string     `json:"content"`
	ToolCalls []ToolCall `json:"tool_calls"`
}

type ToolCall struct {
	Index int64             `json:"index"`
	ID    string            `json:"id"`
	Type  string            `json:"type"`
	Func  map[string]string `json:"function"`
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
		msg = ai.Message{
			Role:    results.Choices[0].Message.Role,
			Content: results.Choices[0].Message.Content,
		}
	}
	return
}

func (k *Kimi) ToolCalls(ctx context.Context, msgs []map[string]interface{}) {
	results := &KimiResp{}
	//body := &KimiRequest{Model: "moonshot-v1-8k", Temperature: 0.3, Messages: msgs, Tools: tools}
	body := map[string]interface{}{
		"model":       "moonshot-v1-8k",
		"temperature": 0.3,
		"tools":       tools,
		"Messages":    msgs,
	}
	client := resty.New()
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", fmt.Sprintf("Bearer %s", k.ApiKey)).
		SetBody(body).SetResult(results).
		Post(fmt.Sprintf("%s/v1/chat/completions", k.BaseUrl))
	if err != nil {
		return
	}
	bytes, _ := json.Marshal(body)
	fmt.Println(string(bytes))
	fmt.Println(resp)
	fmt.Println("==========================")
	if len(results.Choices) == 0 {
		return
	}
	if results.Choices[0].FinishReason == "tool_calls" {
		//if results.Choices[0].Message.Content != "" {
		msgs = append(msgs, map[string]interface{}{
			"role":    results.Choices[0].Message.Role,
			"content": results.Choices[0].Message.Content,
		})
		//}
		msgs = append(msgs, map[string]interface{}{
			"role":         "tool",
			"content":      strings.Replace(strings.Replace(results.Choices[0].Message.ToolCalls[0].Func["arguments"], "\n", "", -1), "\\", "", -1),
			"tool_call_id": results.Choices[0].Message.ToolCalls[0].ID,
			"name":         results.Choices[0].Message.ToolCalls[0].Func["name"],
		})

		k.ToolCalls(ctx, msgs)
	}
	fmt.Println(resp)
	fmt.Println("==========================")
}

var (
	tools = []map[string]interface{}{
		{
			"type": "function",
			"function": map[string]interface{}{
				"name":        "search",
				"description": " 当你无法回答用户的问题时，或用户查询的是实时数据，或用户请求你进行联网搜索时，调用此工具。请从与用户的对话中提取用户想要搜索的内容作为 query 参数的值。搜索结果包含网站的标题、网站的地址（URL）以及网站简介。",
				"parameters": map[string]interface{}{
					"type":     "object",
					"required": []string{"query"},
					"properties": map[string]interface{}{
						"query": map[string]interface{}{
							"type":        "string",
							"description": "用户搜索的内容，请从用户的提问或聊天上下文中提取。",
						},
					},
				},
			},
		}, {
			"type": "function",
			"function": map[string]interface{}{
				"name":        "crawl",
				"description": "根据网站地址（URL）获取网页内容。",
				"parameters": map[string]interface{}{
					"type":     "object",
					"required": []string{"url"},
					"properties": map[string]interface{}{
						"url": map[string]interface{}{
							"type":        "string",
							"description": "需要获取内容的网站地址（URL），通常情况下从搜索结果中可以获取网站的地址。",
						},
					},
				},
			},
		},
	}
)
