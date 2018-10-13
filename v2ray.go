package core

import (
	"context"
	"sync"

	"v2ray.com/core/common"
	"v2ray.com/core/common/serial"
	"v2ray.com/core/common/uuid"
	"v2ray.com/core/features"
	"v2ray.com/core/features/dns"
	"v2ray.com/core/features/inbound"
	"v2ray.com/core/features/outbound"
	"v2ray.com/core/features/policy"
	"v2ray.com/core/features/routing"
	"v2ray.com/core/features/stats"
)

// Server is an instance of V2Ray. At any time, there must be at most one Server instance running.
// Deprecated. Use Instance directly.
type Server interface {
	common.Runnable
}

// Instance combines all functionalities in V2Ray.
type Instance struct {
	dnsClient     syncDNSClient
	policyManager syncPolicyManager
	dispatcher    syncDispatcher
	router        syncRouter
	ihm           syncInboundHandlerManager
	ohm           syncOutboundHandlerManager
	stats         syncStatManager

	access   sync.Mutex
	features []features.Feature
	id       uuid.UUID
	running  bool
}

// New returns a new V2Ray instance based on given configuration.
// The instance is not started at this point.
// To ensure V2Ray instance works properly, the config must contain one Dispatcher, one InboundHandlerManager and one OutboundHandlerManager. Other features are optional.
func New(config *Config) (*Instance, error) {
	var server = &Instance{
		id: uuid.New(),
	}

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
		if _, err := CreateObject(server, settings); err != nil {
			return nil, err
		}
	}

	for _, inboundConfig := range config.Inbound {
		rawHandler, err := CreateObject(server, inboundConfig)
		if err != nil {
			return nil, err
		}
		handler, ok := rawHandler.(inbound.Handler)
		if !ok {
			return nil, newError("not an InboundHandler")
		}
		if err := server.InboundHandlerManager().AddHandler(context.Background(), handler); err != nil {
			return nil, err
		}
	}

	for _, outboundConfig := range config.Outbound {
		rawHandler, err := CreateObject(server, outboundConfig)
		if err != nil {
			return nil, err
		}
		handler, ok := rawHandler.(outbound.Handler)
		if !ok {
			return nil, newError("not an OutboundHandler")
		}
		if err := server.OutboundHandlerManager().AddHandler(context.Background(), handler); err != nil {
			return nil, err
		}
	}

	return server, nil
}

// ID returns a unique ID for this V2Ray instance.
func (s *Instance) ID() uuid.UUID {
	return s.id
}

// Close shutdown the V2Ray instance.
func (s *Instance) Close() error {
	s.access.Lock()
	defer s.access.Unlock()

	s.running = false

	var errors []interface{}
	for _, f := range s.allFeatures() {
		if err := f.Close(); err != nil {
			errors = append(errors, err)
		}
	}
	if len(errors) > 0 {
		return newError("failed to close all features").Base(newError(serial.Concat(errors...)))
	}

	return nil
}

// Start starts the V2Ray instance, including all registered features. When Start returns error, the state of the instance is unknown.
// A V2Ray instance can be started only once. Upon closing, the instance is not guaranteed to start again.
func (s *Instance) Start() error {
	s.access.Lock()
	defer s.access.Unlock()

	s.running = true
	for _, f := range s.allFeatures() {
		if err := f.Start(); err != nil {
			return err
		}
	}

	newError("V2Ray ", Version(), " started").AtWarning().WriteToLog()

	return nil
}

// RegisterFeature registers the given feature into V2Ray.
// If feature is one of the following types, the corresponding feature in this Instance
// will be replaced: DNSClient, PolicyManager, Router, Dispatcher, InboundHandlerManager, OutboundHandlerManager.
func (s *Instance) RegisterFeature(instance features.Feature) error {
	running := false

	switch instance.Type().(type) {
	case dns.Client, *dns.Client:
		s.dnsClient.Set(instance.(dns.Client))
	case policy.Manager, *policy.Manager:
		s.policyManager.Set(instance.(policy.Manager))
	case routing.Router, *routing.Router:
		s.router.Set(instance.(routing.Router))
	case routing.Dispatcher, *routing.Dispatcher:
		s.dispatcher.Set(instance.(routing.Dispatcher))
	case inbound.Manager, *inbound.Manager:
		s.ihm.Set(instance.(inbound.Manager))
	case outbound.Manager, *outbound.Manager:
		s.ohm.Set(instance.(outbound.Manager))
	case stats.Manager, *stats.Manager:
		s.stats.Set(instance.(stats.Manager))
	default:
		s.access.Lock()
		s.features = append(s.features, instance)
		running = s.running
		s.access.Unlock()
	}

	if running {
		return instance.Start()
	}
	return nil
}

func (s *Instance) allFeatures() []features.Feature {
	return append([]features.Feature{s.DNSClient(), s.PolicyManager(), s.Dispatcher(), s.Router(), s.InboundHandlerManager(), s.OutboundHandlerManager(), s.Stats()}, s.features...)
}

// GetFeature returns a feature that was registered in this Instance. Nil if not found.
// The returned Feature must implement common.HasType and whose type equals to the given feature type.
func (s *Instance) GetFeature(featureType interface{}) features.Feature {
	switch featureType.(type) {
	case dns.Client, *dns.Client:
		return s.DNSClient()
	case policy.Manager, *policy.Manager:
		return s.PolicyManager()
	case routing.Router, *routing.Router:
		return s.Router()
	case routing.Dispatcher, *routing.Dispatcher:
		return s.Dispatcher()
	case inbound.Manager, *inbound.Manager:
		return s.InboundHandlerManager()
	case outbound.Manager, *outbound.Manager:
		return s.OutboundHandlerManager()
	case stats.Manager, *stats.Manager:
		return s.Stats()
	default:
		for _, f := range s.features {
			if f.Type() == featureType {
				return f
			}
		}
		return nil
	}
}

// DNSClient returns the dns.Client used by this Instance. The returned dns.Client is always functional.
func (s *Instance) DNSClient() dns.Client {
	return &(s.dnsClient)
}

// PolicyManager returns the policy.Manager used by this Instance. The returned policy.Manager is always functional.
func (s *Instance) PolicyManager() policy.Manager {
	return &(s.policyManager)
}

// Router returns the Router used by this Instance. The returned Router is always functional.
func (s *Instance) Router() routing.Router {
	return &(s.router)
}

// Dispatcher returns the Dispatcher used by this Instance. If Dispatcher was not registered before, the returned value doesn't work, although it is not nil.
func (s *Instance) Dispatcher() routing.Dispatcher {
	return &(s.dispatcher)
}

// InboundHandlerManager returns the InboundHandlerManager used by this Instance. If InboundHandlerManager was not registered before, the returned value doesn't work.
func (s *Instance) InboundHandlerManager() inbound.Manager {
	return &(s.ihm)
}

// OutboundHandlerManager returns the OutboundHandlerManager used by this Instance. If OutboundHandlerManager was not registered before, the returned value doesn't work.
func (s *Instance) OutboundHandlerManager() outbound.Manager {
	return &(s.ohm)
}

// Stats returns the stats.Manager used by this Instance. If StatManager was not registered before, the returned value doesn't work.
func (s *Instance) Stats() stats.Manager {
	return &(s.stats)
}
