package log

import (
	"v2ray.com/core/common/errors"
	"v2ray.com/core/common/log/internal"
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
	if level >= LogLevel_Debug {
		debugLogger = streamLoggerInstance
	}

	infoLogger = new(internal.NoOpLogWriter)
	if level >= LogLevel_Info {
		infoLogger = streamLoggerInstance
	}

	warningLogger = new(internal.NoOpLogWriter)
	if level >= LogLevel_Warning {
		warningLogger = streamLoggerInstance
	}

	errorLogger = new(internal.NoOpLogWriter)
	if level >= LogLevel_Error {
		errorLogger = streamLoggerInstance
	}
}

func InitErrorLogger(file string) error {
	logger, err := internal.NewFileLogWriter(file)
	if err != nil {
		return errors.Base(err).Message("Log: Failed to create error logger on file (", file, ")")
	}
	streamLoggerInstance = logger
	return nil
}

// Debug outputs a debug log with given format and optional arguments.
func Debug(val ...interface{}) {
	debugLogger.Log(&internal.ErrorLog{
		Prefix: "[Debug]",
		Values: val,
	})
}

// Info outputs an info log with given format and optional arguments.
func Info(val ...interface{}) {
	infoLogger.Log(&internal.ErrorLog{
		Prefix: "[Info]",
		Values: val,
	})
}

// Warning outputs a warning log with given format and optional arguments.
func Warning(val ...interface{}) {
	warningLogger.Log(&internal.ErrorLog{
		Prefix: "[Warning]",
		Values: val,
	})
}

// Error outputs an error log with given format and optional arguments.
func Error(val ...interface{}) {
	errorLogger.Log(&internal.ErrorLog{
		Prefix: "[Error]",
		Values: val,
	})
}

func Close() {
	streamLoggerInstance.Close()
	accessLoggerInstance.Close()
}
