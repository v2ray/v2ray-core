package utils

// LogLevel of quic-go
type LogLevel uint8

const (
	// LogLevelNothing disables
	LogLevelNothing LogLevel = iota
	// LogLevelError enables err logs
	LogLevelError
	// LogLevelInfo enables info logs (e.g. packets)
	LogLevelInfo
	// LogLevelDebug enables debug logs (e.g. packet contents)
	LogLevelDebug
)

// A Logger logs.
type Logger interface {
	SetLogLevel(LogLevel)
	SetLogTimeFormat(format string)
	WithPrefix(prefix string) Logger
	Debug() bool

	Errorf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Debugf(format string, args ...interface{})
}

// DefaultLogger is used by quic-go for logging.
var DefaultLogger Logger

type defaultLogger struct{}

var _ Logger = &defaultLogger{}

// SetLogLevel sets the log level
func (l *defaultLogger) SetLogLevel(level LogLevel) {
}

// SetLogTimeFormat sets the format of the timestamp
// an empty string disables the logging of timestamps
func (l *defaultLogger) SetLogTimeFormat(format string) {}

// Debugf logs something
func (l *defaultLogger) Debugf(format string, args ...interface{}) {}

// Infof logs something
func (l *defaultLogger) Infof(format string, args ...interface{}) {}

// Errorf logs something
func (l *defaultLogger) Errorf(format string, args ...interface{}) {}

func (l *defaultLogger) WithPrefix(prefix string) Logger {
	return l
}

// Debug returns true if the log level is LogLevelDebug
func (l *defaultLogger) Debug() bool {
	return false
}

func init() {
	DefaultLogger = &defaultLogger{}
}
