package vmess

import (
	"crypto/md5"
	"errors"

	"github.com/v2ray/v2ray-core/common/uuid"
)

const (
	IDBytesLen = 16
)

var (
	InvalidID = errors.New("Invalid ID.")
)

// The ID of en entity, in the form of an UUID.
type ID struct {
	uuid   *uuid.UUID
	cmdKey [IDBytesLen]byte
}

func (this *ID) Equals(another *ID) bool {
	return this.uuid.Equals(another.uuid)
}

func (this *ID) Bytes() []byte {
	return this.uuid.Bytes()
}

func (this *ID) String() string {
	return this.uuid.String()
}

func (this *ID) UUID() *uuid.UUID {
	return this.uuid
}

func (v ID) CmdKey() []byte {
	return v.cmdKey[:]
}

func NewID(uuid *uuid.UUID) *ID {
	id := &ID{uuid: uuid}
	md5hash := md5.New()
	md5hash.Write(uuid.Bytes())
	md5hash.Write([]byte("c48619fe-8f02-49e0-b9e9-edf763e17e21"))
	md5.Sum(id.cmdKey[:0])
	return id
}
