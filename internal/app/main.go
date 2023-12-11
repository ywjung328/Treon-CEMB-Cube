package main

import (
	. "cube_config_handler"
	. "cube_modbus_handler"
	"encoding/json"
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
	// cubes := global.Conf.Cubes

	ticker := time.NewTicker(time.Duration(global.Conf.ScalarCycle) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			go testPrint()
		}
	}
}

func testPrint() {
	now := time.Now()
	global.Logger.Info(fmt.Sprintf("timestamp: %v", now))
	for _, cube := range global.Conf.Cubes {
		realTimeMeasurements, err := GetRealTimeMeasurements(cube)
		if err != nil {
			global.Logger.Warn(fmt.Sprintf("Fetching realtime measurements from cube '%v' failed: %v", cube.Name, err))
			continue
		}
		// fmt.Printf("Acc Max: %v / Vel Max: %v / Temperature: %v", realTimeMeasurements.AccMax, realTimeMeasurements.VelMax, realTimeMeasurements.T)
		res, _ := json.Marshal(realTimeMeasurements)
		global.Logger.Info(string(res))

		/*
			analogDigitalOutputs, err := GetAnalogDigitalOutputs(cube)
			if err != nil {
				global.Logger.Warn(fmt.Sprintf("Fetching analog digital outputs from cube '%v' failed: %v", cube.Name, err))
				continue
			}
			// fmt.Printf("Acc Max: %v / Vel Max: %v / Temperature: %v", realTimeMeasurements.AccMax, realTimeMeasurements.VelMax, realTimeMeasurements.T)
			res, _ = json.Marshal(analogDigitalOutputs)
			global.Logger.Info(string(res))

			deviceStatuses, err := GetDeviceStatuses(cube)
			if err != nil {
				global.Logger.Warn(fmt.Sprintf("Fetching device statuses from cube '%v' failed: %v", cube.Name, err))
				continue
			}
			// fmt.Printf("Acc Max: %v / Vel Max: %v / Temperature: %v", realTimeMeasurements.AccMax, realTimeMeasurements.VelMax, realTimeMeasurements.T)
			res, _ = json.Marshal(deviceStatuses)
			global.Logger.Info(string(res))

			channelStatuses, err := GetChannelStatuses(cube)
			if err != nil {
				global.Logger.Warn(fmt.Sprintf("Fetching channel statuses from cube '%v' failed: %v", cube.Name, err))
				continue
			}
			// fmt.Printf("Acc Max: %v / Vel Max: %v / Temperature: %v", realTimeMeasurements.AccMax, realTimeMeasurements.VelMax, realTimeMeasurements.T)
			res, _ = json.Marshal(channelStatuses)
			global.Logger.Info(string(res))

			measurementsStatuses, err := GetMeasurementsStatuses(cube)
			if err != nil {
				global.Logger.Warn(fmt.Sprintf("Fetching measurements statuses from cube '%v' failed: %v", cube.Name, err))
				continue
			}
			// fmt.Printf("Acc Max: %v / Vel Max: %v / Temperature: %v", realTimeMeasurements.AccMax, realTimeMeasurements.VelMax, realTimeMeasurements.T)
			res, _ = json.Marshal(measurementsStatuses)
			global.Logger.Info(string(res))
		*/
	}
}
