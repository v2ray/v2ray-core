package log

import (
	"log"
	"os"

	"github.com/v2ray/v2ray-core/common/platform"
)

// AccessStatus is the status of an access request from clients.
type AccessStatus string

const (
	AccessAccepted = AccessStatus("accepted")
	AccessRejected = AccessStatus("rejected")
)

type accessLogger interface {
	Log(from, to string, status AccessStatus, reason string)
}

type noOpAccessLogger struct {
}

func (this *noOpAccessLogger) Log(from, to string, status AccessStatus, reason string) {
	// Swallow
}

type accessLog struct {
	From   string
	To     string
	Status AccessStatus
	Reason string
}

type fileAccessLogger struct {
	queue  chan *accessLog
	logger *log.Logger
	file   *os.File
}

func (this *fileAccessLogger) close() {
	this.file.Close()
}

func (logger *fileAccessLogger) Log(from, to string, status AccessStatus, reason string) {
	select {
	case logger.queue <- &accessLog{
		From:   from,
		To:     to,
		Status: status,
		Reason: reason,
	}:
	default:
		// We don't expect this to happen, but don't want to block main thread as well.
	}
}

func (this *fileAccessLogger) Run() {
	for entry := range this.queue {
		this.logger.Println(entry.From + " " + string(entry.Status) + " " + entry.To + " " + entry.Reason)
	}
}

func newFileAccessLogger(path string) accessLogger {
	file, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		log.Printf("Unable to create or open file (%s): %v%s", path, err, platform.LineSeparator())
		return nil
	}
	return &fileAccessLogger{
		queue:  make(chan *accessLog, 16),
		logger: log.New(file, "", log.Ldate|log.Ltime),
		file:   file,
	}
}

var accessLoggerInstance accessLogger = &noOpAccessLogger{}

// InitAccessLogger initializes the access logger to write into the give file.
func InitAccessLogger(file string) {
	logger := newFileAccessLogger(file)
	if logger != nil {
		go logger.(*fileAccessLogger).Run()
		accessLoggerInstance = logger
	}
}

// Access writes an access log.
func Access(from, to string, status AccessStatus, reason string) {
	accessLoggerInstance.Log(from, to, status, reason)
}
