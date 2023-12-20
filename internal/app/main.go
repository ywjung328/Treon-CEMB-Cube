package main

import (
	. "cube_config_handler"
	. "cube_data_handler"
	"encoding/json"
	"fmt"
	"global"
	"log"

	// . "network_handler"
	"os"
	"time"

	"github.com/pebbe/zmq4"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func init() {
	err := os.MkdirAll("logs", os.ModePerm)
	if err != nil {
		panic(err)
	}
	// now := time.Now()
	// global.LogFile, err = os.Create(fmt.Sprintf("logs/log_%v.log", now.Format(time.RFC3339)))
	global.LogFile, err = os.Create("logs/log.log")
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
	// err = InitZeroMQ()
	if err != nil {
		global.Logger.Warn(fmt.Sprintf("ZeroMQ initiation failed: %v", err))
	} else {
		global.Logger.Info("ZeroMQ initiated successfully.")
	}
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

	// INITIATING ZMQ4
	context, err := zmq4.NewContext()
	if err != nil {
		global.Logger.Error(fmt.Sprintf("Creating zmq4 context failed: %v", err))
	}
	global.Publisher, err = context.NewSocket(zmq4.PUB)
	if err != nil {
		global.Logger.Error(fmt.Sprintf("Creating zmq4 publisher failed: %v", err))
	}
	global.Subscriber, err = context.NewSocket(zmq4.SUB)
	if err != nil {
		global.Logger.Error(fmt.Sprintf("Creating zmq4 subscriber failed: %v", err))
	}
	err = global.Publisher.Bind(fmt.Sprintf("tcp://*:%v", global.Conf.PublishPort))
	if err != nil {
		global.Logger.Error(fmt.Sprintf("Binding zmq4 publisher to tcp://*:%v failed: %v", global.Conf.PublishPort, err))
	}
	err = global.Subscriber.Connect(fmt.Sprintf("tcp://localhost:%v", global.Conf.SubscribePort))
	if err != nil {
		global.Logger.Error(fmt.Sprintf("Connecting zmq4 subscriber to tcp://localhost:%v failed: %v", global.Conf.SubscribePort, err))
	}
	err = global.Subscriber.SetSubscribe(global.Conf.Filter)
	if err != nil {
		global.Logger.Error(fmt.Sprintf("Setting subscribtion (filter: %v) failed: %v", global.Conf.Filter, err))
	}

	defer context.Term()
	defer global.Publisher.Close()
	defer global.Subscriber.Close()

	go subscribe()
	for {
		select {
		case <-scalarTicker.C:
			go scalarPublish()
		case <-vectorTicker.C:
			go vectorPublish()
		}
	}
}

/*
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
					fmt.Println(err)
					global.Logger.Warn(fmt.Sprintf("Sending realtime measurements from cube '%v' via ZeroMQ failed: %v", cube.Name, err))
					continue
				}
			}
		}
	}
}
*/

func subscribe() {
	for {
		message, err := global.Subscriber.Recv(0)
		if err != nil {
			global.Logger.Warn(fmt.Sprintf("ZeroMQ: Error while receiving message: %v", err))
			continue
		}
		global.Logger.Info(fmt.Sprintf("ZeroMQ: Received: %v", message))
		fmt.Println(fmt.Sprintf("ZeroMQ: Received: %v", message))
	}
}

func scalarPublish() {
	// now := time.Now()
	// global.Logger.Info(fmt.Sprintf("timestamp: %v", now))
	for _, cube := range global.Conf.Cubes {
		serialNumber, err := GetSerialNumber(cube)
		if err != nil {
			global.Logger.Warn(fmt.Sprintf("Fetching S/N from cube '%v' failed: %v", cube.Name, err))
			continue
		}
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
		data["Timestamp"] = time.Now().Unix()
		data["SerialNumber"] = serialNumber
		data["RealTimeMeasurements"] = realTimeMeasurements
		data["AnalogDigitalOutputs"] = analogDigitalOutputs
		data["DeviceStatuses"] = deviceStatuses
		data["ChannelStatuses"] = channelStatuses
		data["MeasurementStatuses"] = measurementsStatuses
		res, _ := json.Marshal(data)
		// global.Logger.Info(string(res))
		// _, err = global.Publisher.Send(string(res), 0)
		_, err = global.Publisher.Send(fmt.Sprintf("%v %v", global.Conf.Filter, string(res)), 0)
		if err != nil {
			fmt.Println(err)
			global.Logger.Warn(fmt.Sprintf("Sending realtime measurements from cube '%v' via ZeroMQ failed: %v", cube.Name, err))
			continue
		}
		global.Logger.Info(fmt.Sprintf("SCALAR PUBLISHED: %v", string(res)))
	}
}

func vectorPublish() {
	for _, cube := range global.Conf.Cubes {
		serialNumber, err := GetSerialNumber(cube)
		if err != nil {
			global.Logger.Warn(fmt.Sprintf("Fetching S/N from cube '%v' failed: %v", cube.Name, err))
			continue
		}
		vectorialMeasures, err := GetVectorialMeasures(cube)
		if err != nil {
			global.Logger.Warn(fmt.Sprintf("Fetching vectorial measures from cube '%v' failed: %v", cube.Name, err))
			continue
		}
		// fmt.Printf("Acc Max: %v / Vel Max: %v / Temperature: %v", realTimeMeasurements.AccMax, realTimeMeasurements.VelMax, realTimeMeasurements.T)
		data := make(map[string]interface{})
		data["Timestamp"] = time.Now().Unix()
		data["SerialNumber"] = serialNumber
		data["VectorialMeasures"] = vectorialMeasures
		// res, _ := json.Marshal(vectorialMeasures)
		res, _ := json.Marshal(data)
		// global.Logger.Info(string(res))
		// _, err = global.Publisher.Send(string(res), 0)
		_, err = global.Publisher.Send(fmt.Sprintf("%v %v", global.Conf.Filter, string(res)), 0)
		if err != nil {
			fmt.Println(err)
			global.Logger.Warn(fmt.Sprintf("Sending realtime measurements from cube '%v' via ZeroMQ failed: %v", cube.Name, err))
			continue
		}
		global.Logger.Info(fmt.Sprintf("VECTOR PUBLISHED: %v", string(res)))
	}
}
