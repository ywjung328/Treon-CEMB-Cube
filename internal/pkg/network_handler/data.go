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

	var data map[string]interface{}
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
