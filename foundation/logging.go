package foundation

import (
	"errors"
	"os"
	"strconv"

	"github.com/sirupsen/logrus"
)

func NewLogger() *logrus.Logger {
	return &logrus.Logger{
		Out:       os.Stdout,
		Formatter: new(logrus.JSONFormatter),
		Hooks:     make(logrus.LevelHooks),
		Level:     logrus.DebugLevel,
	}
}

var (
	errEmptyPort           = errors.New("restPort number empty")
	errPortNotAValidNumber = errors.New("restPort number is not a valid number")
	errPortReserved        = errors.New("restPort is a reserved number")
	errPortBeyondRange     = errors.New("restPort is beyond the allowed range")
)

func ValidatePort(port string) error {
	if port == "" {
		return errEmptyPort
	}

	addr, err := strconv.ParseInt(port, 10, 64)
	if err != nil {
		return errPortNotAValidNumber
	}

	if addr < 1024 {
		return errPortReserved
	}

	if addr > 65535 {
		return errPortBeyondRange
	}
	return nil
}
