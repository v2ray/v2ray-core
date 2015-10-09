package log

import (
	"log"
	"os"
)

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

func (logger *noOpAccessLogger) Log(from, to string, status AccessStatus, reason string) {
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
}

func (logger *fileAccessLogger) Log(from, to string, status AccessStatus, reason string) {
	logger.queue <- &accessLog{
		From:   from,
		To:     to,
		Status: status,
		Reason: reason,
	}
}

func (logger *fileAccessLogger) Run() {
	for entry := range logger.queue {
		logger.logger.Println(entry.From + " " + string(entry.Status) + " " + entry.To + " " + entry.Reason)
	}
}

func newFileAccessLogger(path string) accessLogger {
	file, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		log.Printf("Unable to create or open file (%s): %v\n", path, err)
		return nil
	}
	return &fileAccessLogger{
		queue:  make(chan *accessLog, 16),
		logger: log.New(file, "", log.Ldate|log.Ltime),
	}
}

var accessLoggerInstance accessLogger = &noOpAccessLogger{}

func InitAccessLogger(file string) {
	logger := newFileAccessLogger(file)
	if logger != nil {
		go logger.(*fileAccessLogger).Run()
		accessLoggerInstance = logger
	}
}

func Access(from, to string, status AccessStatus, reason string) {
	accessLoggerInstance.Log(from, to, status, reason)
}
