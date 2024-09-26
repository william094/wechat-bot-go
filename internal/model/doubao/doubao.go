package doubao

import (
	"context"
	"fmt"
	"sync"
	ai "wehcat-bot-go/internal/model"

	"github.com/volcengine/volcengine-go-sdk/service/arkruntime"
	"github.com/volcengine/volcengine-go-sdk/service/arkruntime/model"
	"github.com/volcengine/volcengine-go-sdk/volcengine"
	"go.uber.org/zap"
)

var one sync.Once

type Doubao struct {
	ApiKey string
	BotId  string
	Log    *zap.Logger
	Client *arkruntime.Client
}

func NewDoubao(d *Doubao, log *zap.Logger) *Doubao {
	d.Log = log
	d.GetAskClient()
	return d
}

func (d *Doubao) GetAskClient() {
	if d.Client == nil {
		one.Do(func() {
			d.Client = arkruntime.NewClientWithApiKey(d.ApiKey)
		})
	}
}

func (d *Doubao) TextHandler(ctx context.Context, msgs []ai.Message) (msg ai.Message, err error) {
	req := model.BotChatCompletionRequest{
		BotId:    d.BotId,
		Messages: []*model.ChatCompletionMessage{},
	}
	for _, v := range msgs {
		req.Messages = append(req.Messages, &model.ChatCompletionMessage{
			Role: model.ChatMessageRoleUser,
			Content: &model.ChatCompletionMessageContent{
				StringValue: volcengine.String(v.Content),
			},
		})
	}
	resp, err := d.Client.CreateBotChatCompletion(ctx, req)
	if err != nil {
		fmt.Printf("standard chat error: %v\n", err)
		return
	}
	fmt.Println(*resp.Choices[0].Message.Content.StringValue)
	if resp.References != nil {
		for _, ref := range resp.References {
			fmt.Printf("reference url: %s\n", ref.Url)
		}
	}
	msg = ai.Message{
		Content: *resp.Choices[0].Message.Content.StringValue,
		Role:    resp.Choices[0].Message.Role,
	}
	return
}
