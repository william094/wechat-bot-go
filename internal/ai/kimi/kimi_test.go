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
				{Role: "system", Content: "你是 Kimi，由 Moonshot AI 提供的人工智能助手，你更擅长中文和英文的对话。你会为用户提供安全，有帮助，准确的回答。同时，你会拒绝一切涉及恐怖主义，种族歧视，黄色暴力等问题的回答。Moonshot AI 为专有名词，不可翻译成其他语言。"},
				{Role: "user", Content: "你好，我叫李雷，1+1等于多少？"},
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
			gotMsg, err := k.TextHandler(tt.args.ctx, tt.args.msgs)
			if (err != nil) != tt.wantErr {
				t.Errorf("Kimi.TextHandler() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Log(gotMsg)
		})
	}
}
