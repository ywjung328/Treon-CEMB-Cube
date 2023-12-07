package main

import (
	. "cube_config_handler"
	"fmt"
	"global"
	"log"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func init() {
	now := time.Now()
	err := os.MkdirAll("logs", os.ModePerm)
	if err != nil {
		panic(err)
	}
	global.LogFile, err = os.Create(fmt.Sprintf("logs/log_%v.log", now.Format(time.RFC3339)))
	if err != nil {
		panic(err)
	}
	log.SetOutput(global.LogFile)
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "timestamp"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.StacktraceKey = ""
	encoder := zapcore.NewConsoleEncoder(encoderConfig)

	fileCore := zapcore.NewCore(
		encoder,
		zapcore.AddSync(global.LogFile),
		zap.NewAtomicLevelAt(zap.InfoLevel),
	)

	global.Logger = zap.New(fileCore)
}

func main() {
	defer global.LogFile.Close()
	defer global.Logger.Sync()

	conf, err := ReadConfig("./conf.yaml")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(conf)
	global.Logger.Info("WOW")
}
