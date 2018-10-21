package core

import (
	"context"
	"sync"

	"v2ray.com/core/common"
	"v2ray.com/core/common/serial"
	"v2ray.com/core/features"
	"v2ray.com/core/features/dns"
	"v2ray.com/core/features/inbound"
	"v2ray.com/core/features/outbound"
	"v2ray.com/core/features/policy"
	"v2ray.com/core/features/routing"
)

// Server is an instance of V2Ray. At any time, there must be at most one Server instance running.
type Server interface {
	common.Runnable
}

// ServerType returns the type of the server.
func ServerType() interface{} {
	return (*Instance)(nil)
}

type resolution struct {
	deps     []interface{}
	callback func([]features.Feature)
}

func getFeature(allFeatures []features.Feature, t interface{}) features.Feature {
	for _, f := range allFeatures {
		if f.Type() == t {
			return f
		}
	}
	return nil
}

func (r *resolution) resolve(allFeatures []features.Feature) bool {
	var fs []features.Feature
	for _, d := range r.deps {
		f := getFeature(allFeatures, d)
		if f == nil {
			return false
		}
		fs = append(fs, f)
	}

	r.callback(fs)
	return true
}

// Instance combines all functionalities in V2Ray.
type Instance struct {
	access             sync.Mutex
	features           []features.Feature
	featureResolutions []resolution
	running            bool
}

func AddInboundHandler(server *Instance, config *InboundHandlerConfig) error {
	inboundManager := server.GetFeature(inbound.ManagerType()).(inbound.Manager)
	rawHandler, err := CreateObject(server, config)
	if err != nil {
		return err
	}
	handler, ok := rawHandler.(inbound.Handler)
	if !ok {
		return newError("not an InboundHandler")
	}
	if err := inboundManager.AddHandler(context.Background(), handler); err != nil {
		return err
	}
	return nil
}

func addInboundHandlers(server *Instance, configs []*InboundHandlerConfig) error {
	for _, inboundConfig := range configs {
		if err := AddInboundHandler(server, inboundConfig); err != nil {
			return err
		}
	}

	return nil
}

func AddOutboundHandler(server *Instance, config *OutboundHandlerConfig) error {
	outboundManager := server.GetFeature(outbound.ManagerType()).(outbound.Manager)
	rawHandler, err := CreateObject(server, config)
	if err != nil {
		return err
	}
	handler, ok := rawHandler.(outbound.Handler)
	if !ok {
		return newError("not an OutboundHandler")
	}
	if err := outboundManager.AddHandler(context.Background(), handler); err != nil {
		return err
	}
	return nil
}

func addOutboundHandlers(server *Instance, configs []*OutboundHandlerConfig) error {
	for _, outboundConfig := range configs {
		if err := AddOutboundHandler(server, outboundConfig); err != nil {
			return err
		}
	}

	return nil
}

// New returns a new V2Ray instance based on given configuration.
// The instance is not started at this point.
// To ensure V2Ray instance works properly, the config must contain one Dispatcher, one InboundHandlerManager and one OutboundHandlerManager. Other features are optional.
func New(config *Config) (*Instance, error) {
	var server = &Instance{}

	if config.Transport != nil {
		features.PrintDeprecatedFeatureWarning("global transport settings")
	}
	if err := config.Transport.Apply(); err != nil {
		return nil, err
	}

	for _, appSettings := range config.App {
		settings, err := appSettings.GetInstance()
		if err != nil {
			return nil, err
		}
		obj, err := CreateObject(server, settings)
		if err != nil {
			return nil, err
		}
		if feature, ok := obj.(features.Feature); ok {
			server.AddFeature(feature)
		}
	}

	if server.GetFeature(dns.ClientType()) == nil {
		server.AddFeature(dns.LocalClient{})
	}

	if server.GetFeature(policy.ManagerType()) == nil {
		server.AddFeature(policy.DefaultManager{})
	}

	if server.GetFeature(routing.RouterType()) == nil {
		server.AddFeature(routing.DefaultRouter{})
	}

	// Add an empty instance at the end, for optional feature requirement.
	server.AddFeature(&Instance{})

	if server.featureResolutions != nil {
		return nil, newError("not all dependency are resolved.")
	}

	if err := addInboundHandlers(server, config.Inbound); err != nil {
		return nil, err
	}

	if err := addOutboundHandlers(server, config.Outbound); err != nil {
		return nil, err
	}

	return server, nil
}

// Type implements common.HasType.
func (s *Instance) Type() interface{} {
	return ServerType()
}

// Close shutdown the V2Ray instance.
func (s *Instance) Close() error {
	s.access.Lock()
	defer s.access.Unlock()

	s.running = false

	var errors []interface{}
	for _, f := range s.features {
		if err := f.Close(); err != nil {
			errors = append(errors, err)
		}
	}
	if len(errors) > 0 {
		return newError("failed to close all features").Base(newError(serial.Concat(errors...)))
	}

	return nil
}

// RequireFeatures registers a callback, which will be called when all dependent features are registered.
func (s *Instance) RequireFeatures(featureTypes []interface{}, callback func([]features.Feature)) {
	r := resolution{
		deps:     featureTypes,
		callback: callback,
	}
	if r.resolve(s.features) {
		return
	}
	s.featureResolutions = append(s.featureResolutions, r)
}

// AddFeature registers a feature into current Instance.
func (s *Instance) AddFeature(feature features.Feature) {
	s.features = append(s.features, feature)

	if s.running {
		if err := feature.Start(); err != nil {
			newError("failed to start feature").Base(err).WriteToLog()
		}
		return
	}

	if s.featureResolutions == nil {
		return
	}

	var pendingResolutions []resolution
	for _, r := range s.featureResolutions {
		if !r.resolve(s.features) {
			pendingResolutions = append(pendingResolutions, r)
		}
	}
	if len(pendingResolutions) == 0 {
		s.featureResolutions = nil
	} else if len(pendingResolutions) < len(s.featureResolutions) {
		s.featureResolutions = pendingResolutions
	}
}

// GetFeature returns a feature of the given type, or nil if such feature is not registered.
func (s *Instance) GetFeature(featureType interface{}) features.Feature {
	return getFeature(s.features, featureType)
}

// Start starts the V2Ray instance, including all registered features. When Start returns error, the state of the instance is unknown.
// A V2Ray instance can be started only once. Upon closing, the instance is not guaranteed to start again.
func (s *Instance) Start() error {
	s.access.Lock()
	defer s.access.Unlock()

	s.running = true
	for _, f := range s.features {
		if err := f.Start(); err != nil {
			return err
		}
	}

	newError("V2Ray ", Version(), " started").AtWarning().WriteToLog()

	return nil
}
