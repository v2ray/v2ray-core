package internal

import (
	"context"
	"log"
	"os"

	"v2ray.com/core/common/platform"
)

type LogEntry interface {
	String() string
}

type LogWriter interface {
	Log(LogEntry)
	Close()
}

type StdOutLogWriter struct {
	logger *log.Logger
}

func NewStdOutLogWriter() LogWriter {
	return &StdOutLogWriter{
		logger: log.New(os.Stdout, "", log.Ldate|log.Ltime),
	}
}

func (w *StdOutLogWriter) Log(log LogEntry) {
	w.logger.Print(log.String() + platform.LineSeparator())
}

func (*StdOutLogWriter) Close() {}

type FileLogWriter struct {
	queue  chan string
	logger *log.Logger
	file   *os.File
	ctx    context.Context
	cancel context.CancelFunc
}

func (w *FileLogWriter) Log(log LogEntry) {
	select {
	case <-w.ctx.Done():
		return
	case w.queue <- log.String():
	default:
		// We don't expect this to happen, but don't want to block main thread as well.
	}
}

func (w *FileLogWriter) run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			w.file.Close()
			return
		case entry := <-w.queue:
			w.logger.Print(entry + platform.LineSeparator())
		}
	}
}

func (w *FileLogWriter) Close() {
	w.cancel()
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
