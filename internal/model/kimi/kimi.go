package kimi

import (
	"context"
	"encoding/json"
	"fmt"
	ai "wehcat-bot-go/internal/model"

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
	Index int64    `json:"index"`
	ID    string   `json:"id"`
	Type  string   `json:"type"`
	Func  Function `json:"function"`
}

type Function struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"`
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
	fmt.Println(resp)
	fmt.Println("==========================")
	if err != nil {
		panic(err)
	}

	if len(results.Choices) == 0 {
		return
	}
	if results.Choices[0].FinishReason == "tool_calls" {
		if results.Choices[0].Message.Content != "" {
			msgs = append(msgs, map[string]interface{}{
				"role":    results.Choices[0].Message.Role,
				"content": results.Choices[0].Message.Content,
			})
		}
		func_name := results.Choices[0].Message.ToolCalls[0].Func.Name
		call_arguments := make(map[string]string, 0)
		err = json.Unmarshal([]byte(results.Choices[0].Message.ToolCalls[0].Func.Arguments), &call_arguments)
		var tool_result []byte
		if func_name == "search" {
			tool_result = k.search(call_arguments["query"])
		} else {
			tool_result = k.get(call_arguments["url"])
		}
		msgs = append(msgs, map[string]interface{}{
			"role":         "tool",
			"content":      string(tool_result),
			"tool_call_id": results.Choices[0].Message.ToolCalls[0].ID,
			"name":         results.Choices[0].Message.ToolCalls[0].Func.Name,
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

type SerachResult struct {
	Results []struct {
		Body  string `json:"body"`
		Href  string `json:"href"`
		Title string `json:"title"`
	}
}

func (k *Kimi) search(query string) []byte {
	resp, err := resty.New().R().
		SetQueryParams(map[string]string{"q": query, "max_resluts": "1"}).
		Get("http://localhost:8000/search")
	if err != nil {
		panic(err)
	}
	return resp.Body()
}

func (k *Kimi) get(url string) []byte {
	resp, err := resty.New().R().Get(url)
	if err != nil {
		panic(err)
	}
	return resp.Body()
}
