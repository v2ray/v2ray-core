package log

import (
	"context"

	"v2ray.com/core/app"
	"v2ray.com/core/app/log/internal"
	"v2ray.com/core/common"
	"v2ray.com/core/common/errors"
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
		return errors.New("failed to create error logger on file (", file, ")").Base(err).Path("App", "Log")
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

func Trace(err error) {
	s := errors.GetSeverity(err)
	switch s {
	case errors.SeverityDebug:
		Debug(err)
	case errors.SeverityInfo:
		Info(err)
	case errors.SeverityWarning:
		Warning(err)
	case errors.SeverityError:
		Error(err)
	default:
		Info(err)
	}
}

type Instance struct {
	config *Config
}

func New(ctx context.Context, config *Config) (*Instance, error) {
	return &Instance{config: config}, nil
}

func (*Instance) Interface() interface{} {
	return (*Instance)(nil)
}

func (g *Instance) Start() error {
	config := g.config
	if config.AccessLogType == LogType_File {
		if err := InitAccessLogger(config.AccessLogPath); err != nil {
			return err
		}
	}

	if config.ErrorLogType == LogType_None {
		SetLogLevel(LogLevel_Disabled)
	} else {
		if config.ErrorLogType == LogType_File {
			if err := InitErrorLogger(config.ErrorLogPath); err != nil {
				return err
			}
		}
		SetLogLevel(config.ErrorLogLevel)
	}

	return nil
}

func (*Instance) Close() {
	streamLoggerInstance.Close()
	accessLoggerInstance.Close()
}

func FromSpace(space app.Space) *Instance {
	v := space.GetApplication((*Instance)(nil))
	if logger, ok := v.(*Instance); ok && logger != nil {
		return logger
	}
	return nil
}

func init() {
	common.Must(common.RegisterConfig((*Config)(nil), func(ctx context.Context, config interface{}) (interface{}, error) {
		return New(ctx, config.(*Config))
	}))
}
