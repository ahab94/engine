package engine

import (
	logs "github.com/sirupsen/logrus"
)

var logger *logs.Logger

func init() {
	SetLogger(logs.New())
}

// SetLogger - sets custom logrus logger
func SetLogger(log *logs.Logger) {
	logger = log
}

func log(id string) *logs.Entry {
	return logger.WithFields(logs.Fields{
		"ID": id,
	})
}
