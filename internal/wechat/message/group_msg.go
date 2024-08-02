package message

import "github.com/eatmoreapple/openwechat"

type GroupMsgHandler struct{}

func (p *GroupMsgHandler) ReceiveHandler(msg *openwechat.Message) error { return nil }

func (p *GroupMsgHandler) ReplyHandler(msg *openwechat.Message) error { return nil }
