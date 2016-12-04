package internal

import (
	"net"
	"reflect"
	"v2ray.com/core/common/errors"
)

var (
	ErrInvalidConn = errors.New("Invalid Connection.")
)

func GetSysFd(conn net.Conn) (int, error) {
	cv := reflect.ValueOf(conn)
	switch ce := cv.Elem(); ce.Kind() {
	case reflect.Struct:
		netfd := ce.FieldByName("conn").FieldByName("fd")
		switch fe := netfd.Elem(); fe.Kind() {
		case reflect.Struct:
			fd := fe.FieldByName("sysfd")
			return int(fd.Int()), nil
		}
	}
	return 0, ErrInvalidConn
}
