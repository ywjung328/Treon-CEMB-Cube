package cube_config_handler

type Cube struct {
	Name    string `yaml:"name"`
	IP      string `yaml:"ip"`
	Port    int    `yaml:"port"`
	SlaveId int    `yaml:"slave_id"`
}

type CubeConfig struct {
	Cubes []Cube `yaml:"cubes"`
}
