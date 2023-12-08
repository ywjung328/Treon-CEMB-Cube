package cube_config_handler

type Cube struct {
	Name    string `yaml:"name"`
	IP      string `yaml:"ip"`
	Port    int    `yaml:"port"`
	SlaveId int    `yaml:"slave_id"`
}

type CubeConfig struct {
	PublishPort   int    `yaml:"publish_port"`
	SubscribePort int    `yaml:"subscribe_port"`
	Cubes         []Cube `yaml:"cubes"`
	ModbusTimeout int    `yaml:"modbus_timeout"`
}
