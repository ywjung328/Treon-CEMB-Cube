package network_handler

import (
	"bytes"
	. "cube_data_handler"
	"encoding/json"
	"fmt"
	"global"
	"io"
	"net/http"
	"time"
)

var err error

/*
func InitZeroMQ() error {
	context, err := zmq4.NewContext()
	if err != nil {
		return fmt.Errorf("Creating zmq4 context failed: %v", err)
	}
	global.Publisher, err = context.NewSocket(zmq4.PUB)
	if err != nil {
		return fmt.Errorf("Creating zmq4 publisher failed: %v", err)
	}
	global.Subscriber, err = context.NewSocket(zmq4.SUB)
	if err != nil {
		return fmt.Errorf("Creating zmq4 subscriber failed: %v", err)
	}
	// err = global.Publisher.Bind(fmt.Sprintf("tcp://*:%v", global.Conf.PublishPort))
	err = global.Subscriber.Connect(fmt.Sprintf("tcp://localhost:%v", global.Conf.SubscribePort))
	if err != nil {
		return fmt.Errorf("Connecting zmq4 subscriber to tcp://localhost:%v failed: %v", global.Conf.SubscribePort, err)
	}
	err = global.Publisher.Connect(fmt.Sprintf("tcp://*:%v", global.Conf.PublishPort))
	if err != nil {
		return fmt.Errorf("Binding zmq4 publisher to tcp://*:%v failed: %v", global.Conf.PublishPort, err)
	}
	err = global.Subscriber.SetSubscribe(global.Conf.Filter)
	if err != nil {
		return fmt.Errorf("Setting subscribtion (filter: %v) failed: %v", global.Conf.Filter, err)
	}
	return nil
}
*/

/*---- WE DON'T USE ZMQ ANYMORE!!! ----*/

/*
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
			// analogDigitalOutputs, err := GetAnalogDigitalOutputs(cube)
			// if err != nil {
			// 	global.Logger.Warn(fmt.Sprintf("Fetching analog digital outputs from cube '%v' failed: %v", cube.Name, err))
			// 	continue
			// }
			// deviceStatuses, err := GetDeviceStatuses(cube)
			// if err != nil {
			// 	global.Logger.Warn(fmt.Sprintf("Fetching device statuses from cube '%v' failed: %v", cube.Name, err))
			// 	continue
			// }
			// channelStatuses, err := GetChannelStatuses(cube)
			// if err != nil {
			// 	global.Logger.Warn(fmt.Sprintf("Fetching channel statuses from cube '%v' failed: %v", cube.Name, err))
			// 	continue
			// }
			// measurementsStatuses, err := GetMeasurementsStatuses(cube)
			// if err != nil {
			// 	global.Logger.Warn(fmt.Sprintf("Fetching measurements statuses from cube '%v' failed: %v", cube.Name, err))
			// 	continue
			// }
		data := make(map[string]interface{})
		data["DeviceType"] = "CUBE"
		data["DataType"] = "Scalar"
		data["Timestamp"] = time.Now().Unix()
		data["SerialNumber"] = serialNumber
		data["RealTimeMeasurements"] = realTimeMeasurements
		// data["AnalogDigitalOutputs"] = analogDigitalOutputs
		// data["DeviceStatuses"] = deviceStatuses
		// data["ChannelStatuses"] = channelStatuses
		// data["MeasurementStatuses"] = measurementsStatuses
		res, _ := json.Marshal(data)
		// global.Logger.Info(string(res))
		// _, err = global.Publisher.Send(string(res), 0)
		_, err = global.Publisher.Send(fmt.Sprintf("%v %v", global.Conf.Filter, string(res)), 0)
		fmt.Println(fmt.Sprintf("%v %v", global.Conf.Filter, string(res)), 0)
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
		data["DeviceType"] = "CUBE"
		data["DataType"] = "Vector"
		data["Timestamp"] = time.Now().Unix()
		data["SerialNumber"] = serialNumber
		data["VectorialMeasures"] = vectorialMeasures
		// res, _ := json.Marshal(vectorialMeasures)
		res, _ := json.Marshal(data)
		// global.Logger.Info(string(res))
		// _, err = global.Publisher.Send(string(res), 0)
		_, err = global.Publisher.Send(fmt.Sprintf("%v %v", global.Conf.Filter, string(res)), 0)
		fmt.Println(fmt.Sprintf("%v %v", global.Conf.Filter, string(res)), 0)
		if err != nil {
			fmt.Println(err)
			global.Logger.Warn(fmt.Sprintf("Sending realtime measurements from cube '%v' via ZeroMQ failed: %v", cube.Name, err))
			continue
		}
		global.Logger.Info(fmt.Sprintf("VECTOR PUBLISHED: %v", string(res)))
	}
}
*/

