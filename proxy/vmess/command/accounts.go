package command

import (
	"io"

	v2net "github.com/v2ray/v2ray-core/common/net"
	"github.com/v2ray/v2ray-core/common/serial"
	"github.com/v2ray/v2ray-core/common/uuid"
	"github.com/v2ray/v2ray-core/proxy/vmess"
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
// 1 byte: level
// 1 bytes: time
type SwitchAccount struct {
	Host     v2net.Address
	Port     v2net.Port
	ID       *uuid.UUID
	AlterIds serial.Uint16Literal
	Level    vmess.UserLevel
	ValidMin byte
}

func (this *SwitchAccount) Marshal(writer io.Writer) {
	hostStr := ""
	if this.Host != nil {
		hostStr = this.Host.String()
	}
	writer.Write([]byte{byte(len(hostStr))})

	if len(hostStr) > 0 {
		writer.Write([]byte(hostStr))
	}

	writer.Write(this.Port.Bytes())

	idBytes := this.ID.Bytes()
	writer.Write(idBytes)

	writer.Write(this.AlterIds.Bytes())
	writer.Write([]byte{byte(this.Level)})

	writer.Write([]byte{this.ValidMin})
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
	levelStart := alterIdStart + 2
	if len(data) < levelStart {
		return transport.CorruptedPacket
	}
	this.Level = vmess.UserLevel(data[levelStart])
	timeStart := levelStart + 1
	if len(data) < timeStart {
		return transport.CorruptedPacket
	}
	this.ValidMin = data[timeStart]
	return nil
}
