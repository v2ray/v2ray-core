package core

//go:generate go run github.com/golang/mock/mockgen -package mocks -destination testing/mocks/io.go -mock_names Reader=Reader,Writer=Writer io Reader,Writer
//go:generate go run github.com/golang/mock/mockgen -package mocks -destination testing/mocks/log.go -mock_names Handler=LogHandler v2ray.com/core/common/log Handler
//go:generate go run github.com/golang/mock/mockgen -package mocks -destination testing/mocks/mux.go -mock_names ClientWorkerFactory=MuxClientWorkerFactory v2ray.com/core/common/mux ClientWorkerFactory
//go:generate go run github.com/golang/mock/mockgen -package mocks -destination testing/mocks/dns.go -mock_names Client=DNSClient v2ray.com/core/features/dns Client
//go:generate go run github.com/golang/mock/mockgen -package mocks -destination testing/mocks/outbound.go -mock_names Manager=OutboundManager,HandlerSelector=OutboundHandlerSelector v2ray.com/core/features/outbound Manager,HandlerSelector
//go:generate go run github.com/golang/mock/mockgen -package mocks -destination testing/mocks/proxy.go -mock_names Inbound=ProxyInbound,Outbound=ProxyOutbound v2ray.com/core/proxy Inbound,Outbound
