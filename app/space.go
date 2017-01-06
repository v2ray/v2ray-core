package app

import (
	"github.com/golang/protobuf/proto"
	"v2ray.com/core/common/errors"
	"v2ray.com/core/common/log"
	"v2ray.com/core/common/serial"
)

type Application interface {
}

type InitializationCallback func() error

type ApplicationFactory interface {
	Create(space Space, config interface{}) (Application, error)
}

type AppGetter interface {
	GetApp(name string) Application
}

var (
	applicationFactoryCache = make(map[string]ApplicationFactory)
)

func RegisterApplicationFactory(defaultConfig proto.Message, factory ApplicationFactory) error {
	if defaultConfig == nil {
		return errors.New("Space: config is nil.")
	}
	name := serial.GetMessageType(defaultConfig)
	if len(name) == 0 {
		return errors.New("Space: cannot get config type.")
	}
	applicationFactoryCache[name] = factory
	return nil
}

// A Space contains all apps that may be available in a V2Ray runtime.
// Caller must check the availability of an app by calling HasXXX before getting its instance.
type Space interface {
	AddApp(config proto.Message) error
	AddAppLegacy(name string, app Application)
	Initialize() error
	OnInitialize(InitializationCallback)
}

type spaceImpl struct {
	initialized bool
	cache       map[string]Application
	appInit     []InitializationCallback
}

func NewSpace() Space {
	return &spaceImpl{
		cache:   make(map[string]Application),
		appInit: make([]InitializationCallback, 0, 32),
	}
}

func (v *spaceImpl) OnInitialize(f InitializationCallback) {
	if v.initialized {
		if err := f(); err != nil {
			log.Error("Space: error after space initialization: ", err)
		}
	} else {
		v.appInit = append(v.appInit, f)
	}
}

func (v *spaceImpl) Initialize() error {
	for _, f := range v.appInit {
		if err := f(); err != nil {
			return err
		}
	}
	v.appInit = nil
	v.initialized = true
	return nil
}

func (v *spaceImpl) GetApp(configType string) Application {
	obj, found := v.cache[configType]
	if !found {
		return nil
	}
	return obj
}

func (v *spaceImpl) AddApp(config proto.Message) error {
	configName := serial.GetMessageType(config)
	factory, found := applicationFactoryCache[configName]
	if !found {
		return errors.New("Space: app not registered: ", configName)
	}
	app, err := factory.Create(v, config)
	if err != nil {
		return err
	}
	v.cache[configName] = app
	return nil
}

func (v *spaceImpl) AddAppLegacy(name string, application Application) {
	v.cache[name] = application
}
