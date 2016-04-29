package log

import (
	"fmt"

	"github.com/v2ray/v2ray-core/common"
	"github.com/v2ray/v2ray-core/common/alloc"
	"github.com/v2ray/v2ray-core/common/serial"
)

const (
	DebugLevel   = LogLevel(0)
	InfoLevel    = LogLevel(1)
	WarningLevel = LogLevel(2)
	ErrorLevel   = LogLevel(3)
	NoneLevel    = LogLevel(999)
)

type LogEntry interface {
	common.Releasable
	serial.String
}

type errorLog struct {
	prefix string
	values []interface{}
}

func (this *errorLog) Release() {
	for index := range this.values {
		this.values[index] = nil
	}
	this.values = nil
}

func (this *errorLog) String() string {
	b := alloc.NewSmallBuffer().Clear()
	defer b.Release()

	b.AppendString(this.prefix)

	for _, value := range this.values {
		switch typedVal := value.(type) {
		case string:
			b.AppendString(typedVal)
		case *string:
			b.AppendString(*typedVal)
		case serial.String:
			b.AppendString(typedVal.String())
		case error:
			b.AppendString(typedVal.Error())
		default:
			b.AppendString(fmt.Sprint(value))
		}
	}
	return b.String()
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
