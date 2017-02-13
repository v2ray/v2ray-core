package internal

import (
	"context"
	"log"
	"os"

	"v2ray.com/core/common/platform"
)

type LogWriter interface {
	Log(LogEntry)
	Close()
}

type NoOpLogWriter struct {
}

func (v *NoOpLogWriter) Log(entry LogEntry) {}

func (v *NoOpLogWriter) Close() {
}

type StdOutLogWriter struct {
	logger *log.Logger
}

func NewStdOutLogWriter() LogWriter {
	return &StdOutLogWriter{
		logger: log.New(os.Stdout, "", log.Ldate|log.Ltime),
	}
}

func (v *StdOutLogWriter) Log(log LogEntry) {
	v.logger.Print(log.String() + platform.LineSeparator())
}

func (v *StdOutLogWriter) Close() {}

type FileLogWriter struct {
	queue  chan string
	logger *log.Logger
	file   *os.File
	ctx    context.Context
	cancel context.CancelFunc
}

func (v *FileLogWriter) Log(log LogEntry) {
	select {
	case <-v.ctx.Done():
		return
	case v.queue <- log.String():
	default:
		// We don't expect this to happen, but don't want to block main thread as well.
	}
}

func (v *FileLogWriter) run(ctx context.Context) {
L:
	for {
		select {
		case <-ctx.Done():
			break L
		case entry := <-v.queue:
			v.logger.Print(entry + platform.LineSeparator())
		}
	}
	v.file.Close()
}

func (v *FileLogWriter) Close() {
	v.cancel()
}

func NewFileLogWriter(path string) (*FileLogWriter, error) {
	file, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithCancel(context.Background())
	logger := &FileLogWriter{
		queue:  make(chan string, 16),
		logger: log.New(file, "", log.Ldate|log.Ltime),
		file:   file,
		ctx:    ctx,
		cancel: cancel,
	}
	go logger.run(ctx)
	return logger, nil
}
