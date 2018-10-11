package core

import (
	"context"

	"v2ray.com/core/common"
	"v2ray.com/core/common/buf"
)

// CreateObject creates a new object based on the given V2Ray instance and config. The V2Ray instance may be nil.
func CreateObject(v *Instance, config interface{}) (interface{}, error) {
	ctx := context.Background()
	if v != nil {
		ctx = context.WithValue(ctx, v2rayKey, v)
	}
	return common.CreateObject(ctx, config)
}

// StartInstance starts a new V2Ray instance with given serialized config, and return a handle for shutting down the instance.
func StartInstance(configFormat string, configBytes []byte) (common.Closable, error) {
	var mb buf.MultiBuffer
	common.Must2(mb.Write(configBytes))
	config, err := LoadConfig(configFormat, "", &mb)
	if err != nil {
		return nil, err
	}
	instance, err := New(config)
	if err != nil {
		return nil, err
	}
	if err := instance.Start(); err != nil {
		return nil, err
	}
	return instance, nil
}

func PrintDeprecatedFeatureWarning(feature string) {
	newError("You are using a deprecated feature: " + feature + ". Please update your config file with latest configuration format, or update your client software.").WriteToLog()
}
