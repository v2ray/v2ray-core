package log

import (
	"io"
	"log"
	"os"

	"github.com/v2ray/v2ray-core/common/platform"
)

func createLogger(writer io.Writer) *log.Logger {
	return log.New(writer, "", log.Ldate|log.Ltime)
}

type logWriter interface {
	Log(LogEntry)
}

type noOpLogWriter struct {
}

func (this *noOpLogWriter) Log(entry LogEntry) {
	entry.Release()
}

type stdOutLogWriter struct {
	logger *log.Logger
}

func newStdOutLogWriter() logWriter {
	return &stdOutLogWriter{
		logger: createLogger(os.Stdout),
	}
}

func (this *stdOutLogWriter) Log(log LogEntry) {
	this.logger.Print(log.String() + platform.LineSeparator())
	log.Release()
}

type fileLogWriter struct {
	queue  chan LogEntry
	logger *log.Logger
	file   *os.File
}

func (this *fileLogWriter) Log(log LogEntry) {
	select {
	case this.queue <- log:
	default:
		log.Release()
		// We don't expect this to happen, but don't want to block main thread as well.
	}
}

func (this *fileLogWriter) run() {
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

func (this *fileLogWriter) close() {
	this.file.Close()
}

func newFileLogWriter(path string) (*fileLogWriter, error) {
	file, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return nil, err
	}
	logger := &fileLogWriter{
		queue:  make(chan LogEntry, 16),
		logger: log.New(file, "", log.Ldate|log.Ltime),
		file:   file,
	}
	go logger.run()
	return logger, nil
}
