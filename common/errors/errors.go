package errors

import (
	"fmt"
)

type AuthenticationError struct {
	AuthDetail interface{}
}

func NewAuthenticationError(detail interface{}) AuthenticationError {
	return AuthenticationError{
		AuthDetail: detail,
	}
}

func (err AuthenticationError) Error() string {
	return fmt.Sprintf("[Error 0x0001] Invalid auth %v", err.AuthDetail)
}

type ProtocolVersionError struct {
	Version int
}

func NewProtocolVersionError(version int) ProtocolVersionError {
	return ProtocolVersionError{
		Version: version,
	}
}

func (err ProtocolVersionError) Error() string {
	return fmt.Sprintf("[Error 0x0002] Invalid version %d", err.Version)
}

type CorruptedPacketError struct {
}

func NewCorruptedPacketError() CorruptedPacketError {
	return CorruptedPacketError{}
}

func (err CorruptedPacketError) Error() string {
	return "[Error 0x0003] Corrupted packet."
}

type IPFormatError struct {
	IP []byte
}

func NewIPFormatError(ip []byte) IPFormatError {
	return IPFormatError{
		IP: ip,
	}
}

func (err IPFormatError) Error() string {
	return fmt.Sprintf("[Error 0x0004] Invalid IP %v", err.IP)
}
