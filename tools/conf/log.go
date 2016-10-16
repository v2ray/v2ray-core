package conf

import (
	"strings"

	"v2ray.com/core/common/log"
)

type LogConfig struct {
	AccessLog string `json:"access"`
	ErrorLog  string `json:"error"`
	LogLevel  string `json:"loglevel"`
}

func (this *LogConfig) Build() *log.Config {
	if this == nil {
		return nil
	}
	config := new(log.Config)
	if len(this.AccessLog) > 0 {
		config.AccessLogPath = this.AccessLog
		config.AccessLogType = log.LogType_File
	}
	if len(this.ErrorLog) > 0 {
		config.ErrorLogPath = this.ErrorLog
		config.ErrorLogType = log.LogType_File
	}

	level := strings.ToLower(this.LogLevel)
	switch level {
	case "debug":
		config.ErrorLogLevel = log.LogLevel_Debug
	case "info":
		config.ErrorLogLevel = log.LogLevel_Info
	case "error":
		config.ErrorLogLevel = log.LogLevel_Error
	case "none":
		config.ErrorLogType = log.LogType_None
	default:
		config.ErrorLogLevel = log.LogLevel_Warning
	}
	return config
}
