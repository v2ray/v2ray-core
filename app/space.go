package app

import (
	"v2ray.com/core/common"
	"v2ray.com/core/common/errors"
)

type ID int

// Context of a function call from proxy to app.
type Context interface {
	CallerTag() string
}

type Caller interface {
	Tag() string
}

type Application interface {
	common.Releasable
}

type ApplicationInitializer func() error
type ApplicationFactory interface {
	Create(space Space, config interface{}) (Application, error)
	AppId() ID
}

var (
	applicationFactoryCache = make(map[string]ApplicationFactory)
)

func RegisterApplicationFactory(name string, factory ApplicationFactory) error {
	applicationFactoryCache[name] = factory
	return nil
}

// A Space contains all apps that may be available in a V2Ray runtime.
// Caller must check the availability of an app by calling HasXXX before getting its instance.
type Space interface {
	Initialize() error
	InitializeApplication(ApplicationInitializer)

	HasApp(ID) bool
	GetApp(ID) Application
	BindApp(ID, Application)
	BindFromConfig(name string, config interface{}) error
}

type spaceImpl struct {
	cache   map[ID]Application
	appInit []ApplicationInitializer
}

func NewSpace() Space {
	return &spaceImpl{
		cache:   make(map[ID]Application),
		appInit: make([]ApplicationInitializer, 0, 32),
	}
}

func (v *spaceImpl) InitializeApplication(f ApplicationInitializer) {
	v.appInit = append(v.appInit, f)
}

func (v *spaceImpl) Initialize() error {
	for _, f := range v.appInit {
		err := f()
		if err != nil {
			return err
		}
	}
	return nil
}

func (v *spaceImpl) HasApp(id ID) bool {
	_, found := v.cache[id]
	return found
}

func (v *spaceImpl) GetApp(id ID) Application {
	obj, found := v.cache[id]
	if !found {
		return nil
	}
	return obj
}

func (v *spaceImpl) BindApp(id ID, application Application) {
	v.cache[id] = application
}

func (v *spaceImpl) BindFromConfig(name string, config interface{}) error {
	factory, found := applicationFactoryCache[name]
	if !found {
		return errors.New("Space: app not registered: ", name)
	}
	app, err := factory.Create(v, config)
	if err != nil {
		return err
	}
	v.BindApp(factory.AppId(), app)
	return nil
}
