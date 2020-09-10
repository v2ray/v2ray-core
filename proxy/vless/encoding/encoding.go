// +build !confonly

package encoding

//go:generate errorgen

import (
	"io"

	"v2ray.com/core/common/buf"
	"v2ray.com/core/common/net"
	"v2ray.com/core/common/protocol"
	"v2ray.com/core/proxy/vless"
)

const (
	Version = byte(0)
)

var addrParser = protocol.NewAddressParser(
	protocol.AddressFamilyByte(byte(protocol.AddressTypeIPv4), net.AddressFamilyIPv4),
	protocol.AddressFamilyByte(byte(protocol.AddressTypeDomain), net.AddressFamilyDomain),
	protocol.AddressFamilyByte(byte(protocol.AddressTypeIPv6), net.AddressFamilyIPv6),
	protocol.PortThenAddress(),
)

// EncodeRequestHeader writes encoded request header into the given writer.
func EncodeRequestHeader(writer io.Writer, request *protocol.RequestHeader, requestAddons *Addons) error {

	buffer := buf.StackNew()
	defer buffer.Release()

	if err := buffer.WriteByte(request.Version); err != nil {
		return newError("failed to write request version").Base(err)
	}

	if _, err := buffer.Write(request.User.Account.(*vless.MemoryAccount).ID.Bytes()); err != nil {
		return newError("failed to write request user id").Base(err)
	}

	if err := EncodeHeaderAddons(&buffer, requestAddons); err != nil {
		return newError("failed to encode request header addons").Base(err)
	}

	if err := buffer.WriteByte(byte(request.Command)); err != nil {
		return newError("failed to write request command").Base(err)
	}

	if request.Command != protocol.RequestCommandMux {
		if err := addrParser.WriteAddressPort(&buffer, request.Address, request.Port); err != nil {
			return newError("failed to write request address and port").Base(err)
		}
	}

	if _, err := writer.Write(buffer.Bytes()); err != nil {
		return newError("failed to write request header").Base(err)
	}

	return nil
}

// DecodeRequestHeader decodes and returns (if successful) a RequestHeader from an input stream.
func DecodeRequestHeader(reader io.Reader, validator *vless.Validator) (*protocol.RequestHeader, *Addons, error, *buf.Buffer) {

	buffer := buf.StackNew()
	defer buffer.Release()

	pre := buf.New()

	if _, err := buffer.ReadFullFrom(reader, 1); err != nil {
		pre.Write(buffer.Bytes())
		return nil, nil, newError("failed to read request version").Base(err), pre
	}

	request := &protocol.RequestHeader{
		Version: buffer.Byte(0),
	}

	pre.Write(buffer.Bytes())

	switch request.Version {
	case 0:

		buffer.Clear()
		if _, err := buffer.ReadFullFrom(reader, protocol.IDBytesLen); err != nil {
			pre.Write(buffer.Bytes())
			return nil, nil, newError("failed to read request user id").Base(err), pre
		}

		var id [16]byte
		copy(id[:], buffer.Bytes())

		if request.User = validator.Get(id); request.User == nil {
			pre.Write(buffer.Bytes())
			return nil, nil, newError("invalid request user id"), pre
		}

		requestAddons, err := DecodeHeaderAddons(&buffer, reader)
		if err != nil {
			return nil, nil, newError("failed to decode request header addons").Base(err), nil
		}

		buffer.Clear()
		if _, err := buffer.ReadFullFrom(reader, 1); err != nil {
			return nil, nil, newError("failed to read request command").Base(err), nil
		}

		request.Command = protocol.RequestCommand(buffer.Byte(0))
		switch request.Command {
		case protocol.RequestCommandMux:
			request.Address = net.DomainAddress("v1.mux.cool")
			request.Port = 0
		case protocol.RequestCommandTCP, protocol.RequestCommandUDP:
			if addr, port, err := addrParser.ReadAddressPort(&buffer, reader); err == nil {
				request.Address = addr
				request.Port = port
			}
		}

		if request.Address == nil {
			return nil, nil, newError("invalid request address"), nil
		}

		return request, requestAddons, nil, nil

	default:

		return nil, nil, newError("invalid request version"), pre

	}

}

// EncodeResponseHeader writes encoded response header into the given writer.
func EncodeResponseHeader(writer io.Writer, request *protocol.RequestHeader, responseAddons *Addons) error {

	buffer := buf.StackNew()
	defer buffer.Release()

	if err := buffer.WriteByte(request.Version); err != nil {
		return newError("failed to write response version").Base(err)
	}

	if err := EncodeHeaderAddons(&buffer, responseAddons); err != nil {
		return newError("failed to encode response header addons").Base(err)
	}

	if _, err := writer.Write(buffer.Bytes()); err != nil {
		return newError("failed to write response header").Base(err)
	}

	return nil
}

// DecodeResponseHeader decodes and returns (if successful) a ResponseHeader from an input stream.
func DecodeResponseHeader(reader io.Reader, request *protocol.RequestHeader, responseAddons *Addons) error {

	buffer := buf.StackNew()
	defer buffer.Release()

	if _, err := buffer.ReadFullFrom(reader, 1); err != nil {
		return newError("failed to read response version").Base(err)
	}

	if buffer.Byte(0) != request.Version {
		return newError("unexpected response version. Expecting ", int(request.Version), " but actually ", int(buffer.Byte(0)))
	}

	responseAddons, err := DecodeHeaderAddons(&buffer, reader)
	if err != nil {
		return newError("failed to decode response header addons").Base(err)
	}

	return nil
}
