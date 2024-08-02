package message

import (
	"wehcat-bot-go/internal/ai"
	"wehcat-bot-go/internal/ai/kimi"
	"wehcat-bot-go/internal/config"
	"wehcat-bot-go/internal/data"

	"github.com/eatmoreapple/openwechat"
	"go.uber.org/zap"
)

type PrivateMsgHandler struct {
	conf *config.Config
	log  *zap.Logger
}

func NewPrivateMsgHandler(conf *config.Config, log *zap.Logger) *PrivateMsgHandler {
	return &PrivateMsgHandler{conf: conf, log: log}
}

// handle 处理消息
func (p *PrivateMsgHandler) Handle(msg *openwechat.Message) error {
	return p.ReceiveHandler(msg)
}

func (p *PrivateMsgHandler) ReceiveHandler(msg *openwechat.Message) error {
	if msg.IsText() {
		return p.ReceiveHandler(msg)
	}
	msg.ReplyText("很抱歉，我还在学习，暂不支持该消息")
	return nil
}

func (p *PrivateMsgHandler) ReplyHandler(msg *openwechat.Message) error {
	msgs := p.GetContext(msg)
	content, err := p.GetModelSV().TextHandler(msg.Context(), msgs)
	if err != nil {
		msg.ReplyText("我需要休息一下，请晚点再来问我吧")
	}
	_, err = msg.ReplyText(content.Content)
	if err != nil {
		p.log.Sugar().Errorf("reply message err \n", err)
	}
	return nil

}

func (p *PrivateMsgHandler) GetContext(msg *openwechat.Message) (contexts []ai.Message) {
	// TODO 获取上下文
	sender, err := msg.Sender()
	if err != nil || sender != nil {
		p.log.Sugar().Errorf("get msg sender err or sender is nil \n", err)
		msg.ReplyText("我需要休息一下啦，请晚点再来找我吧")
	}
	contexts, err = data.GetUserContext(sender.ID())
	if err != nil || contexts == nil {
		p.log.Sugar().Errorf("get user context error \n", err)
		contexts = make([]ai.Message, 0)
	}
	contexts = append(contexts, ai.Message{Role: "user", Content: msg.Content})
	if len(contexts) > 5 {
		contexts = contexts[2:]
	}
	return
}

func (p *PrivateMsgHandler) GetModelSV() ai.AiHandler {
	return &kimi.Kimi{ApiKey: p.conf.Kimi.ApiKey, BaseUrl: p.conf.Kimi.BaseUrl}
}