func ScalarPublish() {
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
		data := ScalarJson{
			DeviceType:   "cube",
			DataType:     "scalar",
			Timestamp:    time.Now().Unix(),
			SerialNumber: serialNumber,
			GatewayId:    global.Conf.GatewayId,
			Data: ScalarData{
				AccX:        realTimeMeasurements.AccX,
				AccY:        realTimeMeasurements.AccY,
				AccZ:        realTimeMeasurements.AccZ,
				VelX:        realTimeMeasurements.VelX,
				VelY:        realTimeMeasurements.VelY,
				VelZ:        realTimeMeasurements.VelZ,
				Temperature: realTimeMeasurements.T,
			},
		}
		jsonData, err := json.Marshal(data)
		if err != nil {
			global.Logger.Warn(fmt.Sprintf("Error marshaling scalar JSON: %v", err))
			continue
		}
		req, err := http.NewRequest("POST", global.Conf.PublishAddress, bytes.NewBuffer(jsonData))
		if err != nil {
			global.Logger.Warn(fmt.Sprintf("Error creating scalar request: %v", err))
			continue
		}
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			global.Logger.Warn(fmt.Sprintf("Error making scalar request to server: %v", err))
			continue
		}
		defer resp.Body.Close()
		if resp.StatusCode == http.StatusOK {
			continue
		}
		errBody, err := io.ReadAll(resp.Body)
		if err != nil {
			global.Logger.Warn(fmt.Sprintf("Error reading error from scalar response body: %v", err))
			continue
		}
		global.Logger.Warn(fmt.Sprintf("Error from scalar response(%d): %v", resp.StatusCode, string(errBody)))
		continue
	}
}

func VectorPublish() {
	for _, cube := range global.Conf.Cubes {
		serialNumber, err := GetSerialNumber(cube)
		if err != nil {
			global.Logger.Warn(fmt.Sprintf("Fetching S/N from cube '%v' failed: %v", cube.Name, err))
			continue
		}
		vectorialMeasures, err := GetVectorialMeasures(cube)
		if err != nil {
			global.Logger.Warn(fmt.Sprintf("Fetching vectorial measurements from cube '%v' failed: %v", cube.Name, err))
			continue
		}
		data := VectorJson{
			DeviceType:   "cube",
			DataType:     "vector",
			Timestamp:    time.Now().Unix(),
			SerialNumber: serialNumber,
			GatewayId:    global.Conf.GatewayId,
			Data:         vectorialMeasures,
		}
		jsonData, err := json.Marshal(data)
		if err != nil {
			global.Logger.Warn(fmt.Sprintf("Error marshaling vector JSON: %v", err))
			continue
		}
		req, err := http.NewRequest("POST", global.Conf.PublishAddress, bytes.NewBuffer(jsonData))
		if err != nil {
			global.Logger.Warn(fmt.Sprintf("Error creating vector request: %v", err))
			continue
		}
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			global.Logger.Warn(fmt.Sprintf("Error making vector request to server: %v", err))
			continue
		}
		defer resp.Body.Close()
		if resp.StatusCode == http.StatusOK {
			continue
		}
		errBody, err := io.ReadAll(resp.Body)
		if err != nil {
			global.Logger.Warn(fmt.Sprintf("Error reading error from vector response body: %v", err))
			continue
		}
		global.Logger.Warn(fmt.Sprintf("Error from vector response(%d): %v", resp.StatusCode, string(errBody)))
		continue
	}
}
