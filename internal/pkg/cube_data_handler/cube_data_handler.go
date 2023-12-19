package cube_data_handler

import (
	. "cube_config_handler"
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"net/http"
	"reflect"

	"github.com/fxamacker/cbor/v2"
)

func byteArrayToFloat32(data []byte) (float32, error) {
	if len(data) != 4 {
		return float32(0), fmt.Errorf("Invalid byte length: expected 4 bytes for float32, but got %v bytes", len(data))
	}
	// slices.Reverse(data)
	bits := binary.BigEndian.Uint32([]byte{data[2], data[3], data[0], data[1]})

	return math.Float32frombits(bits), nil
}

func byteArrayToUint16(data []byte) (uint16, error) {
	if len(data) != 2 {
		return uint16(0), fmt.Errorf("Invalid byte length: expected 2 bytes for uint16, but got %v bytes", len(data))
	}
	// slices.Reverse(data)
	return binary.BigEndian.Uint16(data), nil
}

func GetRealTimeMeasurements(cube Cube) (RealTimeMeasurements, error) {
	var realTimeMeasurements RealTimeMeasurements

	accResult, err := ReadInputRegisters(cube, uint16(0), uint16(8))
	if err != nil {
		return realTimeMeasurements, err
	} else if len(accResult) != 16 {
		return realTimeMeasurements, fmt.Errorf("Invalid byte length: expected 16 bytes for accResult, but got %v bytes", len(accResult))
	}
	accX, err := byteArrayToFloat32(accResult[0:4])
	if err != nil {
		return realTimeMeasurements, err
	}
	accY, err := byteArrayToFloat32(accResult[4:8])
	if err != nil {
		return realTimeMeasurements, err
	}
	accZ, err := byteArrayToFloat32(accResult[8:12])
	if err != nil {
		return realTimeMeasurements, err
	}
	accMax, err := byteArrayToFloat32(accResult[12:16])
	if err != nil {
		return realTimeMeasurements, err
	}
	realTimeMeasurements.AccX = accX
	realTimeMeasurements.AccY = accY
	realTimeMeasurements.AccZ = accZ
	realTimeMeasurements.AccMax = accMax

	velResult, err := ReadInputRegisters(cube, uint16(10), uint16(8))
	if err != nil {
		return realTimeMeasurements, err
	} else if len(velResult) != 16 {
		return realTimeMeasurements, fmt.Errorf("Invalid byte length: expected 16 bytes for velResult, but got %v bytes", len(velResult))
	}
	velX, err := byteArrayToFloat32(velResult[0:4])
	if err != nil {
		return realTimeMeasurements, err
	}
	velY, err := byteArrayToFloat32(velResult[4:8])
	if err != nil {
		return realTimeMeasurements, err
	}
	velZ, err := byteArrayToFloat32(velResult[8:12])
	if err != nil {
		return realTimeMeasurements, err
	}
	velMax, err := byteArrayToFloat32(velResult[12:16])
	if err != nil {
		return realTimeMeasurements, err
	}
	realTimeMeasurements.VelX = velX
	realTimeMeasurements.VelY = velY
	realTimeMeasurements.VelZ = velZ
	realTimeMeasurements.VelMax = velMax

	otherResult, err := ReadInputRegisters(cube, uint16(20), uint16(12))
	if err != nil {
		return realTimeMeasurements, err
	} else if len(otherResult) != 24 {
		return realTimeMeasurements, fmt.Errorf("Invalid byte length: expected 24 bytes for otherResult, but got %v bytes", len(velResult))
	}
	t, err := byteArrayToFloat32(otherResult[0:4])
	if err != nil {
		return realTimeMeasurements, err
	}
	shaftSpeed, err := byteArrayToFloat32(otherResult[4:8])
	if err != nil {
		return realTimeMeasurements, err
	}
	velDAUnbalanced, err := byteArrayToFloat32(otherResult[8:12])
	if err != nil {
		return realTimeMeasurements, err
	}
	velDAMisalignment, err := byteArrayToFloat32(otherResult[12:16])
	if err != nil {
		return realTimeMeasurements, err
	}
	velDALooseness, err := byteArrayToFloat32(otherResult[16:20])
	if err != nil {
		return realTimeMeasurements, err
	}
	velDAOther, err := byteArrayToFloat32(otherResult[20:24])
	if err != nil {
		return realTimeMeasurements, err
	}
	realTimeMeasurements.T = t
	realTimeMeasurements.ShaftSpeed = shaftSpeed
	realTimeMeasurements.VelDAUnbalanced = velDAUnbalanced
	realTimeMeasurements.VelDAMisalignment = velDAMisalignment
	realTimeMeasurements.VelDALooseness = velDALooseness
	realTimeMeasurements.VelDAOther = velDAOther

	return realTimeMeasurements, nil
}

