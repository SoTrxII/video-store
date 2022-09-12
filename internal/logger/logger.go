package logger

import (
	"github.com/sirupsen/logrus"
	"go.elastic.co/ecslogrus"
)

// BuildLogger Build a new instance
func Build() *logrus.Logger {
	log := logrus.New()
	log.SetFormatter(&ecslogrus.Formatter{})
	return log
}
