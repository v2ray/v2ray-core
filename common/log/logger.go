package log

import (
	"io"
	"log"
	"os"
	"time"

	"v2ray.com/core/common/platform"
	"v2ray.com/core/common/signal"
)

// LogWriter is the interface for writing logs.
type LogWriter interface {
	Write(string) error
	io.Closer
}

// LogWriterCreator is a function to create LogWriters.
type LogWriterCreator func() LogWriter

type generalLogger struct {
	creator LogWriterCreator
	buffer  chan Message
	access  *signal.Semaphore
}

// NewLogger returns a generic log handler that can handle all type of messages.
func NewLogger(logWriterCreator LogWriterCreator) Handler {
	return &generalLogger{
		creator: logWriterCreator,
		buffer:  make(chan Message, 16),
		access:  signal.NewSemaphore(1),
	}
}

func (l *generalLogger) run() {
	defer l.access.Signal()

	dataWritten := false
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	logger := l.creator()
	if logger == nil {
		return
	}
	defer logger.Close()

	for {
		select {
		case msg := <-l.buffer:
			logger.Write(msg.String() + platform.LineSeparator())
			dataWritten = true
		case <-ticker.C:
			if !dataWritten {
				return
			}
			dataWritten = false
		}
	}
}

func (l *generalLogger) Handle(msg Message) {
	select {
	case l.buffer <- msg:
	default:
	}

	select {
	case <-l.access.Wait():
		go l.run()
	default:
	}
}

type consoleLogWriter struct {
	logger *log.Logger
}

func (w *consoleLogWriter) Write(s string) error {
	w.logger.Print(s)
	return nil
}

func (w *consoleLogWriter) Close() error {
	return nil
}

type fileLogWriter struct {
	file   *os.File
	logger *log.Logger
}

func (w *fileLogWriter) Write(s string) error {
	w.logger.Print(s)
	return nil
}

func (w *fileLogWriter) Close() error {
	return w.file.Close()
}

// CreateStdoutLogWriter returns a LogWriterCreator that creates LogWriter for stdout.
func CreateStdoutLogWriter() LogWriterCreator {
	return func() LogWriter {
		return &consoleLogWriter{
			logger: log.New(os.Stdout, "", log.Ldate|log.Ltime),
		}
	}
}

// CreateFileLogWriter returns a LogWriterCreator that creates LogWriter for the given file.
func CreateFileLogWriter(path string) (LogWriterCreator, error) {
	file, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return nil, err
	}
	file.Close()
	return func() LogWriter {
		file, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
		if err != nil {
			return nil
		}
		return &fileLogWriter{
			file:   file,
			logger: log.New(file, "", log.Ldate|log.Ltime),
		}
	}, nil
}

func init() {
	RegisterHandler(NewLogger(CreateStdoutLogWriter()))
}
