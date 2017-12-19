package log

//go:generate go run $GOPATH/src/v2ray.com/core/common/errors/errorgen/main.go -pkg log -path App,Log

import (
	"context"
	"sync"

	"v2ray.com/core/app/log/internal"
	"v2ray.com/core/common"
	"v2ray.com/core/common/log"
)

type Instance struct {
	sync.RWMutex
	config       *Config
	accessLogger internal.LogWriter
	errorLogger  internal.LogWriter
}

func New(ctx context.Context, config *Config) (*Instance, error) {
	return &Instance{
		config: config,
	}, nil
}

func (*Instance) Interface() interface{} {
	return (*Instance)(nil)
}

func (g *Instance) initAccessLogger() error {
	switch g.config.AccessLogType {
	case LogType_File:
		logger, err := internal.NewFileLogWriter(g.config.AccessLogPath)
		if err != nil {
			return err
		}
		g.accessLogger = logger
	case LogType_Console:
		g.accessLogger = internal.NewStdOutLogWriter()
	default:
	}
	return nil
}

func (g *Instance) initErrorLogger() error {
	switch g.config.ErrorLogType {
	case LogType_File:
		logger, err := internal.NewFileLogWriter(g.config.ErrorLogPath)
		if err != nil {
			return err
		}
		g.errorLogger = logger
	case LogType_Console:
		g.errorLogger = internal.NewStdOutLogWriter()
	default:
	}
	return nil
}

func (g *Instance) Start() error {
	if err := g.initAccessLogger(); err != nil {
		return newError("failed to initialize access logger").Base(err).AtWarning()
	}
	if err := g.initErrorLogger(); err != nil {
		return newError("failed to initialize error logger").Base(err).AtWarning()
	}
	log.RegisterHandler(g)
	return nil
}

func (g *Instance) Handle(msg log.Message) {
	switch msg := msg.(type) {
	case *log.AccessMessage:
		g.RLock()
		defer g.RUnlock()
		if g.accessLogger != nil {
			g.accessLogger.Log(msg)
		}
	case *log.GeneralMessage:
		if msg.Severity.SevererThan(g.config.ErrorLogLevel) {
			g.RLock()
			defer g.RUnlock()
			if g.errorLogger != nil {
				g.errorLogger.Log(msg)
			}
		}
	default:
		// Swallow
	}
}

func (g *Instance) Close() {
	g.Lock()
	defer g.Unlock()

	g.accessLogger.Close()
	g.accessLogger = nil

	g.errorLogger.Close()
	g.errorLogger = nil
}

func init() {
	common.Must(common.RegisterConfig((*Config)(nil), func(ctx context.Context, config interface{}) (interface{}, error) {
		return New(ctx, config.(*Config))
	}))
}
