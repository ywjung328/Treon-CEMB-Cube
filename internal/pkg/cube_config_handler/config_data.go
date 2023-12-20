package cube_config_handler

type Cube struct {
	Name    string `yaml:"name"`
	IP      string `yaml:"ip"`
	Port    int    `yaml:"port"`
	SlaveId int    `yaml:"slave_id"`
}

type CubeConfig struct {
	GatewayId     string `yaml:"gateway_id"`
	PublishPort   int    `yaml:"publish_port"`
	SubscribePort int    `yaml:"subscribe_port"`
	Filter        string `yaml:"filter"`
	Cubes         []Cube `yaml:"cubes"`
	ScalarCycle   int    `yaml:"scalar_cycle"`
	VectorCycle   int    `yaml:"vector_cycle"`
	ModbusTimeout int    `yaml:"modbus_timeout"`
}
