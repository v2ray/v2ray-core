package wire

import (
	"bytes"
	"crypto/rand"

	"v2ray.com/core/external/github.com/lucas-clemente/quic-go/internal/protocol"
	"v2ray.com/core/external/github.com/lucas-clemente/quic-go/internal/utils"
)

// ComposeVersionNegotiation composes a Version Negotiation
func ComposeVersionNegotiation(destConnID, srcConnID protocol.ConnectionID, versions []protocol.VersionNumber) ([]byte, error) {
	greasedVersions := protocol.GetGreasedVersions(versions)
	expectedLen := 1 /* type byte */ + 4 /* version field */ + 1 /* connection ID length field */ + destConnID.Len() + srcConnID.Len() + len(greasedVersions)*4
	buf := bytes.NewBuffer(make([]byte, 0, expectedLen))
	r := make([]byte, 1)
	_, _ = rand.Read(r) // ignore the error here. It is not critical to have perfect random here.
	buf.WriteByte(r[0] | 0x80)
	utils.BigEndian.WriteUint32(buf, 0) // version 0
	connIDLen, err := encodeConnIDLen(destConnID, srcConnID)
	if err != nil {
		return nil, err
	}
	buf.WriteByte(connIDLen)
	buf.Write(destConnID)
	buf.Write(srcConnID)
	for _, v := range greasedVersions {
		utils.BigEndian.WriteUint32(buf, uint32(v))
	}
	return buf.Bytes(), nil
}
