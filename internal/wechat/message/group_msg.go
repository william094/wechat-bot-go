package message

import (
	"fmt"
	"strings"
	"wehcat-bot-go/internal/ai"
	"wehcat-bot-go/internal/ai/kimi"
	"wehcat-bot-go/internal/config"
	"wehcat-bot-go/internal/data"

	"github.com/eatmoreapple/openwechat"
	"go.uber.org/zap"
)

type GroupMsgHandler struct {
	conf *config.Config
	log  *zap.Logger
}

func NewGroupMsgHandler(conf *config.Config, log *zap.Logger) *GroupMsgHandler {
	return &GroupMsgHandler{conf: conf, log: log}
}

func (p *GroupMsgHandler) ReceiveHandler(msg *openwechat.Message) error {
	if !msg.IsAt() {
		return nil
	}
	return p.ReplyHandler(msg)
}

func (p *GroupMsgHandler) ReplyHandler(msg *openwechat.Message) error {
	//消息发送者
	sender, _ := msg.Sender()
	//转换获取群相关信息
	group := openwechat.Group{User: sender}
	msgs := p.GetContext(msg)
	message, err := p.GetModelSV().TextHandler(msg.Context(), msgs)
	if err != nil {
		p.log.Sugar().Errorf("get msg  sender is nil \n")
		msg.ReplyText("我需要休息一下啦，请晚点再来找我吧")
		return err
	}
	// 获取@我的用户
	groupSender, err := msg.SenderInGroup()
	if err != nil {
		p.log.Sugar().Error("get sender in group error \n", err)
		return err
	}
	// 回复并@发送者
	replyText := fmt.Sprintf("@%s %s", groupSender.NickName, strings.TrimSpace(message.Content))
	_, err = msg.ReplyText(replyText)
	if err != nil {
		p.log.Sugar().Error("group reply msg error: %v \n", err)
		return err
	}
	//缓存上下文信息
	msgs = append(msgs, message)
	if len(msgs) > 5 {
		msgs = msgs[2:]
	}
	err = data.SetUserContext(fmt.Sprintf("%d-%s", group.ChatRoomId, sender.ID()), msgs)
	if err != nil {
		p.log.Sugar().Errorf("群-%d;用户-%s;上下文缓存异常 \n", group.ChatRoomId, sender.ID(), err)
	}
	return nil
}

func (p *GroupMsgHandler) GetContext(msg *openwechat.Message) (msgs []ai.Message) {
	sender, _ := msg.Sender()
	group := openwechat.Group{User: sender}
	if sender == nil {
		p.log.Sugar().Errorf("get msg  sender is nil \n")
		msg.ReplyText("我需要休息一下啦，请晚点再来找我吧")
		return
	}
	content := strings.TrimSpace(strings.ReplaceAll(msg.Content, "@"+sender.Self().NickName, ""))
	msgs, err := data.GetUserContext(fmt.Sprintf("%s-%d", sender.ID(), group.ChatRoomId))
	if err != nil || msgs == nil {
		p.log.Sugar().Errorf("get user context error \n", err)
		msgs = make([]ai.Message, 0)
	}
	msgs = append(msgs, ai.Message{Role: "user", Content: content})
	p.log.Sugar().Infof("群-%d;用户%s;上下文信息：%v \n", sender.ID(), group.ChatRoomId, msgs)
	return
}

func (p *GroupMsgHandler) GetModelSV() ai.AiHandler {
	return &kimi.Kimi{ApiKey: p.conf.Kimi.ApiKey, BaseUrl: p.conf.Kimi.BaseUrl}
}
