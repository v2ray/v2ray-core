package log

func (v *Config) Apply() error {
	if v == nil {
		return nil
	}
	if v.AccessLogType == LogType_File {
		if err := InitAccessLogger(v.AccessLogPath); err != nil {
			return err
		}
	}

	if v.ErrorLogType == LogType_None {
		SetLogLevel(LogLevel_Disabled)
	} else {
		if v.ErrorLogType == LogType_File {
			if err := InitErrorLogger(v.ErrorLogPath); err != nil {
				return err
			}
		}
		SetLogLevel(v.ErrorLogLevel)
	}

	return nil
}
