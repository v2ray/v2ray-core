package tls

import (
	"errors"
	"strings"

	"v2ray.com/core/common"
	"v2ray.com/core/common/serial"
)

type SniffHeader struct {
	domain string
}

func (h *SniffHeader) Protocol() string {
	return "tls"
}

func (h *SniffHeader) Domain() string {
	return h.domain
}

var errNotTLS = errors.New("not TLS header")
var errNotClientHello = errors.New("not client hello")

func IsValidTLSVersion(major, minor byte) bool {
	return major == 3
}

// ReadClientHello returns server name (if any) from TLS client hello message.
// https://github.com/golang/go/blob/master/src/crypto/tls/handshake_messages.go#L300
func ReadClientHello(data []byte, h *SniffHeader) error {
	if len(data) < 42 {
		return common.ErrNoClue
	}
	sessionIDLen := int(data[38])
	if sessionIDLen > 32 || len(data) < 39+sessionIDLen {
		return common.ErrNoClue
	}
	data = data[39+sessionIDLen:]
	if len(data) < 2 {
		return common.ErrNoClue
	}
	// cipherSuiteLen is the number of bytes of cipher suite numbers. Since
	// they are uint16s, the number must be even.
	cipherSuiteLen := int(data[0])<<8 | int(data[1])
	if cipherSuiteLen%2 == 1 || len(data) < 2+cipherSuiteLen {
		return errNotClientHello
	}
	data = data[2+cipherSuiteLen:]
	if len(data) < 1 {
		return common.ErrNoClue
	}
	compressionMethodsLen := int(data[0])
	if len(data) < 1+compressionMethodsLen {
		return common.ErrNoClue
	}
	data = data[1+compressionMethodsLen:]

	if len(data) == 0 {
		return errNotClientHello
	}
	if len(data) < 2 {
		return errNotClientHello
	}

	extensionsLength := int(data[0])<<8 | int(data[1])
	data = data[2:]
	if extensionsLength != len(data) {
		return errNotClientHello
	}

	for len(data) != 0 {
		if len(data) < 4 {
			return errNotClientHello
		}
		extension := uint16(data[0])<<8 | uint16(data[1])
		length := int(data[2])<<8 | int(data[3])
		data = data[4:]
		if len(data) < length {
			return errNotClientHello
		}

		switch extension {
		case 0x00: /* extensionServerName */
			d := data[:length]
			if len(d) < 2 {
				return errNotClientHello
			}
			namesLen := int(d[0])<<8 | int(d[1])
			d = d[2:]
			if len(d) != namesLen {
				return errNotClientHello
			}
			for len(d) > 0 {
				if len(d) < 3 {
					return errNotClientHello
				}
				nameType := d[0]
				nameLen := int(d[1])<<8 | int(d[2])
				d = d[3:]
				if len(d) < nameLen {
					return errNotClientHello
				}
				if nameType == 0 {
					serverName := string(d[:nameLen])
					// An SNI value may not include a
					// trailing dot. See
					// https://tools.ietf.org/html/rfc6066#section-3.
					if strings.HasSuffix(serverName, ".") {
						return errNotClientHello
					}
					h.domain = serverName
					return nil
				}
				d = d[nameLen:]
			}
		}
		data = data[length:]
	}

	return errNotTLS
}

func SniffTLS(b []byte) (*SniffHeader, error) {
	if len(b) < 5 {
		return nil, common.ErrNoClue
	}

	if b[0] != 0x16 /* TLS Handshake */ {
		return nil, errNotTLS
	}
	if !IsValidTLSVersion(b[1], b[2]) {
		return nil, errNotTLS
	}
	headerLen := int(serial.BytesToUint16(b[3:5]))
	if 5+headerLen > len(b) {
		return nil, common.ErrNoClue
	}

	h := &SniffHeader{}
	err := ReadClientHello(b[5:5+headerLen], h)
	if err == nil {
		return h, nil
	}
	return nil, err
}
