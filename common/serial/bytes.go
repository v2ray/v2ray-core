package serial

import (
	"bytes"
	"encoding/hex"
	"strings"
)

func ByteToHexString(value byte) string {
	return hex.EncodeToString([]byte{value})
}

func BytesToUint16(value []byte) uint16 {
	return uint16(value[0])<<8 + uint16(value[1])
}

func BytesToHexString(value []byte) string {
	strs := make([]string, len(value))
	for i, b := range value {
		strs[i] = hex.EncodeToString([]byte{b})
	}
	return "[" + strings.Join(strs, ",") + "]"
}

type BytesT []byte

func (this BytesT) Value() []byte {
	return []byte(this)
}

func (this BytesT) Equals(another BytesT) bool {
	return bytes.Equal(this.Value(), another.Value())
}

func (this BytesT) Uint8Value() uint8 {
	return this.Value()[0]
}

func (this BytesT) Uint16() Uint16Literal {
	return Uint16Literal(this.Uint16Value())
}

func (this BytesT) Uint16Value() uint16 {
	value := this.Value()
	return uint16(value[0])<<8 + uint16(value[1])
}

func (this BytesT) IntValue() int {
	value := this.Value()
	return int(value[0])<<24 + int(value[1])<<16 + int(value[2])<<8 + int(value[3])
}

func (this BytesT) Uint32Value() uint32 {
	value := this.Value()
	return uint32(value[0])<<24 +
		uint32(value[1])<<16 +
		uint32(value[2])<<8 +
		uint32(value[3])
}

func (this BytesT) Int64Value() int64 {
	value := this.Value()
	return int64(value[0])<<56 +
		int64(value[1])<<48 +
		int64(value[2])<<40 +
		int64(value[3])<<32 +
		int64(value[4])<<24 +
		int64(value[5])<<16 +
		int64(value[6])<<8 +
		int64(value[7])
}

// String returns a string presentation of this ByteLiteral
func (this BytesT) String() string {
	return string(this.Value())
}

// All returns true if all bytes in the ByteLiteral are the same as given value.
func (this BytesT) All(v byte) bool {
	for _, b := range this {
		if b != v {
			return false
		}
	}
	return true
}
