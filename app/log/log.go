package log

//go:generate go run $GOPATH/src/v2ray.com/core/common/errors/errorgen/main.go -pkg log -path App,Log

import (
	"context"
	"sync"

	"v2ray.com/core"
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

	v := core.FromContext(ctx)
	if v != nil {
		common.Must(v.RegisterFeature((*log.Handler)(nil), g))
	}

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

func (*Instance) Type() interface{} {
	return (*Instance)(nil)
}

// Start implements app.Application.Start().
func (g *Instance) Start() error {
	newError("Logger starting").AtDebug().WriteToLog()

	g.Lock()
	defer g.Unlock()

	if g.active {
		return nil
	}

	g.active = true

	if err := g.initAccessLogger(); err != nil {
		return newError("failed to initialize access logger").Base(err).AtWarning()
	}
	if err := g.initErrorLogger(); err != nil {
		return newError("failed to initialize error logger").Base(err).AtWarning()
	}
	log.RegisterHandler(g)

	return nil
}

// Handle implements log.Handler.
func (g *Instance) Handle(msg log.Message) {
	g.RLock()
	defer g.RUnlock()

	if !g.active {
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
	newError("Logger closing").AtDebug().WriteToLog()

	g.Lock()
	defer g.Unlock()

	if !g.active {
		return nil
	}

	g.active = false

	common.Close(g.accessLogger)
	common.Close(g.errorLogger)

	return nil
}

func init() {
	common.Must(common.RegisterConfig((*Config)(nil), func(ctx context.Context, config interface{}) (interface{}, error) {
		return New(ctx, config.(*Config))
	}))
}
