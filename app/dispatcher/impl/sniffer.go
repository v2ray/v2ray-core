package impl

import (
	"bytes"
	"strings"

	"v2ray.com/core/common/serial"
)

var (
	ErrMoreData    = newError("need more data")
	ErrInvalidData = newError("invalid data")
)

func ContainsValidHTTPMethod(b []byte) bool {
	if len(b) == 0 {
		return false
	}

	parts := bytes.Split(b, []byte{' '})
	part0Trimed := strings.ToLower(string(bytes.Trim(parts[0], " ")))
	return part0Trimed == "get" || part0Trimed == "post" ||
		part0Trimed == "head" || part0Trimed == "put" ||
		part0Trimed == "delete" || part0Trimed == "options" || part0Trimed == "connect"
}

func SniffHTTP(b []byte) (string, error) {
	if len(b) == 0 {
		return "", ErrMoreData
	}
	headers := bytes.Split(b, []byte{'\n'})
	if !ContainsValidHTTPMethod(headers[0]) {
		return "", ErrInvalidData
	}
	for i := 1; i < len(headers); i++ {
		header := headers[i]
		if len(header) == 0 {
			return "", ErrInvalidData
		}
		parts := bytes.SplitN(header, []byte{':'}, 2)
		if len(parts) != 2 {
			return "", ErrInvalidData
		}
		key := strings.ToLower(string(parts[0]))
		value := strings.ToLower(string(bytes.Trim(parts[1], " ")))
		if key == "host" {
			domain := strings.Split(value, ":")
			return domain[0], nil
		}
	}
	return "", ErrMoreData
}

func IsValidTLSVersion(major, minor byte) bool {
	return major == 3
}

// ReadClientHello returns server name (if any) from TLS client hello message.
// https://github.com/golang/go/blob/master/src/crypto/tls/handshake_messages.go#L300
func ReadClientHello(data []byte) (string, error) {
	if len(data) < 42 {
		return "", ErrMoreData
	}
	sessionIDLen := int(data[38])
	if sessionIDLen > 32 || len(data) < 39+sessionIDLen {
		return "", ErrInvalidData
	}
	data = data[39+sessionIDLen:]
	if len(data) < 2 {
		return "", ErrMoreData
	}
	// cipherSuiteLen is the number of bytes of cipher suite numbers. Since
	// they are uint16s, the number must be even.
	cipherSuiteLen := int(data[0])<<8 | int(data[1])
	if cipherSuiteLen%2 == 1 || len(data) < 2+cipherSuiteLen {
		return "", ErrInvalidData
	}
	data = data[2+cipherSuiteLen:]
	if len(data) < 1 {
		return "", ErrMoreData
	}
	compressionMethodsLen := int(data[0])
	if len(data) < 1+compressionMethodsLen {
		return "", ErrMoreData
	}
	data = data[1+compressionMethodsLen:]

	if len(data) == 0 {
		return "", ErrInvalidData
	}
	if len(data) < 2 {
		return "", ErrInvalidData
	}

	extensionsLength := int(data[0])<<8 | int(data[1])
	data = data[2:]
	if extensionsLength != len(data) {
		return "", ErrInvalidData
	}

	for len(data) != 0 {
		if len(data) < 4 {
			return "", ErrInvalidData
		}
		extension := uint16(data[0])<<8 | uint16(data[1])
		length := int(data[2])<<8 | int(data[3])
		data = data[4:]
		if len(data) < length {
			return "", ErrInvalidData
		}

		switch extension {
		case 0x00: /* extensionServerName */
			d := data[:length]
			if len(d) < 2 {
				return "", ErrInvalidData
			}
			namesLen := int(d[0])<<8 | int(d[1])
			d = d[2:]
			if len(d) != namesLen {
				return "", ErrInvalidData
			}
			for len(d) > 0 {
				if len(d) < 3 {
					return "", ErrInvalidData
				}
				nameType := d[0]
				nameLen := int(d[1])<<8 | int(d[2])
				d = d[3:]
				if len(d) < nameLen {
					return "", ErrInvalidData
				}
				if nameType == 0 {
					serverName := string(d[:nameLen])
					// An SNI value may not include a
					// trailing dot. See
					// https://tools.ietf.org/html/rfc6066#section-3.
					if strings.HasSuffix(serverName, ".") {
						return "", ErrInvalidData
					}
					return serverName, nil
				}
				d = d[nameLen:]
			}
		}
		data = data[length:]
	}

	return "", ErrInvalidData
}

func SniffTLS(b []byte) (string, error) {
	if len(b) < 5 {
		return "", ErrMoreData
	}

	if b[0] != 0x16 /* TLS Handshake */ {
		return "", ErrInvalidData
	}
	if !IsValidTLSVersion(b[1], b[2]) {
		return "", ErrInvalidData
	}
	headerLen := int(serial.BytesToUint16(b[3:5]))
	if 5+headerLen > len(b) {
		return "", ErrMoreData
	}
	return ReadClientHello(b[5 : 5+headerLen])
}
