package cube_modbus_handler

import (
	"fmt"

	"github.com/goburrow/modbus"
)

/*
func initHandler(cube Cube) (modbus.Client, error) {
	handler := modbus.NewTCPClientHandler(fmt.Sprintf("%v:%v", cube.IP, cube.Port))
	handler.SlaveId = 1
	handler.Timeout = 5

	err := handler.Connect()
	if err != nil {
		return nil, err
	}

	// Is it right?
	defer handler.Close()

	client := modbus.NewClient(handler)

	return client, nil
}

func readRegister(client modbus.Client, address, quantity uint16) ([]byte, error) {
	results, err := client.ReadHoldingRegisters(address, quantity)
	if err != nil {
		return nil, err
	}
	return results, nil
}
*/

func readRegister(cube Cube, address, quantity uint16) ([]byte, error) {
	handler := modbus.NewTCPClientHandler(fmt.Sprintf("%v:%v", cube.IP, cube.Port))
	handler.SlaveId = 1
	handler.Timeout = 5

	err := handler.Connect()

	if err != nil {
		return nil, err
	}

	defer handler.Close()

	client := modbus.NewClient(handler)
	results, err := client.ReadHoldingRegisters(address, quantity)

	if err != nil {
		return nil, err
	}
	return results, nil
}
