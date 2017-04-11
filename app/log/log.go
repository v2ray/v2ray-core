package log

//go:generate go run $GOPATH/src/v2ray.com/core/tools/generrorgen/main.go -pkg log -path App,Log

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
		return newError("failed to create error logger on file (", file, ")").Base(err)
	}
	streamLoggerInstance = logger
	return nil
}

func getLoggerAndPrefix(s errors.Severity) (internal.LogWriter, string) {
	switch s {
	case errors.SeverityDebug:
		return debugLogger, "[Debug]"
	case errors.SeverityInfo:
		return infoLogger, "[Info]"
	case errors.SeverityWarning:
		return warningLogger, "[Warning]"
	case errors.SeverityError:
		return errorLogger, "[Error]"
	default:
		return infoLogger, "[Info]"
	}
}

// Trace logs an error message based on its severity.
func Trace(err error) {
	logger, prefix := getLoggerAndPrefix(errors.GetSeverity(err))
	logger.Log(&internal.ErrorLog{
		Prefix: prefix,
		Error:  err,
	})
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
