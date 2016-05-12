package log

import (
	"github.com/v2ray/v2ray-core/common/log/internal"
)

type LogLevel int

const (
	DebugLevel   = LogLevel(0)
	InfoLevel    = LogLevel(1)
	WarningLevel = LogLevel(2)
	ErrorLevel   = LogLevel(3)
	NoneLevel    = LogLevel(999)
)

var (
	streamLoggerInstance internal.LogWriter = internal.NewStdOutLogWriter()

	debugLogger   internal.LogWriter = streamLoggerInstance
	infoLogger    internal.LogWriter = streamLoggerInstance
	warningLogger internal.LogWriter = streamLoggerInstance
	errorLogger   internal.LogWriter = streamLoggerInstance
)

func SetLogLevel(level LogLevel) {
	debugLogger = new(internal.NoOpLogWriter)
	if level <= DebugLevel {
		debugLogger = streamLoggerInstance
	}

	infoLogger = new(internal.NoOpLogWriter)
	if level <= InfoLevel {
		infoLogger = streamLoggerInstance
	}

	warningLogger = new(internal.NoOpLogWriter)
	if level <= WarningLevel {
		warningLogger = streamLoggerInstance
	}

	errorLogger = new(internal.NoOpLogWriter)
	if level <= ErrorLevel {
		errorLogger = streamLoggerInstance
	}

	if level == NoneLevel {
		accessLoggerInstance = new(internal.NoOpLogWriter)
	}
}

func InitErrorLogger(file string) error {
	logger, err := internal.NewFileLogWriter(file)
	if err != nil {
		Error("Failed to create error logger on file (", file, "): ", err)
		return err
	}
	streamLoggerInstance = logger
	return nil
}

// Debug outputs a debug log with given format and optional arguments.
func Debug(v ...interface{}) {
	debugLogger.Log(&internal.ErrorLog{
		Prefix: "[Debug]",
		Values: v,
	})
}

// Info outputs an info log with given format and optional arguments.
func Info(v ...interface{}) {
	infoLogger.Log(&internal.ErrorLog{
		Prefix: "[Info]",
		Values: v,
	})
}

// Warning outputs a warning log with given format and optional arguments.
func Warning(v ...interface{}) {
	warningLogger.Log(&internal.ErrorLog{
		Prefix: "[Warning]",
		Values: v,
	})
}

// Error outputs an error log with given format and optional arguments.
func Error(v ...interface{}) {
	errorLogger.Log(&internal.ErrorLog{
		Prefix: "[Error]",
		Values: v,
	})
}
