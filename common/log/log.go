package log

import (
	"fmt"
	"io"
	"os"
)

const (
	DebugLevel   = LogLevel(0)
	InfoLevel    = LogLevel(1)
	WarningLevel = LogLevel(2)
	ErrorLevel   = LogLevel(3)
)

type logger interface {
	WriteLog(prefix, format string, v ...interface{})
}

type noOpLogger struct {
}

func (l *noOpLogger) WriteLog(prefix, format string, v ...interface{}) {
	// Swallow
}

type streamLogger struct {
	writer io.Writer
}

func (l *streamLogger) WriteLog(prefix, format string, v ...interface{}) {
	var data string
	if v == nil || len(v) == 0 {
		data = format
	} else {
		data = fmt.Sprintf(format, v...)
	}
	l.writer.Write([]byte(prefix + data))
	l.writer.Write([]byte{'\n'})
}

var (
	noOpLoggerInstance   logger = &noOpLogger{}
	streamLoggerInstance logger = &streamLogger{
		writer: os.Stdout,
	}

	debugLogger   = noOpLoggerInstance
	infoLogger    = noOpLoggerInstance
	warningLogger = noOpLoggerInstance
	errorLogger   = noOpLoggerInstance
)

type LogLevel int

func SetLogLevel(level LogLevel) {
	debugLogger = noOpLoggerInstance
	if level <= DebugLevel {
		debugLogger = streamLoggerInstance
	}

	infoLogger = noOpLoggerInstance
	if level <= InfoLevel {
		infoLogger = streamLoggerInstance
	}

	warningLogger = noOpLoggerInstance
	if level <= WarningLevel {
		warningLogger = streamLoggerInstance
	}

	errorLogger = noOpLoggerInstance
	if level <= ErrorLevel {
		errorLogger = streamLoggerInstance
	}
}

func Debug(format string, v ...interface{}) {
	debugLogger.WriteLog("[Debug]", format, v...)
}

func Info(format string, v ...interface{}) {
	infoLogger.WriteLog("[Info]", format, v...)
}

func Warning(format string, v ...interface{}) {
	warningLogger.WriteLog("[Warning]", format, v...)
}

func Error(format string, v ...interface{}) {
	errorLogger.WriteLog("[Error]", format, v...)
}
