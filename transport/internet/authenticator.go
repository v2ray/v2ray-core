package internet

import (
	"errors"

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
	ErrDuplicatedAuthenticator = errors.New("Authenticator already registered.")
	ErrAuthenticatorNotFound   = errors.New("Authenticator not found.")

	authenticatorCache = make(map[string]AuthenticatorFactory)
	configCache        loader.ConfigLoader
)

func RegisterAuthenticator(name string, factory AuthenticatorFactory, configCreator loader.ConfigCreator) error {
	if _, found := authenticatorCache[name]; found {
		return ErrDuplicatedAuthenticator
	}
	authenticatorCache[name] = factory
	return configCache.RegisterCreator(name, configCreator)
}

func CreateAuthenticator(name string, config AuthenticatorConfig) (Authenticator, error) {
	factory, found := authenticatorCache[name]
	if !found {
		return nil, ErrAuthenticatorNotFound
	}
	return factory.Create(config.(AuthenticatorConfig)), nil
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

type NoOpAuthenticator struct{}

func (this NoOpAuthenticator) Overhead() int {
	return 0
}
func (this NoOpAuthenticator) Open(payload *alloc.Buffer) bool {
	return true
}
func (this NoOpAuthenticator) Seal(payload *alloc.Buffer) {}

type NoOpAuthenticatorFactory struct{}

func (this NoOpAuthenticatorFactory) Create(config AuthenticatorConfig) Authenticator {
	return NoOpAuthenticator{}
}

type NoOpAuthenticatorConfig struct{}

func init() {
	RegisterAuthenticator("none", NoOpAuthenticatorFactory{}, func() interface{} { return NoOpAuthenticatorConfig{} })
}
