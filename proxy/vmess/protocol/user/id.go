package user

import (
	"crypto/md5"
	"encoding/hex"
	"errors"

	"github.com/v2ray/v2ray-core/common/log"
)

const (
	IDBytesLen = 16
)

var (
	InvalidID = errors.New("Invalid ID.")
)

// The ID of en entity, in the form of an UUID.
type ID struct {
	String string
	Bytes  [IDBytesLen]byte
	cmdKey [IDBytesLen]byte
}

func NewID(id string) (ID, error) {
	idBytes, err := UUIDToID(id)
	if err != nil {
		log.Error("Failed to parse id %s", id)
		return ID{}, InvalidID
	}

	md5hash := md5.New()
	md5hash.Write(idBytes[:])
	md5hash.Write([]byte("c48619fe-8f02-49e0-b9e9-edf763e17e21"))
	cmdKey := md5.Sum(nil)

	return ID{
		String: id,
		Bytes:  idBytes,
		cmdKey: cmdKey,
	}, nil
}

func (v ID) CmdKey() []byte {
	return v.cmdKey[:]
}

var byteGroups = []int{8, 4, 4, 4, 12}

// TODO: leverage a full functional UUID library
func UUIDToID(uuid string) (v [IDBytesLen]byte, err error) {
	text := []byte(uuid)
	if len(text) < 32 {
		log.Error("uuid: invalid UUID string: %s", text)
		err = InvalidID
		return
	}

	b := v[:]

	for _, byteGroup := range byteGroups {
		if text[0] == '-' {
			text = text[1:]
		}

		_, err = hex.Decode(b[:byteGroup/2], text[:byteGroup])

		if err != nil {
			return
		}

		text = text[byteGroup:]
		b = b[byteGroup/2:]
	}

	return
}
