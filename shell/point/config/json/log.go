package json

type LogConfig struct {
	AccessLogValue string `json:"access"`
}

func (config *LogConfig) AccessLog() string {
	return config.AccessLogValue
}
