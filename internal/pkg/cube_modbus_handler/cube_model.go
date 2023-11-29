package cube_modbus_handler

type Cube struct {
	Name string `yaml:"name"`
	IP string `yaml:"ip"`
	Port int `yaml:"port"`
}