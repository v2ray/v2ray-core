package log

// AccessStatus is the status of an access request from clients.
type AccessStatus string

const (
	AccessAccepted = AccessStatus("accepted")
	AccessRejected = AccessStatus("rejected")
)

var (
	accessLoggerInstance logWriter = &noOpLogWriter{}
)

type accessLog struct {
	From   string
	To     string
	Status AccessStatus
	Reason string
}

func (this *accessLog) String() string {
	return this.From + " " + string(this.Status) + " " + this.To + " " + this.Reason
}

// InitAccessLogger initializes the access logger to write into the give file.
func InitAccessLogger(file string) error {
	logger, err := newFileLogWriter(file)
	if err != nil {
		Error("Failed to create access logger on file (", file, "): ", file, err)
		return err
	}
	accessLoggerInstance = logger
	return nil
}

// Access writes an access log.
func Access(from, to string, status AccessStatus, reason string) {
	accessLoggerInstance.Log(&accessLog{
		From:   from,
		To:     to,
		Status: status,
		Reason: reason,
	})
}
