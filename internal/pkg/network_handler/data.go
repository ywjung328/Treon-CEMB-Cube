package network_handler

import (
	. "cube_config_handler"
	. "cube_data_handler"
	"encoding/json"
	"global"
	"time"
)

func TreonRTMConverter(realTimeMeasurements RealTimeMeasurements, cube Cube) string {
	now := time.Now()
	ts := now.Unix()

	data := make(map[string]interface{})
	data["Temperature"] = realTimeMeasurements.T
	data["ShaftSpeed"] = realTimeMeasurements.ShaftSpeed
	data["Timestamp"] = ts
	data["GatewayId"] = global.Conf.GatewayId
	data["CubeId"] = cube.Name
	data["Type"] = "scalar"
	// acc := map[string]float32 {

	// }

	formattedData, err := json.Marshal(data)

	if err != nil {
		// Handle error here
	}

	return string(formattedData)
}

type ScalarData struct {
	AccX        float32 `json:"AccX"`
	AccY        float32 `json:"AccY"`
	AccZ        float32 `json:"AccZ"`
	VelX        float32 `json:"VelX"`
	VelY        float32 `json:"VelY"`
	VelZ        float32 `json:"VelZ"`
	Temperature float32 `json:"Temperature"`
}

type ScalarJson struct {
	DeviceType   string     `json:"DeviceType"`
	DataType     string     `json:"DataType"`
	Timestamp    int64      `json:"Timestamp"`
	SerialNumber string     `json:"SeiralNumber"`
	GatewayId    string     `json:"GatewayId"`
	Data         ScalarData `json:"Data"`
}

type VectorJson struct {
	DeviceType   string            `json:"DeviceType"`
	DataType     string            `json:"DataType"`
	Timestamp    int64             `json:"Timestamp"`
	SerialNumber string            `json:"SeiralNumber"`
	GatewayId    string            `json:"GatewayId"`
	Data         VectorialMeasures `json:"Data"`
}
