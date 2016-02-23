package protocol

import (
	"github.com/v2ray/v2ray-core/common/serial"
)

type Timestamp int64

func (this Timestamp) Bytes() []byte {
	return serial.Int64Literal(this).Bytes()
}
