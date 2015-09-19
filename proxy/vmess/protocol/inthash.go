package protocol

import (
	"crypto/md5"
)

func Int64Hash(value int64) []byte {
	md5hash := md5.New()
	buffer := int64ToBytes(value)
	md5hash.Write(buffer)
	md5hash.Write(buffer)
	md5hash.Write(buffer)
	md5hash.Write(buffer)
	return md5hash.Sum(nil)
}
