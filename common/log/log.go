package log

import (
	"fmt"

	"github.com/v2ray/v2ray-core/common/serial"
)

const (
	DebugLevel   = LogLevel(0)
	InfoLevel    = LogLevel(1)
	WarningLevel = LogLevel(2)
	ErrorLevel   = LogLevel(3)
)

type errorLog struct {
	prefix string
	values []interface{}
}

func (this *errorLog) String() string {
	data := ""
	for _, value := range this.values {
		switch typedVal := value.(type) {
		case string:
			data += typedVal
		case *string:
			data += *typedVal
		case serial.String:
			data += typedVal.String()
		case error:
			data += typedVal.Error()
		default:
			data += fmt.Sprintf("%v", value)
		}
	}
	return this.prefix + data
}

var (
	noOpLoggerInstance   logWriter = &noOpLogWriter{}
	streamLoggerInstance logWriter = newStdOutLogWriter()

	debugLogger   = streamLoggerInstance
	infoLogger    = streamLoggerInstance
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
		Error("Failed to create error logger on file (", file, "): ", err)
		return err
	}
	streamLoggerInstance = logger
	return nil
}

// Debug outputs a debug log with given format and optional arguments.
func Debug(v ...interface{}) {
	debugLogger.Log(&errorLog{
		prefix: "[Debug]",
		values: v,
	})
}

// Info outputs an info log with given format and optional arguments.
func Info(v ...interface{}) {
	infoLogger.Log(&errorLog{
		prefix: "[Info]",
		values: v,
	})
}

// Warning outputs a warning log with given format and optional arguments.
func Warning(v ...interface{}) {
	warningLogger.Log(&errorLog{
		prefix: "[Warning]",
		values: v,
	})
}

// Error outputs an error log with given format and optional arguments.
func Error(v ...interface{}) {
	errorLogger.Log(&errorLog{
		prefix: "[Error]",
		values: v,
	})
}
