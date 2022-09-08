package logger

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"path"
)

func init() {
	fileName := path.Join("logs", "verkaufsautomat.log")
	// Create the log file if doesn't exist. And append to it if it already exists.
	f, err := os.OpenFile(fileName, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	Formatter := new(log.TextFormatter)

	Formatter.TimestampFormat = "02-01-2006 15:04:05"
	Formatter.FullTimestamp = true
	log.SetFormatter(Formatter)
	if err != nil {

		fmt.Println(err)
	} else {
		log.SetOutput(f)
	}

}

func Info(msg interface{}) {
	log.Info(msg)
}

func Error(msg interface{}) {
	log.Error(msg)
}
