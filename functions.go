package core

import (
	"bytes"
	"context"

	"v2ray.com/core/common"
	"v2ray.com/core/common/net"
	"v2ray.com/core/features/routing"
)

// CreateObject creates a new object based on the given V2Ray instance and config. The V2Ray instance may be nil.
func CreateObject(v *Instance, config interface{}) (interface{}, error) {
	ctx := context.Background()
	if v != nil {
		ctx = context.WithValue(ctx, v2rayKey, v)
	}
	return common.CreateObject(ctx, config)
}

// StartInstance starts a new V2Ray instance with given serialized config.
// By default V2Ray only support config in protobuf format, i.e., configFormat = "protobuf". Caller need to load other packages to add JSON support.
//
// v2ray:api:stable
func StartInstance(configFormat string, configBytes []byte) (*Instance, error) {
	config, err := LoadConfig(configFormat, "", bytes.NewReader(configBytes))
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

// Dial provides an easy way for upstream caller to create net.Conn through V2Ray.
// It dispatches the request to the given destination by the given V2Ray instance.
// Since it is under a proxy context, the LocalAddr() and RemoteAddr() in returned net.Conn
// will not show real addresses being used for communication.
//
// v2ray:api:stable
func Dial(ctx context.Context, v *Instance, dest net.Destination) (net.Conn, error) {
	dispatcher := v.GetFeature(routing.DispatcherType())
	if dispatcher == nil {
		return nil, newError("routing.Dispatcher is not registered in V2Ray core")
	}
	r, err := dispatcher.(routing.Dispatcher).Dispatch(ctx, dest)
	if err != nil {
		return nil, err
	}
	return net.NewConnection(net.ConnectionInputMulti(r.Writer), net.ConnectionOutputMulti(r.Reader)), nil
}