func GetAnalogDigitalOutputs(cube Cube) (AnalogDigitalOutputs, error) {
	var analogDigitalOutputs AnalogDigitalOutputs

	result, err := ReadInputRegisters(cube, uint16(64), uint16(5))
	if err != nil {
		return analogDigitalOutputs, err
	} else if len(result) != 10 {
		return analogDigitalOutputs, fmt.Errorf("Invalid byte length: expected 10 bytes for result, but got %v bytes", len(result))
	}
	u1Analog, err := byteArrayToFloat32(result[0:4])
	if err != nil {
		return analogDigitalOutputs, err
	}
	u2Analog, err := byteArrayToFloat32(result[4:8])
	if err != nil {
		return analogDigitalOutputs, err
	}
	u2Digital, err := byteArrayToUint16(result[8:10])
	if err != nil {
		return analogDigitalOutputs, err
	}
	analogDigitalOutputs.U1Analog = u1Analog
	analogDigitalOutputs.U2Analog = u2Analog
	analogDigitalOutputs.U2Digital = u2Digital

	return analogDigitalOutputs, nil
}

func GetDeviceStatuses(cube Cube) (DeviceStatuses, error) {
	var deviceStatuses DeviceStatuses
	result, err := ReadDiscreteInputs(cube, uint16(0), uint16(8))
	if err != nil {
		return deviceStatuses, err
	} else if len(result) != 1 {
		return deviceStatuses, fmt.Errorf("Invalid byte length: expected 1 bytes for result, but got %v bytes", len(result))
	}
	var boolValues []bool
	for _, byte := range result {
		for i := 0; i < 8; i++ {
			boolValues = append(boolValues, (byte&(1<<uint(i))) != 0)
		}
	}
	deviceStatuses.GoodOperatingMode = boolValues[0]
	deviceStatuses.AnyHardwareFault = boolValues[1]
	deviceStatuses.MEMSHardwareFault = boolValues[2]
	deviceStatuses.EEPROMHardwareFault = boolValues[3]
	deviceStatuses.ConfigurationFault = boolValues[4]
	deviceStatuses.BootInProgress = boolValues[5]
	deviceStatuses.DigitalOutputDriverFault = boolValues[6]
	deviceStatuses.InvalidSetup = boolValues[7]

	return deviceStatuses, nil
}

func GetChannelStatuses(cube Cube) (ChannelStatuses, error) {
	var channelStatuses ChannelStatuses
	result, err := ReadDiscreteInputs(cube, uint16(16), uint16(52))
	if err != nil {
		return channelStatuses, err
	} else if len(result) != 7 {
		return channelStatuses, fmt.Errorf("Invalid byte length: expected 7 bytes for result, but got %v bytes", len(result))
	}
	var boolValues []bool
	for _, byte := range result {
		for i := 0; i < 8; i++ {
			boolValues = append(boolValues, (byte&(1<<uint(i))) != 0)
		}
	}
	channelStatuses.XOn = boolValues[0]
	channelStatuses.XInvalid = boolValues[1]
	channelStatuses.XFault = boolValues[2]
	channelStatuses.XSaturation = boolValues[3]
	channelStatuses.YOn = boolValues[16]
	channelStatuses.YInvalid = boolValues[17]
	channelStatuses.YFault = boolValues[18]
	channelStatuses.YSaturation = boolValues[19]
	channelStatuses.ZOn = boolValues[32]
	channelStatuses.ZInvalid = boolValues[33]
	channelStatuses.ZFault = boolValues[34]
	channelStatuses.ZSaturation = boolValues[35]
	channelStatuses.TOn = boolValues[48]
	channelStatuses.TInvalid = boolValues[49]
	channelStatuses.TFault = boolValues[50]
	channelStatuses.TSaturation = boolValues[51]

	return channelStatuses, nil
}

