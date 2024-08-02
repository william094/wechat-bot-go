package handlers

import (
	"github.com/eatmoreapple/openwechat"
)

type MsgHandler interface {
	ReceiveHandler(*openwechat.Message) error
	ReplyHandler(*openwechat.Message) error
}
