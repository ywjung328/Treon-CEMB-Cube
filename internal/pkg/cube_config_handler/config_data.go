package cube_config_handler

import . "cube_modbus_handler"

type CubeConfig struct {
	Cubes []Cube `yaml:"cubes"`
}
