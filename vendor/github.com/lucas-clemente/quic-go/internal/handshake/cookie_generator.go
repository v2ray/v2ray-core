package handshake

import (
	"encoding/asn1"
	"fmt"
	"net"
	"time"

	"github.com/lucas-clemente/quic-go/internal/protocol"
)

const (
	cookiePrefixIP byte = iota
	cookiePrefixString
)

// A Cookie is derived from the client address and can be used to verify the ownership of this address.
type Cookie struct {
	RemoteAddr               string
	OriginalDestConnectionID protocol.ConnectionID
	// The time that the Cookie was issued (resolution 1 second)
	SentTime time.Time
}

// token is the struct that is used for ASN1 serialization and deserialization
type token struct {
	RemoteAddr               []byte
	OriginalDestConnectionID []byte

	Timestamp int64
}

// A CookieGenerator generates Cookies
type CookieGenerator struct {
	cookieProtector cookieProtector
}

// NewCookieGenerator initializes a new CookieGenerator
func NewCookieGenerator() (*CookieGenerator, error) {
	cookieProtector, err := newCookieProtector()
	if err != nil {
		return nil, err
	}
	return &CookieGenerator{
		cookieProtector: cookieProtector,
	}, nil
}

// NewToken generates a new Cookie for a given source address
func (g *CookieGenerator) NewToken(raddr net.Addr, origConnID protocol.ConnectionID) ([]byte, error) {
	data, err := asn1.Marshal(token{
		RemoteAddr:               encodeRemoteAddr(raddr),
		OriginalDestConnectionID: origConnID,
		Timestamp:                time.Now().Unix(),
	})
	if err != nil {
		return nil, err
	}
	return g.cookieProtector.NewToken(data)
}

// DecodeToken decodes a Cookie
func (g *CookieGenerator) DecodeToken(encrypted []byte) (*Cookie, error) {
	// if the client didn't send any Cookie, DecodeToken will be called with a nil-slice
	if len(encrypted) == 0 {
		return nil, nil
	}

	data, err := g.cookieProtector.DecodeToken(encrypted)
	if err != nil {
		return nil, err
	}
	t := &token{}
	rest, err := asn1.Unmarshal(data, t)
	if err != nil {
		return nil, err
	}
	if len(rest) != 0 {
		return nil, fmt.Errorf("rest when unpacking token: %d", len(rest))
	}
	cookie := &Cookie{
		RemoteAddr: decodeRemoteAddr(t.RemoteAddr),
		SentTime:   time.Unix(t.Timestamp, 0),
	}
	if len(t.OriginalDestConnectionID) > 0 {
		cookie.OriginalDestConnectionID = protocol.ConnectionID(t.OriginalDestConnectionID)
	}
	return cookie, nil
}

// encodeRemoteAddr encodes a remote address such that it can be saved in the Cookie
func encodeRemoteAddr(remoteAddr net.Addr) []byte {
	if udpAddr, ok := remoteAddr.(*net.UDPAddr); ok {
		return append([]byte{cookiePrefixIP}, udpAddr.IP...)
	}
	return append([]byte{cookiePrefixString}, []byte(remoteAddr.String())...)
}

// decodeRemoteAddr decodes the remote address saved in the Cookie
func decodeRemoteAddr(data []byte) string {
	// data will never be empty for a Cookie that we generated. Check it to be on the safe side
	if len(data) == 0 {
		return ""
	}
	if data[0] == cookiePrefixIP {
		return net.IP(data[1:]).String()
	}
	return string(data[1:])
}
