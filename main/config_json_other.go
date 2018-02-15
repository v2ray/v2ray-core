// +build !windows

package main

import "syscall"

func getSysProcAttr() *syscall.SysProcAttr {
	return nil
}
