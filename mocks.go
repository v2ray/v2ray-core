package core

//go:generate go get -u github.com/golang/mock/gomock
//go:generate go install github.com/golang/mock/mockgen
//go:generate mockgen -package mocks -destination testing/mocks/dns.go -mock_names Client=DNSClient v2ray.com/core/features/dns Client
//go:generate mockgen -package mocks -destination testing/mocks/proxy.go -mock_names Inbound=ProxyInbound,Outbound=ProxyOutbound v2ray.com/core/proxy Inbound,Outbound
