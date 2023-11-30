package cube_modbus_handler

import (
	. "cube_config_handler"
	"fmt"

	"github.com/goburrow/modbus"
)

func ReadHoldingRegister(cube Cube, address, quantity uint16) ([]byte, error) {
	handler := modbus.NewTCPClientHandler(fmt.Sprintf("%v:%v", cube.IP, cube.Port))
	handler.SlaveId = 1
	handler.Timeout = 5

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

func ReadInputRegister(cube Cube, address, quantity uint16) ([]byte, error) {
	handler := modbus.NewTCPClientHandler(fmt.Sprintf("%v:%v", cube.IP, cube.Port))
	handler.SlaveId = byte(cube.SlaveId)
	handler.Timeout = 5

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
