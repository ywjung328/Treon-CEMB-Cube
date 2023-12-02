package main

import (
	. "cube_config_handler"
	"fmt"
)

func init() {
	fmt.Println("THIS IS INIT")
}

func main() {
	fmt.Println("THIS IS MAIN")
	conf, err := ReadConfig("./conf.yaml")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(conf)
}
