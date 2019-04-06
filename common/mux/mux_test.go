package mux_test

import (
	"io"
	"testing"

	"github.com/google/go-cmp/cmp"

	"v2ray.com/core/common"
	"v2ray.com/core/common/buf"
	. "v2ray.com/core/common/mux"
	"v2ray.com/core/common/net"
	"v2ray.com/core/common/protocol"
	"v2ray.com/core/transport/pipe"
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
		mb = append(mb, b...)
	}
	return mb, nil
}

func TestReaderWriter(t *testing.T) {
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
		return writer.WriteMultiBuffer(buf.MultiBuffer{b})
	}

	common.Must(writePayload(writer, 'a', 'b', 'c', 'd'))
	common.Must(writePayload(writer2))

	common.Must(writePayload(writer, 'e', 'f', 'g', 'h'))
	common.Must(writePayload(writer3, 'x'))

	writer.Close()
	writer3.Close()

	common.Must(writePayload(writer2, 'y'))
	writer2.Close()

	bytesReader := &buf.BufferedReader{Reader: pReader}

	{
		var meta FrameMetadata
		common.Must(meta.Unmarshal(bytesReader))
		if r := cmp.Diff(meta, FrameMetadata{
			SessionID:     1,
			SessionStatus: SessionStatusNew,
			Target:        dest,
			Option:        OptionData,
		}); r != "" {
			t.Error("metadata: ", r)
		}

		data, err := readAll(NewStreamReader(bytesReader))
		common.Must(err)
		if s := data.String(); s != "abcd" {
			t.Error("data: ", s)
		}
	}

	{
		var meta FrameMetadata
		common.Must(meta.Unmarshal(bytesReader))
		if r := cmp.Diff(meta, FrameMetadata{
			SessionStatus: SessionStatusNew,
			SessionID:     2,
			Option:        0,
			Target:        dest2,
		}); r != "" {
			t.Error("meta: ", r)
		}
	}

	{
		var meta FrameMetadata
		common.Must(meta.Unmarshal(bytesReader))
		if r := cmp.Diff(meta, FrameMetadata{
			SessionID:     1,
			SessionStatus: SessionStatusKeep,
			Option:        1,
		}); r != "" {
			t.Error("meta: ", r)
		}

		data, err := readAll(NewStreamReader(bytesReader))
		common.Must(err)
		if s := data.String(); s != "efgh" {
			t.Error("data: ", s)
		}
	}

	{
		var meta FrameMetadata
		common.Must(meta.Unmarshal(bytesReader))
		if r := cmp.Diff(meta, FrameMetadata{
			SessionID:     3,
			SessionStatus: SessionStatusNew,
			Option:        1,
			Target:        dest3,
		}); r != "" {
			t.Error("meta: ", r)
		}

		data, err := readAll(NewStreamReader(bytesReader))
		common.Must(err)
		if s := data.String(); s != "x" {
			t.Error("data: ", s)
		}
	}

	{
		var meta FrameMetadata
		common.Must(meta.Unmarshal(bytesReader))
		if r := cmp.Diff(meta, FrameMetadata{
			SessionID:     1,
			SessionStatus: SessionStatusEnd,
			Option:        0,
		}); r != "" {
			t.Error("meta: ", r)
		}
	}

	{
		var meta FrameMetadata
		common.Must(meta.Unmarshal(bytesReader))
		if r := cmp.Diff(meta, FrameMetadata{
			SessionID:     3,
			SessionStatus: SessionStatusEnd,
			Option:        0,
		}); r != "" {
			t.Error("meta: ", r)
		}
	}

	{
		var meta FrameMetadata
		common.Must(meta.Unmarshal(bytesReader))
		if r := cmp.Diff(meta, FrameMetadata{
			SessionID:     2,
			SessionStatus: SessionStatusKeep,
			Option:        1,
		}); r != "" {
			t.Error("meta: ", r)
		}

		data, err := readAll(NewStreamReader(bytesReader))
		common.Must(err)
		if s := data.String(); s != "y" {
			t.Error("data: ", s)
		}
	}

	{
		var meta FrameMetadata
		common.Must(meta.Unmarshal(bytesReader))
		if r := cmp.Diff(meta, FrameMetadata{
			SessionID:     2,
			SessionStatus: SessionStatusEnd,
			Option:        0,
		}); r != "" {
			t.Error("meta: ", r)
		}
	}

	pWriter.Close()

	{
		var meta FrameMetadata
		err := meta.Unmarshal(bytesReader)
		if err == nil {
			t.Error("nil error")
		}
	}
}
