package log

import (
	"v2ray.com/core/app/log/internal"
	"v2ray.com/core/common/errors"
)

// AccessStatus is the status of an access request from clients.
type AccessStatus string

const (
	AccessAccepted = AccessStatus("accepted")
	AccessRejected = AccessStatus("rejected")
)

var (
	accessLoggerInstance internal.LogWriter = new(internal.NoOpLogWriter)
)

// InitAccessLogger initializes the access logger to write into the give file.
func InitAccessLogger(file string) error {
	logger, err := internal.NewFileLogWriter(file)
	if err != nil {
		return errors.Base(err).Message("Failed to create access logger on file: ", file)
	}
	accessLoggerInstance = logger
	return nil
}

// Access writes an access log.
func Access(from, to interface{}, status AccessStatus, reason interface{}) {
	accessLoggerInstance.Log(&internal.AccessLog{
		From:   from,
		To:     to,
		Status: string(status),
		Reason: reason,
	})
}
