package handlers

import (
	"log"
	"wehcat-bot-go/internal/app"
	"wehcat-bot-go/internal/wechat/message"

	"github.com/eatmoreapple/openwechat"
)

var (
	groupHandler   = message.NewGroupMsgHandler(app.Conf, app.Log)
	privateHandler = message.NewPrivateMsgHandler(app.Conf, app.Log)
)

// Handler 全局处理入口
func MessageHandler(msg *openwechat.Message) {
	log.Printf("hadler Received msg : %v", msg.Content)
	// 处理群消息
	if msg.IsSendByGroup() {
		groupHandler.ReceiveHandler(msg)

	} else {
		privateHandler.ReceiveHandler(msg)
	}
}

func MessageErrorHandler(err error) error {
	//获取消息发生错误的handle, 返回err == nil 则尝试继续监听
	if err != nil {
		app.Log.Sugar().Errorf("get msg err", err)
	}
	return nil
}
