package cube_modbus_handler

import (
	. "cube_config_handler"
	"encoding/binary"
	"fmt"
	"math"
)

func byteArrayToFloat32(data []byte) (float32, error) {
	if len(data) != 4 {
		return float32(0), fmt.Errorf("Invalid byte length: expected 4 bytes for float32, but got %v bytes", len(data))
	}
	bits := binary.BigEndian.Uint32(data)

	return math.Float32frombits(bits), nil
}

func byteArrayToUint16(data []byte) (uint16, error) {
	if len(data) != 2 {
		return uint16(0), fmt.Errorf("Invalid byte length: expected 2 bytes for uint16, but got %v bytes", len(data))
	}

	return binary.BigEndian.Uint16(data), nil
}

func GetRealTimeMeasurements(cube Cube) (RealTimeMeasurements, error) {
	var realTimeMeasurements RealTimeMeasurements

	accResult, err := ReadInputRegister(cube, uint16(0), uint16(8))
	if err != nil {
		return realTimeMeasurements, err
	} else if len(accResult) != 16 {
		return realTimeMeasurements, fmt.Errorf("Invalid byte length: expected 16 bytes for accResult, but got %v bytes", len(accResult))
	}
	accX, err := byteArrayToFloat32(accResult[0:3])
	if err != nil {
		return realTimeMeasurements, err
	}
	accY, err := byteArrayToFloat32(accResult[4:7])
	if err != nil {
		return realTimeMeasurements, err
	}
	accZ, err := byteArrayToFloat32(accResult[8:11])
	if err != nil {
		return realTimeMeasurements, err
	}
	accMax, err := byteArrayToFloat32(accResult[12:15])
	if err != nil {
		return realTimeMeasurements, err
	}
	realTimeMeasurements.AccX = accX
	realTimeMeasurements.AccY = accY
	realTimeMeasurements.AccZ = accZ
	realTimeMeasurements.AccMax = accMax

	velResult, err := ReadInputRegister(cube, uint16(10), uint16(8))
	if err != nil {
		return realTimeMeasurements, err
	} else if len(velResult) != 16 {
		return realTimeMeasurements, fmt.Errorf("Invalid byte length: expected 16 bytes for velResult, but got %v bytes", len(velResult))
	}
	velX, err := byteArrayToFloat32(accResult[0:3])
	if err != nil {
		return realTimeMeasurements, err
	}
	velY, err := byteArrayToFloat32(accResult[4:7])
	if err != nil {
		return realTimeMeasurements, err
	}
	velZ, err := byteArrayToFloat32(accResult[8:11])
	if err != nil {
		return realTimeMeasurements, err
	}
	velMax, err := byteArrayToFloat32(accResult[12:15])
	if err != nil {
		return realTimeMeasurements, err
	}
	realTimeMeasurements.VelX = velX
	realTimeMeasurements.VelY = velY
	realTimeMeasurements.VelZ = velZ
	realTimeMeasurements.VelMax = velMax

	otherResult, err := ReadInputRegister(cube, uint16(20), uint16(12))
	if err != nil {
		return realTimeMeasurements, err
	} else if len(otherResult) != 24 {
		return realTimeMeasurements, fmt.Errorf("Invalid byte length: expected 24 bytes for otherResult, but got %v bytes", len(velResult))
	}
	t, err := byteArrayToFloat32(otherResult[0:3])
	if err != nil {
		return realTimeMeasurements, err
	}
	shaftSpeed, err := byteArrayToFloat32(otherResult[4:7])
	if err != nil {
		return realTimeMeasurements, err
	}
	velDAUnbalanced, err := byteArrayToFloat32(otherResult[8:11])
	if err != nil {
		return realTimeMeasurements, err
	}
	velDAMisalignment, err := byteArrayToFloat32(otherResult[12:15])
	if err != nil {
		return realTimeMeasurements, err
	}
	velDALooseness, err := byteArrayToFloat32(otherResult[16:19])
	if err != nil {
		return realTimeMeasurements, err
	}
	velDAOther, err := byteArrayToFloat32(otherResult[20:23])
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

	result, err := ReadInputRegister(cube, uint16(64), uint16(5))
	if err != nil {
		return analogDigitalOutputs, err
	} else if len(result) != 10 {
		return analogDigitalOutputs, fmt.Errorf("Invalid byte length: expected 10 bytes for result, but got %v bytes", len(result))
	}
	u1Analog, err := byteArrayToFloat32(result[0:3])
	if err != nil {
		return analogDigitalOutputs, err
	}
	u2Analog, err := byteArrayToFloat32(result[4:7])
	if err != nil {
		return analogDigitalOutputs, err
	}
	u2Digital, err := byteArrayToUint16(result[8:9])
	if err != nil {
		return analogDigitalOutputs, err
	}
	analogDigitalOutputs.U1Analog = u1Analog
	analogDigitalOutputs.U2Analog = u2Analog
	analogDigitalOutputs.U2Digital = u2Digital

	return analogDigitalOutputs, nil
}
