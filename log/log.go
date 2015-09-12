package log

import (
	"errors"
	"fmt"
	"log"
)

const (
	DebugLevel   = LogLevel(0)
	InfoLevel    = LogLevel(1)
	WarningLevel = LogLevel(2)
	ErrorLevel   = LogLevel(3)
)

var logLevel = WarningLevel

type LogLevel int

func SetLogLevel(level LogLevel) {
	logLevel = level
}

func writeLog(data string, level LogLevel) {
	if level < logLevel {
		return
	}
	log.Print(data)
}

func Debug(format string, v ...interface{}) {
	data := fmt.Sprintf(format, v)
	writeLog("[Debug]"+data, DebugLevel)
}

func Info(format string, v ...interface{}) {
	data := fmt.Sprintf(format, v)
	writeLog("[Info]"+data, InfoLevel)
}

func Warning(format string, v ...interface{}) {
	data := fmt.Sprintf(format, v)
	writeLog("[Warning]"+data, WarningLevel)
}

func Error(format string, v ...interface{}) error {
	data := fmt.Sprintf(format, v)
	writeLog("[Error]"+data, ErrorLevel)
	return errors.New(data)
}
