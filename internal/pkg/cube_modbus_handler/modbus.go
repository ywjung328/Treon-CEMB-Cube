package cube_modbus_handler

import (
	. "cube_config_handler"
	"fmt"
	. "global"
	"time"

	"github.com/goburrow/modbus"
)

func ReadHoldingRegisters(cube Cube, address, quantity uint16) ([]byte, error) {
	handler := modbus.NewTCPClientHandler(fmt.Sprintf("%v:%v", cube.IP, cube.Port))
	handler.SlaveId = byte(cube.SlaveId)
	handler.Timeout = time.Duration(Conf.ModbusTimeout) * time.Second

	err := handler.Connect()
	defer handler.Close()

	if err != nil {
		return nil, err
	}

	client := modbus.NewClient(handler)
	results, err := client.ReadHoldingRegisters(address, quantity)

	if err != nil {
		return nil, err
	}
	return results, nil
}

func ReadInputRegisters(cube Cube, address, quantity uint16) ([]byte, error) {
	handler := modbus.NewTCPClientHandler(fmt.Sprintf("%v:%v", cube.IP, cube.Port))
	handler.SlaveId = byte(cube.SlaveId)
	handler.Timeout = time.Duration(Conf.ModbusTimeout) * time.Second

	err := handler.Connect()
	defer handler.Close()

	if err != nil {
		return nil, err
	}

	client := modbus.NewClient(handler)
	// results, err := client.ReadHoldingRegisters(address, quantity)
	results, err := client.ReadInputRegisters(address, quantity)

	if err != nil {
		return nil, err
	}
	return results, nil
}

func ReadCoils(cube Cube, address, quantity uint16) ([]byte, error) {
	handler := modbus.NewTCPClientHandler(fmt.Sprintf("%v:%v", cube.IP, cube.Port))
	handler.SlaveId = byte(cube.SlaveId)
	handler.Timeout = time.Duration(Conf.ModbusTimeout) * time.Second

	err := handler.Connect()
	defer handler.Close()

	if err != nil {
		return nil, err
	}

	client := modbus.NewClient(handler)
	results, err := client.ReadCoils(address, quantity)

	if err != nil {
		return nil, err
	}
	return results, nil
}

func ReadDiscreteInputs(cube Cube, address, quantity uint16) ([]byte, error) {
	handler := modbus.NewTCPClientHandler(fmt.Sprintf("%v:%v", cube.IP, cube.Port))
	handler.SlaveId = byte(cube.SlaveId)
	handler.Timeout = time.Duration(Conf.ModbusTimeout) * time.Second

	err := handler.Connect()
	defer handler.Close()

	if err != nil {
		return nil, err
	}

	client := modbus.NewClient(handler)
	results, err := client.ReadDiscreteInputs(address, quantity)

	if err != nil {
		return nil, err
	}
	return results, nil
}
