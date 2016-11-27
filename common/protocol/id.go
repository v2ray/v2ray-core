package protocol

import (
	"crypto/hmac"
	"crypto/md5"
	"errors"
	"hash"

	"v2ray.com/core/common/uuid"
)

const (
	IDBytesLen = 16
)

var (
	InvalidID = errors.New("Invalid ID.")
)

type IDHash func(key []byte) hash.Hash

func DefaultIDHash(key []byte) hash.Hash {
	return hmac.New(md5.New, key)
}

// The ID of en entity, in the form of an UUID.
type ID struct {
	uuid   *uuid.UUID
	cmdKey [IDBytesLen]byte
}

func (v *ID) Equals(another *ID) bool {
	return v.uuid.Equals(another.uuid)
}

func (v *ID) Bytes() []byte {
	return v.uuid.Bytes()
}

func (v *ID) String() string {
	return v.uuid.String()
}

func (v *ID) UUID() *uuid.UUID {
	return v.uuid
}

func (v ID) CmdKey() []byte {
	return v.cmdKey[:]
}

func NewID(uuid *uuid.UUID) *ID {
	id := &ID{uuid: uuid}
	md5hash := md5.New()
	md5hash.Write(uuid.Bytes())
	md5hash.Write([]byte("c48619fe-8f02-49e0-b9e9-edf763e17e21"))
	md5hash.Sum(id.cmdKey[:0])
	return id
}

func NewAlterIDs(primary *ID, alterIDCount uint16) []*ID {
	alterIDs := make([]*ID, alterIDCount)
	prevID := primary.UUID()
	for idx := range alterIDs {
		newid := prevID.Next()
		// TODO: check duplicates
		alterIDs[idx] = NewID(newid)
		prevID = newid
	}
	return alterIDs
}
