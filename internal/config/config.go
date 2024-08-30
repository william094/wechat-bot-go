package config

import (
	"wehcat-bot-go/internal/model/doubao"
	"wehcat-bot-go/internal/model/kimi"
	"wehcat-bot-go/pkg/db"
)

type Config struct {
	Redis   db.RedisConf
	Kimi    kimi.Kimi
	Doubao  doubao.Doubao
	LogPath string
}
