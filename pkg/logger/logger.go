package logger

import (
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func InitLogger(logPath string) *zap.Logger {
	writeSyncer := getLogWriter(logPath)
	encoder := getEncoder()
	core := zapcore.NewCore(encoder, writeSyncer, zapcore.DebugLevel)
	return zap.New(core, zap.AddCaller())
}

// 日志编码
func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05.000")
	return zapcore.NewConsoleEncoder(encoderConfig)
}

// 日志文件输出、切割
func getLogWriter(logPath string) zapcore.WriteSyncer {
	// 使用lumberjack进行日志文件的轮转
	lumberjackLogger := &lumberjack.Logger{
		Filename:   logPath, // 日志文件的位置
		MaxSize:    100,     // 文件最大尺寸（以MB为单位）
		MaxBackups: 10,      // 保留的最大旧文件数量
		MaxAge:     7,       // 保留旧文件的最大天数
		Compress:   true,    // 是否压缩/归档旧文件
		LocalTime:  true,    // 使用本地时间创建时间戳
	}
	return zapcore.AddSync(lumberjackLogger)
}
