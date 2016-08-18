package internet

import (
	"github.com/v2ray/v2ray-core/common"
	"github.com/v2ray/v2ray-core/common/alloc"
	"github.com/v2ray/v2ray-core/common/loader"
)

type Authenticator interface {
	Seal(*alloc.Buffer)
	Open(*alloc.Buffer) bool
	Overhead() int
}

type AuthenticatorFactory interface {
	Create(AuthenticatorConfig) Authenticator
}

type AuthenticatorConfig interface {
}

var (
	authenticatorCache = make(map[string]AuthenticatorFactory)
	configCache        loader.ConfigLoader
)

func RegisterAuthenticator(name string, factory AuthenticatorFactory, configCreator loader.ConfigCreator) error {
	if _, found := authenticatorCache[name]; found {
		return common.ErrDuplicatedName
	}
	authenticatorCache[name] = factory
	return configCache.RegisterCreator(name, configCreator)
}

func CreateAuthenticator(name string, config AuthenticatorConfig) (Authenticator, error) {
	factory, found := authenticatorCache[name]
	if !found {
		return nil, common.ErrObjectNotFound
	}
	return factory.Create(config), nil
}

func CreateAuthenticatorConfig(rawConfig []byte) (string, AuthenticatorConfig, error) {
	config, name, err := configCache.Load(rawConfig)
	if err != nil {
		return name, nil, err
	}
	return name, config, nil
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
