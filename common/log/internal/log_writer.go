package internal

import (
	"log"
	"os"
	"time"

	"v2ray.com/core/common/platform"
	"v2ray.com/core/common/signal"
)

type LogWriter interface {
	Log(LogEntry)
	Close()
}

type NoOpLogWriter struct {
}

func (v *NoOpLogWriter) Log(entry LogEntry) {
	entry.Release()
}

func (v *NoOpLogWriter) Close() {
}

type StdOutLogWriter struct {
	logger *log.Logger
	cancel *signal.CancelSignal
}

func NewStdOutLogWriter() LogWriter {
	return &StdOutLogWriter{
		logger: log.New(os.Stdout, "", log.Ldate|log.Ltime),
		cancel: signal.NewCloseSignal(),
	}
}

func (v *StdOutLogWriter) Log(log LogEntry) {
	v.logger.Print(log.String() + platform.LineSeparator())
	log.Release()
}

func (v *StdOutLogWriter) Close() {
	time.Sleep(500 * time.Millisecond)
}

type FileLogWriter struct {
	queue  chan string
	logger *log.Logger
	file   *os.File
	cancel *signal.CancelSignal
}

func (v *FileLogWriter) Log(log LogEntry) {
	select {
	case v.queue <- log.String():
	default:
		// We don't expect this to happen, but don't want to block main thread as well.
	}
	log.Release()
}

func (v *FileLogWriter) run() {
	v.cancel.WaitThread()
	defer v.cancel.FinishThread()

	for {
		entry, open := <-v.queue
		if !open {
			break
		}
		v.logger.Print(entry + platform.LineSeparator())
	}
}

func (v *FileLogWriter) Close() {
	close(v.queue)
	v.cancel.WaitForDone()
	v.file.Close()
}

func NewFileLogWriter(path string) (*FileLogWriter, error) {
	file, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return nil, err
	}
	logger := &FileLogWriter{
		queue:  make(chan string, 16),
		logger: log.New(file, "", log.Ldate|log.Ltime),
		file:   file,
		cancel: signal.NewCloseSignal(),
	}
	go logger.run()
	return logger, nil
}
