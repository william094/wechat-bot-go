package cmd

import (
	"wehcat-bot-go/internal/app"
	"wehcat-bot-go/internal/wechat/handlers"

	"github.com/eatmoreapple/openwechat"
	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "run wechat bot ",
	Long:  `run wechat bot `,
	Run: func(cmd *cobra.Command, args []string) {
		app.InitConfig(cfgFile)
		Run()
	},
}

func Run() {
	// 桌面模式
	bot := openwechat.DefaultBot(openwechat.Desktop)
	// 注册消息处理函数
	bot.MessageHandler = handlers.MessageHandler
	// 注册消息获取错误处理函数
	bot.MessageErrorHandler = handlers.MessageErrorHandler
	// 注册登陆二维码回调
	bot.UUIDCallback = openwechat.PrintlnQrcodeUrl
	// 创建热存储容器对象
	reloadStorage := openwechat.NewFileHotReloadStorage("storage.json")
	// 执行热登录
	err := bot.HotLogin(reloadStorage)
	if err != nil {
		if err = bot.Login(); err != nil {
			app.Log.Sugar().Errorf("bot start failed err: %v", err)
			return
		}
	}
	// 阻塞主goroutine, 直到发生异常或者用户主动退出
	bot.Block()
}

func init() {
	RootCmd.AddCommand(runCmd)
}
