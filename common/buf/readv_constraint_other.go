// +build !windows

package buf

import "syscall"

func checkReadVConstraint(conn syscall.RawConn) (bool, error) {
	return true, nil
}
