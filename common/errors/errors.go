package errors

import (
	"fmt"
)

func HasCode(err error, code int) bool {
	if errWithCode, ok := err.(ErrorWithCode); ok {
		return errWithCode.Code() == code
	}
	return false
}

type ErrorWithCode interface {
	error
	Code() int
}

type ErrorCode int

func (code ErrorCode) Code() int {
	return int(code)
}

func (code ErrorCode) Prefix() string {
	return fmt.Sprintf("[Error 0x%04X] ", code.Code())
}

type AuthenticationError struct {
	ErrorCode
	AuthDetail interface{}
}

func NewAuthenticationError(detail interface{}) AuthenticationError {
	return AuthenticationError{
		ErrorCode:  1,
		AuthDetail: detail,
	}
}

func (err AuthenticationError) Error() string {
	return fmt.Sprintf("%sInvalid auth %v", err.Prefix(), err.AuthDetail)
}

type ProtocolVersionError struct {
	ErrorCode
	Version int
}

func NewProtocolVersionError(version int) ProtocolVersionError {
	return ProtocolVersionError{
		ErrorCode: 2,
		Version:   version,
	}
}

func (err ProtocolVersionError) Error() string {
	return fmt.Sprintf("%sInvalid version %d", err.Prefix(), err.Version)
}

type CorruptedPacketError struct {
	ErrorCode
}

var corruptedPacketErrorInstance = CorruptedPacketError{ErrorCode: 3}

func NewCorruptedPacketError() CorruptedPacketError {
	return corruptedPacketErrorInstance
}

func (err CorruptedPacketError) Error() string {
	return err.Prefix() + "Corrupted packet."
}

type IPFormatError struct {
	ErrorCode
	IP []byte
}

func NewIPFormatError(ip []byte) IPFormatError {
	return IPFormatError{
		ErrorCode: 4,
		IP:        ip,
	}
}

func (err IPFormatError) Error() string {
	return fmt.Sprintf("%sInvalid IP %v", err.Prefix(), err.IP)
}

type ConfigurationError struct {
	ErrorCode
}

var configurationErrorInstance = ConfigurationError{ErrorCode: 5}

func NewConfigurationError() ConfigurationError {
	return configurationErrorInstance
}

func (r ConfigurationError) Error() string {
	return r.Prefix() + "Invalid configuration."
}

type InvalidOperationError struct {
	ErrorCode
	Operation string
}

func NewInvalidOperationError(operation string) InvalidOperationError {
	return InvalidOperationError{
		ErrorCode: 6,
		Operation: operation,
	}
}

func (r InvalidOperationError) Error() string {
	return r.Prefix() + "Invalid operation: " + r.Operation
}
