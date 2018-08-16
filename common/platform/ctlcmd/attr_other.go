// +build !windows

package ctlcmd

import "syscall"

func getSysProcAttr() *syscall.SysProcAttr {
	return nil
}
