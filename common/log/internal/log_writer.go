package internal

import (
	"log"
	"os"

	"github.com/v2ray/v2ray-core/common/platform"
)

type LogWriter interface {
	Log(LogEntry)
}

type NoOpLogWriter struct {
}

func (this *NoOpLogWriter) Log(entry LogEntry) {
	entry.Release()
}

type StdOutLogWriter struct {
	logger *log.Logger
}

func NewStdOutLogWriter() LogWriter {
	return &StdOutLogWriter{
		logger: log.New(os.Stdout, "", log.Ldate|log.Ltime),
	}
}

func (this *StdOutLogWriter) Log(log LogEntry) {
	this.logger.Print(log.String() + platform.LineSeparator())
	log.Release()
}

type FileLogWriter struct {
	queue  chan LogEntry
	logger *log.Logger
	file   *os.File
}

func (this *FileLogWriter) Log(log LogEntry) {
	select {
	case this.queue <- log:
	default:
		log.Release()
		// We don't expect this to happen, but don't want to block main thread as well.
	}
}

func (this *FileLogWriter) run() {
	for {
		entry, open := <-this.queue
		if !open {
			break
		}
		this.logger.Print(entry.String() + platform.LineSeparator())
		entry.Release()
		entry = nil
	}
}

func (this *FileLogWriter) Close() {
	this.file.Close()
}

func NewFileLogWriter(path string) (*FileLogWriter, error) {
	file, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return nil, err
	}
	logger := &FileLogWriter{
		queue:  make(chan LogEntry, 16),
		logger: log.New(file, "", log.Ldate|log.Ltime),
		file:   file,
	}
	go logger.run()
	return logger, nil
}
