package encoding

import (
	"io"

	"v2ray.com/core/common"
	"v2ray.com/core/common/buf"
	"v2ray.com/core/common/net"
	"v2ray.com/core/common/protocol"
	"v2ray.com/core/common/serial"
	"v2ray.com/core/common/uuid"
)

var (
	ErrCommandTypeMismatch = newError("Command type mismatch.")
	ErrUnknownCommand      = newError("Unknown command.")
	ErrCommandTooLarge     = newError("Command too large.")
)

func MarshalCommand(command interface{}, writer io.Writer) error {
	if command == nil {
		return ErrUnknownCommand
	}

	var cmdID byte
	var factory CommandFactory
	switch command.(type) {
	case *protocol.CommandSwitchAccount:
		factory = new(CommandSwitchAccountFactory)
		cmdID = 1
	default:
		return ErrUnknownCommand
	}

	buffer := buf.NewLocal(512)
	defer buffer.Release()

	err := factory.Marshal(command, buffer)
	if err != nil {
		return err
	}

	auth := Authenticate(buffer.Bytes())
	len := buffer.Len() + 4
	if len > 255 {
		return ErrCommandTooLarge
	}

	common.Must2(writer.Write([]byte{cmdID, byte(len), byte(auth >> 24), byte(auth >> 16), byte(auth >> 8), byte(auth)}))
	common.Must2(writer.Write(buffer.Bytes()))
	return nil
}

func UnmarshalCommand(cmdID byte, data []byte) (protocol.ResponseCommand, error) {
	if len(data) <= 4 {
		return nil, newError("insufficient length")
	}
	expectedAuth := Authenticate(data[4:])
	actualAuth := serial.BytesToUint32(data[:4])
	if expectedAuth != actualAuth {
		return nil, newError("invalid auth")
	}

	var factory CommandFactory
	switch cmdID {
	case 1:
		factory = new(CommandSwitchAccountFactory)
	default:
		return nil, ErrUnknownCommand
	}
	return factory.Unmarshal(data[4:])
}

type CommandFactory interface {
	Marshal(command interface{}, writer io.Writer) error
	Unmarshal(data []byte) (interface{}, error)
}

type CommandSwitchAccountFactory struct {
}

func (f *CommandSwitchAccountFactory) Marshal(command interface{}, writer io.Writer) error {
	cmd, ok := command.(*protocol.CommandSwitchAccount)
	if !ok {
		return ErrCommandTypeMismatch
	}

	hostStr := ""
	if cmd.Host != nil {
		hostStr = cmd.Host.String()
	}
	common.Must2(writer.Write([]byte{byte(len(hostStr))}))

	if len(hostStr) > 0 {
		common.Must2(writer.Write([]byte(hostStr)))
	}

	common.Must2(writer.Write(cmd.Port.Bytes(nil)))

	idBytes := cmd.ID.Bytes()
	common.Must2(writer.Write(idBytes))

	common.Must2(writer.Write(serial.Uint16ToBytes(cmd.AlterIds, nil)))
	common.Must2(writer.Write([]byte{byte(cmd.Level)}))

	common.Must2(writer.Write([]byte{cmd.ValidMin}))
	return nil
}

func (f *CommandSwitchAccountFactory) Unmarshal(data []byte) (interface{}, error) {
	cmd := new(protocol.CommandSwitchAccount)
	if len(data) == 0 {
		return nil, newError("insufficient length.")
	}
	lenHost := int(data[0])
	if len(data) < lenHost+1 {
		return nil, newError("insufficient length.")
	}
	if lenHost > 0 {
		cmd.Host = net.ParseAddress(string(data[1 : 1+lenHost]))
	}
	portStart := 1 + lenHost
	if len(data) < portStart+2 {
		return nil, newError("insufficient length.")
	}
	cmd.Port = net.PortFromBytes(data[portStart : portStart+2])
	idStart := portStart + 2
	if len(data) < idStart+16 {
		return nil, newError("insufficient length.")
	}
	cmd.ID, _ = uuid.ParseBytes(data[idStart : idStart+16])
	alterIDStart := idStart + 16
	if len(data) < alterIDStart+2 {
		return nil, newError("insufficient length.")
	}
	cmd.AlterIds = serial.BytesToUint16(data[alterIDStart : alterIDStart+2])
	levelStart := alterIDStart + 2
	if len(data) < levelStart+1 {
		return nil, newError("insufficient length.")
	}
	cmd.Level = uint32(data[levelStart])
	timeStart := levelStart + 1
	if len(data) < timeStart {
		return nil, newError("insufficient length.")
	}
	cmd.ValidMin = data[timeStart]
	return cmd, nil
}
