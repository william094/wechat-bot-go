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
	application.LoadConfig("../configs", "config-dev", &Conf)
	// 初始化 logger
	Log = logger.InitLogger(Conf.LogPath)
	// 初始化redis
	Rdb = db.InitRedis(&Conf.Redis)
	Log.Sugar().Info("init success")
}