func GetMeasurementsStatuses(cube Cube) (MeasurementsStatuses, error) {
	var measurementsStatuses MeasurementsStatuses
	result, err := ReadDiscreteInputs(cube, uint16(80), uint16(16))
	if err != nil {
		return measurementsStatuses, err
	} else if len(result) != 2 {
		return measurementsStatuses, fmt.Errorf("Invalid byte length: expected 2 bytes for result, but got %v bytes", len(result))
	}
	var boolValues []bool
	for _, byte := range result {
		for i := 0; i < 8; i++ {
			boolValues = append(boolValues, (byte&(1<<uint(i))) != 0)
		}
	}
	measurementsStatuses.AnyAlert = boolValues[0]
	measurementsStatuses.AnyVibrationAlert = boolValues[1]
	measurementsStatuses.AnyAccAlert = boolValues[2]
	measurementsStatuses.AnyVelAlert = boolValues[3]
	measurementsStatuses.AnyXAlert = boolValues[4]
	measurementsStatuses.AnyYAlert = boolValues[5]
	measurementsStatuses.AnyZAlert = boolValues[6]
	measurementsStatuses.AccXAlert = boolValues[7]
	measurementsStatuses.AccYAlert = boolValues[8]
	measurementsStatuses.AccZAlert = boolValues[9]
	measurementsStatuses.VelXAlert = boolValues[10]
	measurementsStatuses.VelYAlert = boolValues[11]
	measurementsStatuses.VelZAlert = boolValues[12]
	measurementsStatuses.AccMaxAlert = boolValues[13]
	measurementsStatuses.VelMaxAlert = boolValues[14]
	measurementsStatuses.TAlert = boolValues[15]

	return measurementsStatuses, nil
}

func GetVectorialMeasures(cube Cube) (VectorialMeasures, error) {
	var vectorialMeasures VectorialMeasures
	url := fmt.Sprintf("http://%v/rawdata.cbor", cube.IP)
	resp, err := http.Get(url)
	if err != nil {
		return vectorialMeasures, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return vectorialMeasures, fmt.Errorf("Bad status code while fetching vectorial measures: %v", resp.StatusCode)
	}
	cborData, err := io.ReadAll(resp.Body)
	if err != nil {
		return vectorialMeasures, err
	}
	var decodedData map[string]interface{}
	err = cbor.Unmarshal(cborData, &decodedData)
	if err != nil {
		return vectorialMeasures, err
	}

	// DATA CASTING
	vmUnit, ok := decodedData["unit"].(string)
	if !ok {
		return vectorialMeasures, fmt.Errorf("Failed parsing 'unit' from vectorial measures")
	}
	vmValuesX, ok := decodedData["values_x"].(cbor.Tag)
	if !ok {
		return vectorialMeasures, fmt.Errorf("Failed parsing 'values_x' from vectorial measures: failed casting")
	} else if reflect.TypeOf(vmValuesX.Content).String() != "[]uint8" {
		return vectorialMeasures, fmt.Errorf("Failed parsing 'values_x' from vectorial measures: type not []uint8")
	}
	vmValuesY, ok := decodedData["values_y"].(cbor.Tag)
	if !ok {
		return vectorialMeasures, fmt.Errorf("Failed parsing 'values_y' from vectorial measures: failed casting")
	} else if reflect.TypeOf(vmValuesY.Content).String() != "[]uint8" {
		return vectorialMeasures, fmt.Errorf("Failed parsing 'values_y' from vectorial measures: type not []uint8")
	}
	vmValuesZ, ok := decodedData["values_z"].(cbor.Tag)
	if !ok {
		return vectorialMeasures, fmt.Errorf("Failed parsing 'values_z' from vectorial measures: failed casting")
	} else if reflect.TypeOf(vmValuesY.Content).String() != "[]uint8" {
		return vectorialMeasures, fmt.Errorf("Failed parsing 'values_z' from vectorial measures: type not []uint8")
	}
	vmDtRaw, ok := decodedData["dt"].(map[interface{}]interface{})
	if !ok {
		return vectorialMeasures, fmt.Errorf("Failed parsing 'dt' from vectorial measures")
	}
	vmDtUnit, ok := vmDtRaw["unit"].(string)
	if !ok {
		return vectorialMeasures, fmt.Errorf("Failed parsing / casting 'dt'->'unit' from vectorial measures")
	}
	vmDtValue, ok := vmDtRaw["value"].(float64)
	if !ok {
		return vectorialMeasures, fmt.Errorf("Failed parsing / casting 'dt'->'value' from vectorial measures")
	}

	vectorialMeasures.Unit = vmUnit
	vectorialMeasures.ValuesX = vmValuesX.Content.([]uint8)
	vectorialMeasures.ValuesY = vmValuesY.Content.([]uint8)
	vectorialMeasures.ValuesZ = vmValuesZ.Content.([]uint8)
	vectorialMeasures.Dt = dt{Unit: vmDtUnit, Value: vmDtValue}

	return vectorialMeasures, nil
}

func GetSerialNumber(cube Cube) (string, error) {
	var serialNumber string
	result, err := ReadHoldingRegisters(cube, uint16(17), uint16(4))
	if err != nil {
		return serialNumber, err
	}
	// for i := 0; i < len(result); i += 2 {
	// 	temp := result[i]
	// 	result[i] = result[i+1]
	// 	result[i+1] = temp
	// }
	serialNumber = string(result)
	return serialNumber, nil
}
