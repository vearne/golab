package main

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"time"
	"github.com/vearne/golab/myencoder"
)


func MyEncoding(config zapcore.EncoderConfig) (zapcore.Encoder, error) {
	return myencoder.NewConsole2Encoder(config, true), nil
}


func main() {
	// 注册一个Encoder
	zap.RegisterEncoder("console2", MyEncoding)

	// 默认是Info级别
	logcfg := zap.NewProductionConfig()
	// 启用自定义的Encoding
	logcfg.Encoding = "console2"

	logger, err := logcfg.Build()
	if err != nil {
		fmt.Println("err", err)
	}

	defer logger.Sync()
	for i := 0; i < 3; i++ {
		time.Sleep(1 * time.Second)
		logger.Info("some message", zap.String("name", "buick2008"),
			zap.Int("age", 15))
	}
}
