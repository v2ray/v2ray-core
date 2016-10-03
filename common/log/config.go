package log

func (this *Config) Apply() error {
	if this == nil {
		return nil
	}
	if this.AccessLogType == LogType_File {
		if err := InitAccessLogger(this.AccessLogPath); err != nil {
			return err
		}
	}

	if this.ErrorLogType == LogType_None {
		SetLogLevel(LogLevel_Disabled)
	} else {
		if this.ErrorLogType == LogType_File {
			if err := InitErrorLogger(this.ErrorLogPath); err != nil {
				return err
			}
		}
		SetLogLevel(this.ErrorLogLevel)
	}

	return nil
}
