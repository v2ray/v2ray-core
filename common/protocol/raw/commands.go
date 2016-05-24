package raw

import (
	"errors"
	"io"

	"github.com/v2ray/v2ray-core/common/alloc"
	v2net "github.com/v2ray/v2ray-core/common/net"
	"github.com/v2ray/v2ray-core/common/protocol"
	"github.com/v2ray/v2ray-core/common/serial"
	"github.com/v2ray/v2ray-core/common/uuid"
	"github.com/v2ray/v2ray-core/transport"
)

var (
	ErrorCommandTypeMismatch = errors.New("Command type mismatch.")
	ErrorUnknownCommand      = errors.New("Unknown command.")
	ErrorCommandTooLarge     = errors.New("Command too large.")
)

func MarshalCommand(command interface{}, writer io.Writer) error {
	if command == nil {
		return ErrorUnknownCommand
	}

	var cmdId byte
	var factory CommandFactory
	switch command.(type) {
	case *protocol.CommandSwitchAccount:
		factory = new(CommandSwitchAccountFactory)
		cmdId = 1
	default:
		return ErrorUnknownCommand
	}

	buffer := alloc.NewSmallBuffer().Clear()
	err := factory.Marshal(command, buffer)
	if err != nil {
		return err
	}

	auth := Authenticate(buffer.Value)
	len := buffer.Len() + 4
	if len > 255 {
		return ErrorCommandTooLarge
	}

	writer.Write([]byte{cmdId, byte(len), byte(auth >> 24), byte(auth >> 16), byte(auth >> 8), byte(auth)})
	writer.Write(buffer.Value)
	return nil
}

func UnmarshalCommand(cmdId byte, data []byte) (protocol.ResponseCommand, error) {
	if len(data) <= 4 {
		return nil, transport.ErrorCorruptedPacket
	}
	expectedAuth := Authenticate(data[4:])
	actualAuth := serial.BytesT(data[:4]).Uint32Value()
	if expectedAuth != actualAuth {
		return nil, transport.ErrorCorruptedPacket
	}

	var factory CommandFactory
	switch cmdId {
	case 1:
		factory = new(CommandSwitchAccountFactory)
	default:
		return nil, ErrorUnknownCommand
	}
	return factory.Unmarshal(data[4:])
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

	writer.Write(serial.Uint16ToBytes(cmd.AlterIds))
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
	cmd.AlterIds = serial.BytesToUint16(data[alterIdStart : alterIdStart+2])
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
