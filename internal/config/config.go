package config

import (
	"wehcat-bot-go/internal/ai/kimi"
	"wehcat-bot-go/pkg/db"
)

type Config struct {
	Redis   db.RedisConf
	Kimi    kimi.Kimi
	LogPath string
}
