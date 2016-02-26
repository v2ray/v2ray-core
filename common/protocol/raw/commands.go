package raw

import (
	"errors"
	"io"

	v2net "github.com/v2ray/v2ray-core/common/net"
	"github.com/v2ray/v2ray-core/common/protocol"
	"github.com/v2ray/v2ray-core/common/serial"
	"github.com/v2ray/v2ray-core/common/uuid"
	"github.com/v2ray/v2ray-core/transport"
)

var (
	ErrorCommandTypeMismatch = errors.New("Command type mismatch.")
	ErrorUnknownCommand      = errors.New("Unknown command.")
)

func MarshalCommand(command interface{}, writer io.Writer) error {
	var factory CommandFactory
	switch command.(type) {
	case *protocol.CommandSwitchAccount:
		factory = new(CommandSwitchAccountFactory)
	default:
		return ErrorUnknownCommand
	}
	return factory.Marshal(command, writer)
}

func UnmarshalCommand(cmdId byte, data []byte) (protocol.ResponseCommand, error) {
	var factory CommandFactory
	switch cmdId {
	case 1:
		factory = new(CommandSwitchAccountFactory)
	default:
		return nil, ErrorUnknownCommand
	}
	return factory.Unmarshal(data)
}

type CommandFactory interface {
	Marshal(command interface{}, writer io.Writer) error
	Unmarshal(data []byte) (interface{}, error)
}

type CommandSwitchAccountFactory struct {
}

func (this *CommandSwitchAccountFactory) Marshal(command interface{}, writer io.Writer) error {
	cmd, ok := command.(*protocol.CommandSwitchAccount)
	if !ok {
		return ErrorCommandTypeMismatch
	}

	hostStr := ""
	if cmd.Host != nil {
		hostStr = cmd.Host.String()
	}
	writer.Write([]byte{byte(len(hostStr))})

	if len(hostStr) > 0 {
		writer.Write([]byte(hostStr))
	}

	writer.Write(cmd.Port.Bytes())

	idBytes := cmd.ID.Bytes()
	writer.Write(idBytes)

	writer.Write(cmd.AlterIds.Bytes())
	writer.Write([]byte{byte(cmd.Level)})

	writer.Write([]byte{cmd.ValidMin})
	return nil
}

func (this *CommandSwitchAccountFactory) Unmarshal(data []byte) (interface{}, error) {
	cmd := new(protocol.CommandSwitchAccount)
	if len(data) == 0 {
		return nil, transport.ErrorCorruptedPacket
	}
	lenHost := int(data[0])
	if len(data) < lenHost+1 {
		return nil, transport.ErrorCorruptedPacket
	}
	if lenHost > 0 {
		cmd.Host = v2net.ParseAddress(string(data[1 : 1+lenHost]))
	}
	portStart := 1 + lenHost
	if len(data) < portStart+2 {
		return nil, transport.ErrorCorruptedPacket
	}
	cmd.Port = v2net.PortFromBytes(data[portStart : portStart+2])
	idStart := portStart + 2
	if len(data) < idStart+16 {
		return nil, transport.ErrorCorruptedPacket
	}
	cmd.ID, _ = uuid.ParseBytes(data[idStart : idStart+16])
	alterIdStart := idStart + 16
	if len(data) < alterIdStart+2 {
		return nil, transport.ErrorCorruptedPacket
	}
	cmd.AlterIds = serial.BytesLiteral(data[alterIdStart : alterIdStart+2]).Uint16()
	levelStart := alterIdStart + 2
	if len(data) < levelStart+1 {
		return nil, transport.ErrorCorruptedPacket
	}
	cmd.Level = protocol.UserLevel(data[levelStart])
	timeStart := levelStart + 1
	if len(data) < timeStart {
		return nil, transport.ErrorCorruptedPacket
	}
	cmd.ValidMin = data[timeStart]
	return cmd, nil
}
