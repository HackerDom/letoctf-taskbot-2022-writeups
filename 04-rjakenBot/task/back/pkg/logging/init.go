package logging

import (
	"github.com/sirupsen/logrus"
	"os"
)

func Init() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.DebugLevel)
}
