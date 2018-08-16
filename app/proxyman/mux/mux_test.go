package mux_test

import (
	"io"
	"testing"

	. "v2ray.com/core/app/proxyman/mux"
	"v2ray.com/core/common/buf"
	"v2ray.com/core/common/net"
	"v2ray.com/core/common/protocol"
	"v2ray.com/core/transport/pipe"
	. "v2ray.com/ext/assert"
)

func readAll(reader buf.Reader) (buf.MultiBuffer, error) {
	var mb buf.MultiBuffer
	for {
		b, err := reader.ReadMultiBuffer()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		mb.AppendMulti(b)
	}
	return mb, nil
}

func TestReaderWriter(t *testing.T) {
	assert := With(t)

	pReader, pWriter := pipe.New(pipe.WithSizeLimit(1024))

	dest := net.TCPDestination(net.DomainAddress("v2ray.com"), 80)
	writer := NewWriter(1, dest, pWriter, protocol.TransferTypeStream)

	dest2 := net.TCPDestination(net.LocalHostIP, 443)
	writer2 := NewWriter(2, dest2, pWriter, protocol.TransferTypeStream)

	dest3 := net.TCPDestination(net.LocalHostIPv6, 18374)
	writer3 := NewWriter(3, dest3, pWriter, protocol.TransferTypeStream)

	writePayload := func(writer *Writer, payload ...byte) error {
		b := buf.New()
		b.Write(payload)
		return writer.WriteMultiBuffer(buf.NewMultiBufferValue(b))
	}

	assert(writePayload(writer, 'a', 'b', 'c', 'd'), IsNil)
	assert(writePayload(writer2), IsNil)

	assert(writePayload(writer, 'e', 'f', 'g', 'h'), IsNil)
	assert(writePayload(writer3, 'x'), IsNil)

	writer.Close()
	writer3.Close()

	assert(writePayload(writer2, 'y'), IsNil)
	writer2.Close()

	bytesReader := &buf.BufferedReader{Reader: pReader}

	meta, err := ReadMetadata(bytesReader)
	assert(err, IsNil)
	assert(meta.SessionID, Equals, uint16(1))
	assert(byte(meta.SessionStatus), Equals, byte(SessionStatusNew))
	assert(meta.Target, Equals, dest)
	assert(byte(meta.Option), Equals, byte(OptionData))

	data, err := readAll(NewStreamReader(bytesReader))
	assert(err, IsNil)
	assert(len(data), Equals, 1)
	assert(data[0].String(), Equals, "abcd")

	meta, err = ReadMetadata(bytesReader)
	assert(err, IsNil)
	assert(byte(meta.SessionStatus), Equals, byte(SessionStatusNew))
	assert(meta.SessionID, Equals, uint16(2))
	assert(byte(meta.Option), Equals, byte(0))
	assert(meta.Target, Equals, dest2)

	meta, err = ReadMetadata(bytesReader)
	assert(err, IsNil)
	assert(byte(meta.SessionStatus), Equals, byte(SessionStatusKeep))
	assert(meta.SessionID, Equals, uint16(1))
	assert(byte(meta.Option), Equals, byte(1))

	data, err = readAll(NewStreamReader(bytesReader))
	assert(err, IsNil)
	assert(len(data), Equals, 1)
	assert(data[0].String(), Equals, "efgh")

	meta, err = ReadMetadata(bytesReader)
	assert(err, IsNil)
	assert(byte(meta.SessionStatus), Equals, byte(SessionStatusNew))
	assert(meta.SessionID, Equals, uint16(3))
	assert(byte(meta.Option), Equals, byte(1))
	assert(meta.Target, Equals, dest3)

	data, err = readAll(NewStreamReader(bytesReader))
	assert(err, IsNil)
	assert(len(data), Equals, 1)
	assert(data[0].String(), Equals, "x")

	meta, err = ReadMetadata(bytesReader)
	assert(err, IsNil)
	assert(byte(meta.SessionStatus), Equals, byte(SessionStatusEnd))
	assert(meta.SessionID, Equals, uint16(1))
	assert(byte(meta.Option), Equals, byte(0))

	meta, err = ReadMetadata(bytesReader)
	assert(err, IsNil)
	assert(byte(meta.SessionStatus), Equals, byte(SessionStatusEnd))
	assert(meta.SessionID, Equals, uint16(3))
	assert(byte(meta.Option), Equals, byte(0))

	meta, err = ReadMetadata(bytesReader)
	assert(err, IsNil)
	assert(byte(meta.SessionStatus), Equals, byte(SessionStatusKeep))
	assert(meta.SessionID, Equals, uint16(2))
	assert(byte(meta.Option), Equals, byte(1))

	data, err = readAll(NewStreamReader(bytesReader))
	assert(err, IsNil)
	assert(len(data), Equals, 1)
	assert(data[0].String(), Equals, "y")

	meta, err = ReadMetadata(bytesReader)
	assert(err, IsNil)
	assert(byte(meta.SessionStatus), Equals, byte(SessionStatusEnd))
	assert(meta.SessionID, Equals, uint16(2))
	assert(byte(meta.Option), Equals, byte(0))

	pWriter.Close()

	meta, err = ReadMetadata(bytesReader)
	assert(err, IsNotNil)
	assert(meta, IsNil)
}
