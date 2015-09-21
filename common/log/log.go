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

func writeLog(level LogLevel, prefix, format string, v ...interface{}) string {
	if level < logLevel {
		return ""
	}
	var data string
	if v == nil || len(v) == 0 {
		data = format
	} else {
		data = fmt.Sprintf(format, v...)
	}
	log.Println(prefix + data)
	return data
}

func Debug(format string, v ...interface{}) {
	writeLog(DebugLevel, "[Debug]", format, v...)
}

func Info(format string, v ...interface{}) {
	writeLog(InfoLevel, "[Info]", format, v...)
}

func Warning(format string, v ...interface{}) {
	writeLog(WarningLevel, "[Warning]", format, v...)
}

func Error(format string, v ...interface{}) error {
	data := writeLog(ErrorLevel, "[Error]", format, v...)
	return errors.New(data)
}
