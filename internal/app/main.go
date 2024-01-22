package main

import (
	. "cube_config_handler"
	"fmt"
	"global"
	"log"
	network "network_handler"

	// . "network_handler"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func init() {
	if global.Conf.SaveLog {
		err := os.MkdirAll("logs", os.ModePerm)
		if err != nil {
			panic(err)
		}
		now := time.Now()
		global.LogFile, err = os.Create(fmt.Sprintf("logs/log_%v.log", now.Format(time.RFC3339)))
		if err != nil {
			panic(err)
		}
		log.SetOutput(global.LogFile)
	}
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
	// if err != nil {
	// 	global.Logger.Warn(fmt.Sprintf("ZeroMQ initiation failed: %v", err))
	// } else {
	// 	global.Logger.Info("ZeroMQ initiated successfully.")
	// }
}

func main() {
	defer global.LogFile.Close()
	defer global.Logger.Sync()
	// defer global.Publisher.Close()
	// defer global.Subscriber.Close()

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

	/*---- WE DON'T USE ZMQ ANYMORE!!! ----*/

	// // INITIATING ZMQ4
	// context, err := zmq4.NewContext()
	// if err != nil {
	// 	global.Logger.Error(fmt.Sprintf("Creating zmq4 context failed: %v", err))
	// }
	// global.Publisher, err = context.NewSocket(zmq4.PUB)
	// if err != nil {
	// 	global.Logger.Error(fmt.Sprintf("Creating zmq4 publisher failed: %v", err))
	// }
	// global.Subscriber, err = context.NewSocket(zmq4.SUB)
	// if err != nil {
	// 	global.Logger.Error(fmt.Sprintf("Creating zmq4 subscriber failed: %v", err))
	// }
	// // err = global.Publisher.Connect(fmt.Sprintf("tcp://127.0.0.1:%v", global.Conf.PublishPort))
	// err = global.Publisher.Bind(fmt.Sprintf("tcp://127.0.0.1:%v", global.Conf.PublishPort))
	// if err != nil {
	// 	global.Logger.Error(fmt.Sprintf("Connecting zmq4 publisher to tcp://127.0.0.1:%v failed: %v", global.Conf.PublishPort, err))
	// }
	// err = global.Subscriber.Connect(fmt.Sprintf("tcp://127.0.0.1:%v", global.Conf.SubscribePort))
	// if err != nil {
	// 	global.Logger.Error(fmt.Sprintf("Connecting zmq4 subscriber to tcp://127.0.0.1:%v failed: %v", global.Conf.SubscribePort, err))
	// }
	// err = global.Subscriber.SetSubscribe(global.Conf.Filter)
	// if err != nil {
	// 	global.Logger.Error(fmt.Sprintf("Setting subscribtion (filter: %v) failed: %v", global.Conf.Filter, err))
	// }

	// defer context.Term()
	// defer global.Publisher.Close()
	// defer global.Subscriber.Close()

	// go subscribe()
	for {
		select {
		case <-scalarTicker.C:
			go network.ScalarPublish()
		case <-vectorTicker.C:
			go network.VectorPublish()
		}
	}
}
