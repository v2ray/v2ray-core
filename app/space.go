package app

import (
	"context"
	"reflect"

	"v2ray.com/core/common"
	"v2ray.com/core/common/event"
)

type Application interface {
	Interface() interface{}
	Start() error
	Close()
}

func CreateAppFromConfig(ctx context.Context, config interface{}) (Application, error) {
	application, err := common.CreateObject(ctx, config)
	if err != nil {
		return nil, err
	}
	switch a := application.(type) {
	case Application:
		return a, nil
	default:
		return nil, newError("not an application")
	}
}

// A Space contains all apps that may be available in a V2Ray runtime.
type Space interface {
	event.Registry
	GetApplication(appInterface interface{}) Application
	AddApplication(application Application) error
	Initialize() error
	Start() error
	Close()
}

const (
	// SpaceInitializing is an event to be fired when Space is being initialized.
	SpaceInitializing event.Event = iota
)

type spaceImpl struct {
	event.Listener
	cache       map[reflect.Type]Application
	initialized bool
}

// NewSpace creates a new Space.
func NewSpace() Space {
	return &spaceImpl{
		cache: make(map[reflect.Type]Application),
	}
}

func (s *spaceImpl) On(e event.Event, h event.Handler) {
	if e == SpaceInitializing && s.initialized {
		_ = h(nil) // Ignore error
		return
	}
	s.Listener.On(e, h)
}

func (s *spaceImpl) Initialize() error {
	if s.initialized {
		return nil
	}
	s.initialized = true
	return s.Fire(SpaceInitializing, nil)
}

func (s *spaceImpl) GetApplication(appInterface interface{}) Application {
	if s == nil {
		return nil
	}
	appType := reflect.TypeOf(appInterface)
	return s.cache[appType]
}

func (s *spaceImpl) AddApplication(app Application) error {
	if s == nil {
		return newError("nil space").AtError()
	}
	appType := reflect.TypeOf(app.Interface())
	s.cache[appType] = app
	return nil
}

func (s *spaceImpl) Start() error {
	for _, app := range s.cache {
		if err := app.Start(); err != nil {
			return err
		}
	}
	return nil
}

func (s *spaceImpl) Close() {
	for _, app := range s.cache {
		app.Close()
	}
}

type contextKey int

const (
	spaceKey = contextKey(0)
)

func AddApplicationToSpace(ctx context.Context, appConfig interface{}) error {
	space := SpaceFromContext(ctx)
	if space == nil {
		return newError("no space in context").AtError()
	}
	application, err := CreateAppFromConfig(ctx, appConfig)
	if err != nil {
		return err
	}
	return space.AddApplication(application)
}

func SpaceFromContext(ctx context.Context) Space {
	return ctx.Value(spaceKey).(Space)
}

func ContextWithSpace(ctx context.Context, space Space) context.Context {
	return context.WithValue(ctx, spaceKey, space)
}
