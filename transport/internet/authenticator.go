package internet

import (
	"errors"

	"v2ray.com/core/common"
	"v2ray.com/core/common/alloc"
	"v2ray.com/core/common/loader"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
)

type Authenticator interface {
	Seal(*alloc.Buffer)
	Open(*alloc.Buffer) bool
	Overhead() int
}

type AuthenticatorFactory interface {
	Create(interface{}) Authenticator
}

func (this *AuthenticatorConfig) GetInternalConfig() (interface{}, error) {
	config, err := configCache.CreateConfig(this.Name)
	if err != nil {
		return nil, err
	}
	if err := ptypes.UnmarshalAny(this.Settings, config.(proto.Message)); err != nil {
		return nil, err
	}
	return config, nil
}

func NewAuthenticatorConfig(name string, config interface{}) (*AuthenticatorConfig, error) {
	pbMsg, ok := config.(proto.Message)
	if !ok {
		return nil, errors.New("Internet|Authenticator: Failed to convert config into proto message.")
	}
	anyConfig, err := ptypes.MarshalAny(pbMsg)
	if err != nil {
		return nil, err
	}
	return &AuthenticatorConfig{
		Name:     name,
		Settings: anyConfig,
	}, nil
}

func (this *AuthenticatorConfig) CreateAuthenticator() (Authenticator, error) {
	config, err := this.GetInternalConfig()
	if err != nil {
		return nil, err
	}
	return CreateAuthenticator(this.Name, config)
}

var (
	authenticatorCache = make(map[string]AuthenticatorFactory)
	configCache        = loader.ConfigCreatorCache{}
)

func RegisterAuthenticator(name string, factory AuthenticatorFactory) error {
	if _, found := authenticatorCache[name]; found {
		return common.ErrDuplicatedName
	}
	authenticatorCache[name] = factory
	return nil
}

func RegisterAuthenticatorConfig(name string, configCreator loader.ConfigCreator) error {
	return configCache.RegisterCreator(name, configCreator)
}

func CreateAuthenticator(name string, config interface{}) (Authenticator, error) {
	factory, found := authenticatorCache[name]
	if !found {
		return nil, common.ErrObjectNotFound
	}
	return factory.Create(config), nil
}

type AuthenticatorChain struct {
	authenticators []Authenticator
}

func NewAuthenticatorChain(auths ...Authenticator) Authenticator {
	return &AuthenticatorChain{
		authenticators: auths,
	}
}

func (this *AuthenticatorChain) Overhead() int {
	total := 0
	for _, auth := range this.authenticators {
		total += auth.Overhead()
	}
	return total
}

func (this *AuthenticatorChain) Open(payload *alloc.Buffer) bool {
	for _, auth := range this.authenticators {
		if !auth.Open(payload) {
			return false
		}
	}
	return true
}

func (this *AuthenticatorChain) Seal(payload *alloc.Buffer) {
	for i := len(this.authenticators) - 1; i >= 0; i-- {
		auth := this.authenticators[i]
		auth.Seal(payload)
	}
}
