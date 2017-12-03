// Package kcp - A Fast and Reliable ARQ Protocol
//
// Acknowledgement:
//    skywind3000@github for inventing the KCP protocol
//    xtaci@github for translating to Golang
package kcp

//go:generate go run $GOPATH/src/v2ray.com/core/common/errors/errorgen/main.go -pkg kcp -path Transport,Internet,mKCP
