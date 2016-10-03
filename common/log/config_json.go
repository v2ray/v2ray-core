package log

import (
	"encoding/json"
	"errors"
	"strings"
)

func (this *Config) UnmarshalJSON(data []byte) error {
	type JsonLogConfig struct {
		AccessLog string `json:"access"`
		ErrorLog  string `json:"error"`
		LogLevel  string `json:"loglevel"`
	}
	jsonConfig := new(JsonLogConfig)
	if err := json.Unmarshal(data, jsonConfig); err != nil {
		return errors.New("Log: Failed to parse log config: " + err.Error())
	}
	if len(jsonConfig.AccessLog) > 0 {
		this.AccessLogPath = jsonConfig.AccessLog
		this.AccessLogType = LogType_File
	}
	if len(jsonConfig.ErrorLog) > 0 {
		this.ErrorLogPath = jsonConfig.ErrorLog
		this.ErrorLogType = LogType_File
	}

	level := strings.ToLower(jsonConfig.LogLevel)
	switch level {
	case "debug":
		this.ErrorLogLevel = LogLevel_Debug
	case "info":
		this.ErrorLogLevel = LogLevel_Info
	case "error":
		this.ErrorLogLevel = LogLevel_Error
	case "none":
		this.ErrorLogType = LogType_None
	default:
		this.ErrorLogLevel = LogLevel_Warning
	}
	return nil
}
