// Package shadowsocks provides compatible functionality to Shadowsocks.
//
// Shadowsocks client and server are implemented as outbound and inbound respectively in V2Ray's term.
//
// Shadowsocks OTA is fully supported. By default both client and server enable OTA, but it can be optionally disabled.
//
// Supperted Ciphers:
// * AES-256-CFB
// * AES-128-CFB
// * Chacha20
// * Chacha20-IEFT
//
// R.I.P Shadowsocks
package shadowsocks

//go:generate go run $GOPATH/src/v2ray.com/core/common/errors/errorgen/main.go -pkg shadowsocks -path Proxy,Shadowsocks
