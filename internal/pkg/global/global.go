package global

import (
	"os"

	. "cube_config_handler"

	"go.uber.org/zap"
)

var (
	LogFile *os.File
	Logger  *zap.Logger
	Conf    CubeConfig
	// Publisher, Subscriber *zmq4.Socket
)
