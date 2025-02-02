package logger

import (
	"github.com/sirupsen/logrus"
	"os"
)

var log = logrus.New()

// init инициализирует логгер
func init() {
	// устанавливем уровень логирования
	log.SetLevel(logrus.InfoLevel)
	// устанавливем вывод логов в файл
	file, err := os.OpenFile("server.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		log.SetOutput(file)
	} else {
		log.Info("logrus.go/init - failed while opening log file")
	}
}

func Fatal(args ...interface{}) {
	log.Fatal(args)
}

func Warn(args ...interface{}) {
	log.Warn(args)
}

func Info(args ...interface{}) {
	log.Info(args)
}

func Infof(format string, args ...interface{}) {
	log.Infof(format, args)
}
