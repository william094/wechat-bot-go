package handlers

import (
	"log"
	"wehcat-bot-go/internal/app"
	"wehcat-bot-go/internal/wechat/message"

	"github.com/eatmoreapple/openwechat"
)

// Handler 全局处理入口
func MessageHandler(msg *openwechat.Message) {
	log.Printf("hadler Received msg : %v", msg.Content)
	// 处理群消息
	if msg.IsSendByGroup() {

		return
	}

	//处理私聊消息
	message.NewPrivateMsgHandler(app.Conf, app.Log).Handle(msg)
}

func MessageErrorHandler(err error) error {
	//获取消息发生错误的handle, 返回err == nil 则尝试继续监听
	if err != nil {
		app.Log.Sugar().Errorf("get msg err", err)
	}
	return nil
}
