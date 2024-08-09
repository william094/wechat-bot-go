package kimi

import (
	"context"
	"testing"
	"wehcat-bot-go/internal/ai"
)

func TestKimi_TextHandler(t *testing.T) {
	type fields struct {
		ApiKey  string
		BaseUrl string
	}
	type args struct {
		ctx  context.Context
		msgs []ai.Message
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantMsg ai.Message
		wantErr bool
	}{
		{
			name:   "success",
			fields: fields{ApiKey: "sk-MyQJaP6HFPl50uk8d976frV9bNKHrsUAbFLq4C0RBVEMFiD6", BaseUrl: "https://api.moonshot.cn"},
			args: args{ctx: context.Background(), msgs: []ai.Message{
				//{Role: "system", Content: "你是 Kimi，由 Moonshot AI 提供的人工智能助手，你更擅长中文和英文的对话。你会为用户提供安全，有帮助，准确的回答。同时，你会拒绝一切涉及恐怖主义，种族歧视，黄色暴力等问题的回答。Moonshot AI 为专有名词，不可翻译成其他语言。"},
				{Role: "user", Content: "2024奥运会中国获奖情况"},
			}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := &Kimi{
				ApiKey:  tt.fields.ApiKey,
				BaseUrl: tt.fields.BaseUrl,
			}
			k.TextHandler(tt.args.ctx, tt.args.msgs)

		})
	}
}

func TestKimi_ToolCallHandler(t *testing.T) {
	type fields struct {
		ApiKey  string
		BaseUrl string
	}
	type args struct {
		ctx  context.Context
		msgs []map[string]interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantMsg ai.Message
		wantErr bool
	}{
		{
			name:   "success",
			fields: fields{ApiKey: "sk-MyQJaP6HFPl50uk8d976frV9bNKHrsUAbFLq4C0RBVEMFiD6", BaseUrl: "https://api.moonshot.cn"},
			args: args{ctx: context.Background(), msgs: []map[string]interface{}{
				//{Role: "system", Content: "你是 Kimi，由 Moonshot AI 提供的人工智能助手，你更擅长中文和英文的对话。你会为用户提供安全，有帮助，准确的回答。同时，你会拒绝一切涉及恐怖主义，种族歧视，黄色暴力等问题的回答。Moonshot AI 为专有名词，不可翻译成其他语言。"},
				{"role": "user", "content": "今日上海静安区天气情况"},
			}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := &Kimi{
				ApiKey:  tt.fields.ApiKey,
				BaseUrl: tt.fields.BaseUrl,
			}
			k.ToolCalls(tt.args.ctx, tt.args.msgs)

		})
	}
}
