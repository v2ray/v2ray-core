package conf

import (
	"strings"

	"v2ray.com/core/app/log"
	clog "v2ray.com/core/common/log"
)

func DefaultLogConfig() *log.Config {
	return &log.Config{
		AccessLogType: log.LogType_None,
		ErrorLogType:  log.LogType_Console,
		ErrorLogLevel: clog.Severity_Warning,
	}
}

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

	if v.AccessLog == "none" {
		config.AccessLogType = log.LogType_None
	} else if len(v.AccessLog) > 0 {
		config.AccessLogPath = v.AccessLog
		config.AccessLogType = log.LogType_File
	}
	if v.ErrorLog == "none" {
		config.ErrorLogType = log.LogType_None
	} else if len(v.ErrorLog) > 0 {
		config.ErrorLogPath = v.ErrorLog
		config.ErrorLogType = log.LogType_File
	}

	level := strings.ToLower(v.LogLevel)
	switch level {
	case "debug":
		config.ErrorLogLevel = clog.Severity_Debug
	case "info":
		config.ErrorLogLevel = clog.Severity_Info
	case "error":
		config.ErrorLogLevel = clog.Severity_Error
	case "none":
		config.ErrorLogType = log.LogType_None
		config.AccessLogType = log.LogType_None
	default:
		config.ErrorLogLevel = clog.Severity_Warning
	}
	return config
}
