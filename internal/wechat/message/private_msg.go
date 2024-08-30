package message

import (
	"wehcat-bot-go/internal/config"
	"wehcat-bot-go/internal/data"
	"wehcat-bot-go/internal/model"
	"wehcat-bot-go/internal/model/doubao"

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

func (p *PrivateMsgHandler) ReceiveHandler(msg *openwechat.Message) error {
	p.log.Sugar().Infof("收到私聊消息---消息类型：%s,消息子类型：%s,消息发送者：%s", msg.MsgType, msg.SubMsgType, msg.Content)
	if msg.IsText() {
		return p.ReplyHandler(msg)
	}
	msg.ReplyText("很抱歉，我还在学习，暂时只支持文本消息")
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
		p.log.Sugar().Errorf("reply message err \n %+v", err)
	}
	return nil

}

func (p *PrivateMsgHandler) GetContext(msg *openwechat.Message) (contexts []model.Message) {
	sender, err := msg.Sender()
	if err != nil || sender == nil {
		p.log.Sugar().Errorf("get msg sender err or sender is nil \n %+v", err)
		msg.ReplyText("我需要休息一下啦，请晚点再来找我吧")
		return
	}
	contexts, err = data.GetUserContext(sender.AvatarID())
	if err != nil || contexts == nil {
		p.log.Sugar().Errorf("get user context error \n %+v", err)
		contexts = make([]model.Message, 0)
	}
	contexts = append(contexts, model.Message{Role: "user", Content: msg.Content})
	p.log.Sugar().Infof("用户-%s,上下文信息：%v \n", sender.AvatarID(), contexts)
	if len(contexts) > 5 {
		contexts = contexts[2:]
	}
	err = data.SetUserContext(sender.AvatarID(), contexts)
	if err != nil {
		p.log.Sugar().Errorf("用户-%s上下文缓存异常 \n %+v", sender.AvatarID(), err)
	}
	return
}

func (p *PrivateMsgHandler) GetModelSV() model.AiHandler {
	// return &kimi.Kimi{ApiKey: p.conf.Kimi.ApiKey, BaseUrl: p.conf.Kimi.BaseUrl}
	return doubao.NewDoubao(&p.conf.Doubao, p.log)
}
