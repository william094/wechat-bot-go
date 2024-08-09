package handlers

import (
	"github.com/eatmoreapple/openwechat"
)

type MsgHandler interface {
	Dispatch(*openwechat.Message) error
	ReceiveHandler(*openwechat.Message) error
	ReplyHandler(*openwechat.Message) error
}
