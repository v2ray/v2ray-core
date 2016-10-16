package app

import (
	"errors"

	"v2ray.com/core/common"
)

var (
	ErrMissingApplication = errors.New("App: Failed to found one or more applications.")
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
	applicationFactoryCache map[string]ApplicationFactory
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

func (this *spaceImpl) InitializeApplication(f ApplicationInitializer) {
	this.appInit = append(this.appInit, f)
}

func (this *spaceImpl) Initialize() error {
	for _, f := range this.appInit {
		err := f()
		if err != nil {
			return err
		}
	}
	return nil
}

func (this *spaceImpl) HasApp(id ID) bool {
	_, found := this.cache[id]
	return found
}

func (this *spaceImpl) GetApp(id ID) Application {
	obj, found := this.cache[id]
	if !found {
		return nil
	}
	return obj
}

func (this *spaceImpl) BindApp(id ID, application Application) {
	this.cache[id] = application
}

func (this *spaceImpl) BindFromConfig(name string, config interface{}) error {
	factory, found := applicationFactoryCache[name]
	if !found {
		return errors.New("Space: app not registered: " + name)
	}
	app, err := factory.Create(this, config)
	if err != nil {
		return err
	}
	this.BindApp(factory.AppId(), app)
	return nil
}
