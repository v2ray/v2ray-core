package core

//go:generate go get -u github.com/golang/mock/gomock
//go:generate go install github.com/golang/mock/mockgen
//go:generate mockgen -package mocks -destination v2ray.com/core/testing/mocks/dns.go -mock_names Client=MockDNSClient v2ray.com/core/features/dns Client
//go:generate mockgen -package mocks -destination v2ray.com/core/testing/mocks/proxy.go -mock_names Inbound=MockProxyInbound,Outbound=MockProxyOutbound v2ray.com/core/proxy Inbound,Outbound
