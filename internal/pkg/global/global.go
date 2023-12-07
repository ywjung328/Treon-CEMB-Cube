package global

import (
	"os"

	"go.uber.org/zap"
)

var (
	LogFile *os.File
	Logger  *zap.Logger
)
