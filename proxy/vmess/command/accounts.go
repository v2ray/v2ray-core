package command

import (
	"io"
	"time"

	"github.com/v2ray/v2ray-core/common/serial"
	"github.com/v2ray/v2ray-core/common/uuid"
	"github.com/v2ray/v2ray-core/transport"
)

func init() {
	RegisterResponseCommand(1, func() Command { return new(SwitchAccount) })
}

// Size: 16 + 8 = 24
type SwitchAccount struct {
	ID         *uuid.UUID
	ValidUntil time.Time
}

func (this *SwitchAccount) Marshal(writer io.Writer) (int, error) {
	idBytes := this.ID.Bytes()
	timestamp := this.ValidUntil.Unix()
	timeBytes := serial.Int64Literal(timestamp).Bytes()

	writer.Write(idBytes)
	writer.Write(timeBytes)

	return 24, nil
}

func (this *SwitchAccount) Unmarshal(data []byte) error {
	if len(data) != 24 {
		return transport.CorruptedPacket
	}
	this.ID, _ = uuid.ParseBytes(data[0:16])
	this.ValidUntil = time.Unix(serial.BytesLiteral(data[16:24]).Int64Value(), 0)
	return nil
}
