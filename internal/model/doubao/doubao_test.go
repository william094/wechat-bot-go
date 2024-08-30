package doubao

import (
	"context"
	"os"
	"reflect"
	"testing"
	"wehcat-bot-go/internal/model"

	"github.com/go-rod/rod"
	"github.com/russross/blackfriday/v2"
	"github.com/volcengine/volcengine-go-sdk/service/arkruntime"
	doubao "github.com/volcengine/volcengine-go-sdk/service/arkruntime/model"
)

func TestDoubao_TextHandler(t *testing.T) {
	type fields struct {
		ApiKey string
		BotId  string
		Client *arkruntime.Client
	}
	type args struct {
		ctx  context.Context
		msgs []model.Message
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantMsg model.Message
		wantErr bool
	}{
		{
			name: "success",
			fields: fields{
				ApiKey: "",
				BotId:  "",
			},
			args: args{
				ctx: context.Background(),
				msgs: []model.Message{
					{Role: doubao.ChatMessageRoleUser, Content: "上海最近七天天气怎么样"},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Doubao{
				ApiKey: tt.fields.ApiKey,
				BotId:  tt.fields.BotId,
			}
			gotMsg, err := d.TextHandler(tt.args.ctx, tt.args.msgs)
			if (err != nil) != tt.wantErr {
				t.Errorf("Doubao.TextHandler() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotMsg, tt.wantMsg) {
				t.Errorf("Doubao.TextHandler() = %v, want %v", gotMsg, tt.wantMsg)
			}
			markdownProcess(gotMsg.Content)
		})
	}
}

func markdownProcess(text string) {
	// 使用blackfriday解析Markdown
	html := blackfriday.Run([]byte(text))
	// 添加HTML结构和字符编码
	htmls := `<!DOCTYPE html>
    <html>
    <head>
        <meta charset="UTF-8">
        <style>
            body { font-family: Arial, sans-serif; }
        </style>
    </head>
    <body>` + string(html) + `</body></html>`
	// 使用Rod浏览器将HTML渲染为图像
	page := rod.New().MustConnect().MustPage("data:text/html," + string(htmls))
	page.MustWaitLoad()

	// 截图并保存为图像
	imgBytes := page.MustScreenshot()
	os.WriteFile("output.png", imgBytes, 0644)
}
