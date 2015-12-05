package log

import (
	"fmt"
)

const (
	DebugLevel   = LogLevel(0)
	InfoLevel    = LogLevel(1)
	WarningLevel = LogLevel(2)
	ErrorLevel   = LogLevel(3)
)

type errorLog struct {
	prefix string
	format string
	values []interface{}
}

func (this *errorLog) String() string {
	var data string
	if len(this.values) == 0 {
		data = this.format
	} else {
		data = fmt.Sprintf(this.format, this.values...)
	}
	return this.prefix + data
}

var (
	noOpLoggerInstance   logWriter = &noOpLogWriter{}
	streamLoggerInstance logWriter = newStdOutLogWriter()

	debugLogger   = noOpLoggerInstance
	infoLogger    = noOpLoggerInstance
	warningLogger = streamLoggerInstance
	errorLogger   = streamLoggerInstance
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

func InitErrorLogger(file string) error {
	logger, err := newFileLogWriter(file)
	if err != nil {
		Error("Failed to create error logger on file (%s): %v", file, err)
		return err
	}
	streamLoggerInstance = logger
	return nil
}

// Debug outputs a debug log with given format and optional arguments.
func Debug(format string, v ...interface{}) {
	debugLogger.Log(&errorLog{
		prefix: "[Debug]",
		format: format,
		values: v,
	})
}

// Info outputs an info log with given format and optional arguments.
func Info(format string, v ...interface{}) {
	infoLogger.Log(&errorLog{
		prefix: "[Info]",
		format: format,
		values: v,
	})
}

// Warning outputs a warning log with given format and optional arguments.
func Warning(format string, v ...interface{}) {
	warningLogger.Log(&errorLog{
		prefix: "[Warning]",
		format: format,
		values: v,
	})
}

// Error outputs an error log with given format and optional arguments.
func Error(format string, v ...interface{}) {
	errorLogger.Log(&errorLog{
		prefix: "[Error]",
		format: format,
		values: v,
	})
}
