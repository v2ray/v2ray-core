package serial

import (
	"encoding/hex"
	"strings"
)

func ByteToHexString(value byte) string {
	return hex.EncodeToString([]byte{value})
}

// BytesToUint16 deserializes a byte array to an uint16 in big endian order. The byte array must have at least 2 elements.
func BytesToUint16(value []byte) uint16 {
	return uint16(value[0])<<8 | uint16(value[1])
}

// BytesToUint32 deserializes a byte array to an uint32 in big endian order. The byte array must have at least 4 elements.
func BytesToUint32(value []byte) uint32 {
	return uint32(value[0])<<24 |
		uint32(value[1])<<16 |
		uint32(value[2])<<8 |
		uint32(value[3])
}

// BytesToInt64 deserializes a byte array to an int64 in big endian order. The byte array must have at least 8 elements.
func BytesToInt64(value []byte) int64 {
	return int64(value[0])<<56 |
		int64(value[1])<<48 |
		int64(value[2])<<40 |
		int64(value[3])<<32 |
		int64(value[4])<<24 |
		int64(value[5])<<16 |
		int64(value[6])<<8 |
		int64(value[7])
}

func BytesToHexString(value []byte) string {
	strs := make([]string, len(value))
	for i, b := range value {
		strs[i] = hex.EncodeToString([]byte{b})
	}
	return "[" + strings.Join(strs, ",") + "]"
}
