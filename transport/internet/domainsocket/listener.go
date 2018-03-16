package domainsocket

import (
	"context"
	"net"
	"os"
	"syscall"
)

type Listener struct {
	ln           net.Listener
	listenerChan <-chan net.Conn
	ctx          context.Context
	path         string
	lockfile     os.File
}

func ListenDS(ctx context.Context, path string) (*Listener, error) {

	vln := &Listener{path: path}
	return vln, nil
}

func (ls *Listener) Down() error {
	err := ls.ln.Close()
	if err != nil {
		newError(err).AtDebug().WriteToLog()
	}
	return err
}

//Setup systen level Listener
func (ls *Listener) LowerUP() error {

	if isUnixDomainSocketFileSystemBased(ls.path) && !___DEBUG_IGNORE_FLOCK {

	}

	addr := new(net.UnixAddr)
	addr.Name = ls.path
	addr.Net = "unix"
	li, err := net.ListenUnix("unix", addr)
	if err != nil {
		return err
	}

}

func isUnixDomainSocketFileSystemBased(path string) bool {
	//No Branching
	return path[0] != 0
}

func AcquireLock(lockfilepath string) (*os.File, error) {
	f, err := os.Create(lockfilepath)
	if err != nil {
		newError(err).AtDebug().WriteToLog()
		return f, err
	}
	err = syscall.Flock(int(f.Fd()), syscall.LOCK_EX)
	if err != nil {
		newError(err).AtDebug().WriteToLog()
		err = f.Close()
		if err != nil {
			if ___DEBUG_PANIC_WHEN_ENCOUNTED_IMPOSSIBLE_ERROR {
				panic(err)
			}
			newError(err).AtDebug().WriteToLog()
		}
		return nil, err
	}
}

//DEBUG CONSTS
const ___DEBUG_IGNORE_FLOCK = false
const ___DEBUG_PANIC_WHEN_ERROR_UNPROPAGATEABLE = false
const ___DEBUG_PANIC_WHEN_ENCOUNTED_IMPOSSIBLE_ERROR = false
