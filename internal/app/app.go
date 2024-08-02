package app

import (
	"wehcat-bot-go/internal/config"
	application "wehcat-bot-go/pkg/config"
	"wehcat-bot-go/pkg/db"
	"wehcat-bot-go/pkg/logger"

	"github.com/go-redis/redis"
	"go.uber.org/zap"
)

var (
	Conf *config.Config
	Rdb  *redis.Client
	Log  *zap.Logger
)

func init() {
	// 记载配置文件
	application.LoadConfig("/Users/zhr/myself/wechat-bot-go/configs", "config", &Conf)
	// 初始化 logger
	logger.InitLogger(Conf.LogPath)
	// 初始化redis
	db.InitRedis(&Conf.Redis)
}
