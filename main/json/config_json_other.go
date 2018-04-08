// +build !windows

package json

import "syscall"

func getSysProcAttr() *syscall.SysProcAttr {
	return nil
}
