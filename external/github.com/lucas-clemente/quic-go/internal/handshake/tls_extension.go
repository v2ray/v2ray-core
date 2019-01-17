package handshake

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"

	"v2ray.com/core/external/github.com/lucas-clemente/quic-go/internal/protocol"
	"v2ray.com/core/external/github.com/lucas-clemente/quic-go/internal/utils"
)

const quicTLSExtensionType = 0xffa5

type clientHelloTransportParameters struct {
	InitialVersion protocol.VersionNumber
	Parameters     TransportParameters
}

func (p *clientHelloTransportParameters) Marshal() []byte {
	const lenOffset = 4
	b := &bytes.Buffer{}
	utils.BigEndian.WriteUint32(b, uint32(p.InitialVersion))
	b.Write([]byte{0, 0}) // length. Will be replaced later
	p.Parameters.marshal(b)
	data := b.Bytes()
	binary.BigEndian.PutUint16(data[lenOffset:lenOffset+2], uint16(len(data)-lenOffset-2))
	return data
}

func (p *clientHelloTransportParameters) Unmarshal(data []byte) error {
	if len(data) < 6 {
		return errors.New("transport parameter data too short")
	}
	p.InitialVersion = protocol.VersionNumber(binary.BigEndian.Uint32(data[:4]))
	paramsLen := int(binary.BigEndian.Uint16(data[4:6]))
	data = data[6:]
	if len(data) != paramsLen {
		return fmt.Errorf("expected transport parameters to be %d bytes long, have %d", paramsLen, len(data))
	}
	return p.Parameters.unmarshal(data, protocol.PerspectiveClient)
}

type encryptedExtensionsTransportParameters struct {
	NegotiatedVersion protocol.VersionNumber
	SupportedVersions []protocol.VersionNumber
	Parameters        TransportParameters
}

func (p *encryptedExtensionsTransportParameters) Marshal() []byte {
	b := &bytes.Buffer{}
	utils.BigEndian.WriteUint32(b, uint32(p.NegotiatedVersion))
	b.WriteByte(uint8(4 * len(p.SupportedVersions)))
	for _, v := range p.SupportedVersions {
		utils.BigEndian.WriteUint32(b, uint32(v))
	}
	lenOffset := b.Len()
	b.Write([]byte{0, 0}) // length. Will be replaced later
	p.Parameters.marshal(b)
	data := b.Bytes()
	binary.BigEndian.PutUint16(data[lenOffset:lenOffset+2], uint16(len(data)-lenOffset-2))
	return data
}

func (p *encryptedExtensionsTransportParameters) Unmarshal(data []byte) error {
	if len(data) < 5 {
		return errors.New("transport parameter data too short")
	}
	p.NegotiatedVersion = protocol.VersionNumber(binary.BigEndian.Uint32(data[:4]))
	numVersions := int(data[4])
	if numVersions%4 != 0 {
		return fmt.Errorf("invalid length for version list: %d", numVersions)
	}
	numVersions /= 4
	data = data[5:]
	if len(data) < 4*numVersions+2 /*length field for the parameter list */ {
		return errors.New("transport parameter data too short")
	}
	p.SupportedVersions = make([]protocol.VersionNumber, numVersions)
	for i := 0; i < numVersions; i++ {
		p.SupportedVersions[i] = protocol.VersionNumber(binary.BigEndian.Uint32(data[:4]))
		data = data[4:]
	}
	paramsLen := int(binary.BigEndian.Uint16(data[:2]))
	data = data[2:]
	if len(data) != paramsLen {
		return fmt.Errorf("expected transport parameters to be %d bytes long, have %d", paramsLen, len(data))
	}
	return p.Parameters.unmarshal(data, protocol.PerspectiveServer)
}
