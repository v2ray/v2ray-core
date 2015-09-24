package protocol

import (
	_ "fmt"

	"github.com/v2ray/v2ray-core/common/errors"
	_ "github.com/v2ray/v2ray-core/common/log"
)

type SocksVersion4Error struct {
	errors.ErrorCode
}

var socksVersion4ErrorInstance = SocksVersion4Error{ErrorCode: 1000}

func NewSocksVersion4Error() SocksVersion4Error {
	return socksVersion4ErrorInstance
}

func (err SocksVersion4Error) Error() string {
	return err.Prefix() + "Request is socks version 4."
}

type Socks4AuthenticationRequest struct {
	Version byte
	Command byte
	Port    uint16
	IP      [4]byte
}
