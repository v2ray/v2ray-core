package conf

import (
	"strings"

	"v2ray.com/core/app/log"
)

type LogConfig struct {
	AccessLog string `json:"access"`
	ErrorLog  string `json:"error"`
	LogLevel  string `json:"loglevel"`
}

func (v *LogConfig) Build() *log.Config {
	if v == nil {
		return nil
	}
	config := &log.Config{
		ErrorLogType:  log.LogType_Console,
		AccessLogType: log.LogType_Console,
	}

	if len(v.AccessLog) > 0 {
		config.AccessLogPath = v.AccessLog
		config.AccessLogType = log.LogType_File
	}
	if len(v.ErrorLog) > 0 {
		config.ErrorLogPath = v.ErrorLog
		config.ErrorLogType = log.LogType_File
	}

	level := strings.ToLower(v.LogLevel)
	switch level {
	case "debug":
		config.ErrorLogLevel = log.LogLevel_Debug
	case "info":
		config.ErrorLogLevel = log.LogLevel_Info
	case "error":
		config.ErrorLogLevel = log.LogLevel_Error
	case "none":
		config.ErrorLogType = log.LogType_None
		config.AccessLogType = log.LogType_None
	default:
		config.ErrorLogLevel = log.LogLevel_Warning
	}
	return config
}
