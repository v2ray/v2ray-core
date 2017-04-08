package internal

import (
	"net"
	"reflect"
)

var (
	errInvalidConn = newError("Invalid Connection.")
)

// GetSysFd returns the underlying fd of a connection.
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
	return 0, errInvalidConn
}
