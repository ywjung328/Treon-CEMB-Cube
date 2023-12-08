package main

import (
	. "cube_config_handler"
	. "cube_modbus_handler"
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

	var err error

	global.Conf, err = ReadConfig("./conf.yaml")
	if err != nil {
		global.Logger.Panic(fmt.Sprintf("Reading config file failed: %v", err))
	}
	global.Logger.Info("Reading config file done.")
	cubes := global.Conf.Cubes

	for _, cube := range cubes {
		realTimeMeasurements, err := GetDeviceStatuses(cube)
		if err != nil {
			global.Logger.Warn(fmt.Sprintf("Fetching realtime measurements from cube '%v' failed: %v", cube.Name, err))
			continue
		}
		fmt.Println(realTimeMeasurements)
	}
}
