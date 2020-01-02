package zaplog

import (
	"go.uber.org/zap"
	"log"
)

func InitLogger2() *zap.Logger {
	// 默认是Info级别
	logcfg := zap.NewProductionConfig()
	// 关闭日志采样采样
	logcfg.Sampling = nil
	//logcfg.Encoding = "json"
	logcfg.Encoding = "console"

	logger, err := logcfg.Build()
	if err != nil {
		log.Println("error", err)
	}

	return logger
}
