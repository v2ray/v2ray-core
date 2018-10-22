package core

//go:generate go get -u github.com/golang/mock/gomock
//go:generate go install github.com/golang/mock/mockgen
//go:generate mockgen -package mocks -destination v2ray.com/core/features/mocks/dns.go v2ray.com/core/features/dns Client
