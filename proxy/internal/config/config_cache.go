package config

import (
	"errors"
)

type ConfigObjectCreator func(data []byte) (interface{}, error)

var (
	configCache map[string]ConfigObjectCreator
)

func getConfigKey(protocol string, proxyType string) string {
	return protocol + "_" + proxyType
}

func registerConfigType(protocol string, proxyType string, creator ConfigObjectCreator) error {
	// TODO: check name
	configCache[getConfigKey(protocol, proxyType)] = creator
	return nil
}

func RegisterInboundConfig(protocol string, creator ConfigObjectCreator) error {
	return registerConfigType(protocol, "inbound", creator)
}

func RegisterOutboundConfig(protocol string, creator ConfigObjectCreator) error {
	return registerConfigType(protocol, "outbound", creator)
}

func CreateInboundConnectionConfig(protocol string, data []byte) (interface{}, error) {
	creator, found := configCache[getConfigKey(protocol, "inbound")]
	if !found {
		return nil, errors.New(protocol + " not found.")
	}
	return creator(data)
}

func CreateOutboundConnectionConfig(protocol string, data []byte) (interface{}, error) {
	creator, found := configCache[getConfigKey(protocol, "outbound")]
	if !found {
		return nil, errors.New(protocol + " not found.")
	}
	return creator(data)
}

func initializeConfigCache() {
	configCache = make(map[string]ConfigObjectCreator)
}

func init() {
	initializeConfigCache()
}
