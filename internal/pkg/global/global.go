package global

import (
	"os"

	. "cube_config_handler"

	"github.com/pebbe/zmq4"
	"go.uber.org/zap"
)

const (
	Filter = "cube"
)

var (
	LogFile               *os.File
	Logger                *zap.Logger
	Conf                  CubeConfig
	Publisher, Subscriber *zmq4.Socket
)
