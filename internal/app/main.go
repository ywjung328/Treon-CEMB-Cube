package main

import (
	. "cube_config_handler"
	. "cube_data_handler"
	"encoding/json"
	"fmt"
	"global"
	"log"
	. "network_handler"
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
	InitZeroMQ()
}

func main() {
	defer global.LogFile.Close()
	defer global.Logger.Sync()
	defer global.Publisher.Close()
	defer global.Subscriber.Close()

	var err error

	global.Conf, err = ReadConfig("./conf.yaml")
	if err != nil {
		global.Logger.Panic(fmt.Sprintf("Reading config file failed: %v", err))
	}
	global.Logger.Info("Reading config file done.")

	scalarTicker := time.NewTicker(time.Duration(global.Conf.ScalarCycle) * time.Millisecond)
	vectorTicker := time.NewTicker(time.Duration(global.Conf.VectorCycle) * time.Millisecond)
	defer scalarTicker.Stop()
	defer vectorTicker.Stop()

	for {
		select {
		case <-scalarTicker.C:
			go scalarPrint()
		case <-vectorTicker.C:
			go vectorPrint()
		}
	}
	// go publish()
	// go subscribe()
}

func publish() {
	ticker := time.NewTicker(time.Duration(global.Conf.ScalarCycle) * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			for _, cube := range global.Conf.Cubes {
				realTimeMeasurements, err := GetRealTimeMeasurements(cube)
				if err != nil {
					global.Logger.Warn(fmt.Sprintf("Fetching realtime measurements fromr cube '%v' failed: %v", cube.Name, err))
					continue
				}
				data := TreonRTMConverter(realTimeMeasurements, cube)
				_, err = global.Publisher.Send(data, 0)
				if err != nil {
					global.Logger.Warn(fmt.Sprintf("Sending realtime measurements from cube '%v' via ZeroMQ failed: %v", cube.Name, err))
					continue
				}
			}
		}
	}
}

func subscribe() {
	for {
		message, err := global.Subscriber.Recv(0)
		if err != nil {
			global.Logger.Warn(fmt.Sprintf("ZeroMQ: Error while receiving message: %v", err))
			continue
		}
		global.Logger.Info(fmt.Sprintf("ZeroMQ: Received: %v", message))
	}
}

func scalarPrint() {
	// now := time.Now()
	// global.Logger.Info(fmt.Sprintf("timestamp: %v", now))
	for _, cube := range global.Conf.Cubes {
		realTimeMeasurements, err := GetRealTimeMeasurements(cube)
		if err != nil {
			global.Logger.Warn(fmt.Sprintf("Fetching realtime measurements from cube '%v' failed: %v", cube.Name, err))
			continue
		}
		// // fmt.Printf("Acc Max: %v / Vel Max: %v / Temperature: %v", realTimeMeasurements.AccMax, realTimeMeasurements.VelMax, realTimeMeasurements.T)
		// res, _ := json.Marshal(realTimeMeasurements)
		// global.Logger.Info(string(res))

		analogDigitalOutputs, err := GetAnalogDigitalOutputs(cube)
		if err != nil {
			global.Logger.Warn(fmt.Sprintf("Fetching analog digital outputs from cube '%v' failed: %v", cube.Name, err))
			continue
		}
		// // fmt.Printf("Acc Max: %v / Vel Max: %v / Temperature: %v", realTimeMeasurements.AccMax, realTimeMeasurements.VelMax, realTimeMeasurements.T)
		// res, _ = json.Marshal(analogDigitalOutputs)
		// global.Logger.Info(string(res))

		deviceStatuses, err := GetDeviceStatuses(cube)
		if err != nil {
			global.Logger.Warn(fmt.Sprintf("Fetching device statuses from cube '%v' failed: %v", cube.Name, err))
			continue
		}
		// // fmt.Printf("Acc Max: %v / Vel Max: %v / Temperature: %v", realTimeMeasurements.AccMax, realTimeMeasurements.VelMax, realTimeMeasurements.T)
		// res, _ = json.Marshal(deviceStatuses)
		// global.Logger.Info(string(res))

		channelStatuses, err := GetChannelStatuses(cube)
		if err != nil {
			global.Logger.Warn(fmt.Sprintf("Fetching channel statuses from cube '%v' failed: %v", cube.Name, err))
			continue
		}
		// // fmt.Printf("Acc Max: %v / Vel Max: %v / Temperature: %v", realTimeMeasurements.AccMax, realTimeMeasurements.VelMax, realTimeMeasurements.T)
		// res, _ = json.Marshal(channelStatuses)
		// global.Logger.Info(string(res))

		measurementsStatuses, err := GetMeasurementsStatuses(cube)
		if err != nil {
			global.Logger.Warn(fmt.Sprintf("Fetching measurements statuses from cube '%v' failed: %v", cube.Name, err))
			continue
		}
		// // fmt.Printf("Acc Max: %v / Vel Max: %v / Temperature: %v", realTimeMeasurements.AccMax, realTimeMeasurements.VelMax, realTimeMeasurements.T)
		// res, _ = json.Marshal(measurementsStatuses)
		// global.Logger.Info(string(res))
		data := make(map[string]interface{})
		data["RealTimeMeasurements"] = realTimeMeasurements
		data["AnalogDigitalOutputs"] = analogDigitalOutputs
		data["DeviceStatuses"] = deviceStatuses
		data["ChannelStatuses"] = channelStatuses
		data["MeasurementStatuses"] = measurementsStatuses
		res, _ := json.Marshal(data)
		global.Logger.Info(string(res))
	}
}

func vectorPrint() {
	for _, cube := range global.Conf.Cubes {
		vectorialMeasures, err := GetVectorialMeasures(cube)
		if err != nil {
			global.Logger.Warn(fmt.Sprintf("Fetching vectorial measures from cube '%v' failed: %v", cube.Name, err))
			continue
		}
		// fmt.Printf("Acc Max: %v / Vel Max: %v / Temperature: %v", realTimeMeasurements.AccMax, realTimeMeasurements.VelMax, realTimeMeasurements.T)
		data := make(map[string]interface{})
		data["VectorialMeasures"] = vectorialMeasures
		// res, _ := json.Marshal(vectorialMeasures)
		res, _ := json.Marshal(data)
		global.Logger.Info(string(res))
	}
}
