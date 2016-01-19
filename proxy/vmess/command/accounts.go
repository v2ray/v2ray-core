package command

import (
	"io"
	"time"

	v2net "github.com/v2ray/v2ray-core/common/net"
	"github.com/v2ray/v2ray-core/common/serial"
	"github.com/v2ray/v2ray-core/common/uuid"
	"github.com/v2ray/v2ray-core/transport"
)

func init() {
	RegisterResponseCommand(1, func() Command { return new(SwitchAccount) })
}

// Structure
// 1 byte: host len N
// N bytes: host
// 2 bytes: port
// 16 bytes: uuid
// 2 bytes: alterid
// 8 bytes: time
type SwitchAccount struct {
	Host       v2net.Address
	Port       v2net.Port
	ID         *uuid.UUID
	AlterIds   serial.Uint16Literal
	ValidUntil time.Time
}

func (this *SwitchAccount) Marshal(writer io.Writer) (int, error) {
	outBytes := 0
	hostStr := ""
	if this.Host != nil {
		hostStr = this.Host.String()
	}
	writer.Write([]byte{byte(len(hostStr))})
	outBytes++

	if len(hostStr) > 0 {
		writer.Write([]byte(hostStr))
		outBytes += len(hostStr)
	}

	writer.Write(this.Port.Bytes())
	outBytes += 2

	idBytes := this.ID.Bytes()
	writer.Write(idBytes)
	outBytes += len(idBytes)

	writer.Write(this.AlterIds.Bytes())
	outBytes += 2

	timestamp := this.ValidUntil.Unix()
	timeBytes := serial.Int64Literal(timestamp).Bytes()

	writer.Write(timeBytes)
	outBytes += len(timeBytes)

	return outBytes, nil
}

func (this *SwitchAccount) Unmarshal(data []byte) error {
	lenHost := int(data[0])
	if len(data) < lenHost+1 {
		return transport.CorruptedPacket
	}
	this.Host = v2net.ParseAddress(string(data[1 : 1+lenHost]))
	portStart := 1 + lenHost
	if len(data) < portStart+2 {
		return transport.CorruptedPacket
	}
	this.Port = v2net.PortFromBytes(data[portStart : portStart+2])
	idStart := portStart + 2
	if len(data) < idStart+16 {
		return transport.CorruptedPacket
	}
	this.ID, _ = uuid.ParseBytes(data[idStart : idStart+16])
	alterIdStart := idStart + 16
	if len(data) < alterIdStart+2 {
		return transport.CorruptedPacket
	}
	this.AlterIds = serial.ParseUint16(data[alterIdStart : alterIdStart+2])
	timeStart := alterIdStart + 2
	if len(data) < timeStart+8 {
		return transport.CorruptedPacket
	}
	this.ValidUntil = time.Unix(serial.BytesLiteral(data[timeStart:timeStart+8]).Int64Value(), 0)
	return nil
}
