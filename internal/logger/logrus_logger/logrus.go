package logrus_logger

import (
	"github.com/sirupsen/logrus"
	"os"
)

var Log = logrus.New()

func init() {
	// устанавливем уровень логирования
	Log.SetLevel(logrus.InfoLevel)
	// устанавливем вывод логов в файл
	file, err := os.OpenFile("server.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		Log.SetOutput(file)
	} else {
		Log.Info("internal/logger/logrus/logrus.go - failed while opening log file")
	}
}
