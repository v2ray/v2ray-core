package internet

import (
	"v2ray.com/core/common"
	"v2ray.com/core/common/alloc"
)

type Authenticator interface {
	Seal(*alloc.Buffer)
	Open(*alloc.Buffer) bool
	Overhead() int
}

type AuthenticatorFactory interface {
	Create(interface{}) Authenticator
}

var (
	authenticatorCache = make(map[string]AuthenticatorFactory)
)

func RegisterAuthenticator(name string, factory AuthenticatorFactory) error {
	if _, found := authenticatorCache[name]; found {
		return common.ErrDuplicatedName
	}
	authenticatorCache[name] = factory
	return nil
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

func (v *AuthenticatorChain) Overhead() int {
	total := 0
	for _, auth := range v.authenticators {
		total += auth.Overhead()
	}
	return total
}

func (v *AuthenticatorChain) Open(payload *alloc.Buffer) bool {
	for _, auth := range v.authenticators {
		if !auth.Open(payload) {
			return false
		}
	}
	return true
}

func (v *AuthenticatorChain) Seal(payload *alloc.Buffer) {
	for i := len(v.authenticators) - 1; i >= 0; i-- {
		auth := v.authenticators[i]
		auth.Seal(payload)
	}
}
