// +build !confonly

package encoding

import (
	"io"

	"github.com/golang/protobuf/proto"

	"v2ray.com/core/common/buf"
	"v2ray.com/core/common/protocol"
)

func EncodeHeaderAddons(buffer *buf.Buffer, addons *Addons) error {

	switch addons.Flow {
	default:

		if err := buffer.WriteByte(0); err != nil {
			return newError("failed to write addons protobuf length").Base(err)
		}

	}

	return nil

}

func DecodeHeaderAddons(buffer *buf.Buffer, reader io.Reader) (*Addons, error) {

	addons := new(Addons)

	buffer.Clear()
	if _, err := buffer.ReadFullFrom(reader, 1); err != nil {
		return nil, newError("failed to read addons protobuf length").Base(err)
	}

	if length := int32(buffer.Byte(0)); length != 0 {

		buffer.Clear()
		if _, err := buffer.ReadFullFrom(reader, length); err != nil {
			return nil, newError("failed to read addons protobuf value").Base(err)
		}

		if err := proto.Unmarshal(buffer.Bytes(), addons); err != nil {
			return nil, newError("failed to unmarshal addons protobuf value").Base(err)
		}

		// Verification.
		switch addons.Flow {
		default:

		}

	}

	return addons, nil

}

// EncodeBodyAddons returns a Writer that auto-encrypt content written by caller.
func EncodeBodyAddons(writer io.Writer, request *protocol.RequestHeader, addons *Addons) buf.Writer {

	switch addons.Flow {
	default:

		return buf.NewWriter(writer)

	}

}

// DecodeBodyAddons returns a Reader from which caller can fetch decrypted body.
func DecodeBodyAddons(reader io.Reader, request *protocol.RequestHeader, addons *Addons) buf.Reader {

	switch addons.Flow {
	default:

		return buf.NewReader(reader)

	}

}
