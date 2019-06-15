// +build windows
package buf

import (
	"syscall"
)

func checkReadVConstraint(conn syscall.RawConn) (bool, error) {
	var isSocketReady = false
	var reason error
	/*
		In Windows, WSARecv system call only support socket connection.

		It it required to check if the given fd is of a socket type

		Fix https://github.com/v2ray/v2ray-core/issues/1666

		Additional Information:
		https://docs.microsoft.com/en-us/windows/desktop/api/winsock2/nf-winsock2-wsarecv
		https://docs.microsoft.com/en-us/windows/desktop/api/winsock/nf-winsock-getsockopt
		https://docs.microsoft.com/en-us/windows/desktop/WinSock/sol-socket-socket-options

	*/
	err := conn.Control(func(fd uintptr) {
		var val [4]byte
		var le = int32(len(val))
		err := syscall.Getsockopt(syscall.Handle(fd), syscall.SOL_SOCKET, syscall.SO_RCVBUF, &val[0], &le)
		if err != nil {
			isSocketReady = false
		} else {
			isSocketReady = true
		}
		reason = err
	})

	return isSocketReady, err
}
