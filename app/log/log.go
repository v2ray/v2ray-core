package log

//go:generate go run $GOPATH/src/v2ray.com/core/common/errors/errorgen/main.go -pkg log -path App,Log

import (
	"context"
	"sync"

	"v2ray.com/core/common"
	"v2ray.com/core/common/log"
)

// Instance is an app.Application that handles logs.
type Instance struct {
	sync.RWMutex
	config       *Config
	accessLogger log.Handler
	errorLogger  log.Handler
	active       bool
}

// New creates a new log.Instance based on the given config.
func New(ctx context.Context, config *Config) (*Instance, error) {
	g := &Instance{
		config: config,
		active: true,
	}

	if err := g.initAccessLogger(); err != nil {
		return nil, newError("failed to initialize access logger").Base(err).AtWarning()
	}
	if err := g.initErrorLogger(); err != nil {
		return nil, newError("failed to initialize error logger").Base(err).AtWarning()
	}
	log.RegisterHandler(g)

	return g, nil
}

func (g *Instance) initAccessLogger() error {
	switch g.config.AccessLogType {
	case LogType_File:
		creator, err := log.CreateFileLogWriter(g.config.AccessLogPath)
		if err != nil {
			return err
		}
		g.accessLogger = log.NewLogger(creator)
	case LogType_Console:
		g.accessLogger = log.NewLogger(log.CreateStdoutLogWriter())
	default:
	}
	return nil
}

func (g *Instance) initErrorLogger() error {
	switch g.config.ErrorLogType {
	case LogType_File:
		creator, err := log.CreateFileLogWriter(g.config.ErrorLogPath)
		if err != nil {
			return err
		}
		g.errorLogger = log.NewLogger(creator)
	case LogType_Console:
		g.errorLogger = log.NewLogger(log.CreateStdoutLogWriter())
	default:
	}
	return nil
}

// Start implements app.Application.Start().
func (g *Instance) Start() error {
	g.Lock()
	defer g.Unlock()
	g.active = true
	return nil
}

func (g *Instance) isActive() bool {
	g.RLock()
	defer g.RUnlock()

	return g.active
}

// Handle implements log.Handler.
func (g *Instance) Handle(msg log.Message) {
	if !g.isActive() {
		return
	}

	switch msg := msg.(type) {
	case *log.AccessMessage:
		if g.accessLogger != nil {
			g.accessLogger.Handle(msg)
		}
	case *log.GeneralMessage:
		if g.errorLogger != nil && msg.Severity <= g.config.ErrorLogLevel {
			g.errorLogger.Handle(msg)
		}
	default:
		// Swallow
	}
}

// Close implement app.Application.Close().
func (g *Instance) Close() error {
	g.Lock()
	defer g.Unlock()

	g.active = false

	return nil
}

func init() {
	common.Must(common.RegisterConfig((*Config)(nil), func(ctx context.Context, config interface{}) (interface{}, error) {
		return New(ctx, config.(*Config))
	}))
}
