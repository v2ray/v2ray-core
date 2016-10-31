package internet

import (
	"io"
	"v2ray.com/core/common"
)

type ConnectionAuthenticator interface {
	Seal(io.Writer) io.Writer
	Open(io.Reader) (io.Reader, error)
}

type ConnectionAuthenticatorFactory interface {
	Create(interface{}) ConnectionAuthenticator
}

var (
	connectionAuthenticatorCache = make(map[string]ConnectionAuthenticatorFactory)
)

func RegisterConnectionAuthenticator(name string, factory ConnectionAuthenticatorFactory) error {
	if _, found := connectionAuthenticatorCache[name]; found {
		return common.ErrDuplicatedName
	}
	connectionAuthenticatorCache[name] = factory
	return nil
}

func CreateConnectionAuthenticator(name string, config interface{}) (ConnectionAuthenticator, error) {
	factory, found := connectionAuthenticatorCache[name]
	if !found {
		return nil, common.ErrObjectNotFound
	}
	return factory.Create(config), nil
}
