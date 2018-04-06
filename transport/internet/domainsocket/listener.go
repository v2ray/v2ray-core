package domainsocket

import (
	"context"
	"net"
	"os"
	"syscall"
	"time"

	"v2ray.com/core/common/bitmask"
)

type Listener struct {
	ln           net.Listener
	listenerChan chan<- net.Conn
	ctx          context.Context
	path         string
	lockfile     *os.File
	state        bitmask.Byte
	cancal       func()
}

const (
	STATE_UNDEFINED   = 0
	STATE_INITIALIZED = 1 << iota
	STATE_LOWERUP     = 1 << iota
	STATE_UP          = 1 << iota
	STATE_TAINT       = 1 << iota
)

func ListenDS(ctx context.Context, path string) (*Listener, error) {

	vln := &Listener{path: path, state: STATE_INITIALIZED, ctx: ctx}
	return vln, nil
}

func (ls *Listener) Down() error {
	var err error
	if !ls.state.Has(STATE_LOWERUP | STATE_UP) {
		err = newError(ls.state).Base(newError("Invalid State:Down"))
		if ___DEBUG_PANIC_WHEN_ENCOUNTED_IMPOSSIBLE_ERROR {
			panic(err)
		}
		return err
	}

	ls.cancal()
	closeerr := ls.ln.Close()
	var lockerr error
	if isUnixDomainSocketFileSystemBased(ls.path) {
		lockerr = giveupLock(ls.lockfile)
	}
	if closeerr != nil && lockerr != nil {
		if ___DEBUG_PANIC_WHEN_ERROR_UNPROPAGATEABLE {
			panic(closeerr.Error() + lockerr.Error())
		}
	}

	if closeerr != nil {
		return newError("Cannot Close Unix domain socket listener").Base(closeerr)
	}
	if lockerr != nil {
		return newError("Cannot release lock for Unix domain socket listener").Base(lockerr)
	}
	ls.state.Clear(STATE_LOWERUP | STATE_UP)
	return nil
}

//LowerUP Setup systen level Listener
func (ls *Listener) LowerUP() error {
	var err error

	if !ls.state.Has(STATE_INITIALIZED) || ls.state.Has(STATE_LOWERUP) {
		err = newError(ls.state).Base(newError("Invalid State:LowerUP"))
		if ___DEBUG_PANIC_WHEN_ENCOUNTED_IMPOSSIBLE_ERROR {
			panic(err)
		}
		return err
	}
	
	//If the unix domain socket is filesystem based, an file lock must be used to claim the right for listening on respective file.
	//https://gavv.github.io/blog/unix-socket-reuse/
	if isUnixDomainSocketFileSystemBased(ls.path) && !___DEBUG_IGNORE_FLOCK {
		ls.lockfile, err = acquireLock(ls.path + ".lock")
		if err != nil {
			newError(err).AtDebug().WriteToLog()
			return newError("Unable to acquire lock for filesystem based unix domain socket").Base(err)
		}

		err = cleansePath(ls.path)
		if err != nil {
			return newError("Unable to cleanse path for the creation of unix domain socket").Base(err)
		}
	}

	addr := new(net.UnixAddr)
	addr.Name = ls.path
	addr.Net = "unix"
	li, err := net.ListenUnix("unix", addr)
	ls.ln = li
	if err != nil {
		return newError("Unable to listen unix domain socket").Base(err)
	}

	ls.state.Set(STATE_LOWERUP)

	return nil
}

func (ls *Listener) UP(listener chan<- net.Conn, allowkick bool) error {
	var err error
	if !ls.state.Has(STATE_INITIALIZED|STATE_LOWERUP) || (ls.state.Has(STATE_UP) && !allowkick) {
		err = newError(ls.state).Base(newError("Invalid State:UP"))
		if ___DEBUG_PANIC_WHEN_ENCOUNTED_IMPOSSIBLE_ERROR {
			panic(err)
		}
		return err
	}
	ls.listenerChan = listener
	if !ls.state.Has(STATE_UP) {
		cctx, cancel := context.WithCancel(ls.ctx)
		ls.cancal = cancel
		go ls.uploop(cctx)
	}
	ls.state.Set(STATE_UP)
	return nil
}

func (ls *Listener) uploop(cctx context.Context) {
	var lasterror error
	errortolerance := 5
	for {
		if cctx.Err() != nil {
			close(ls.listenerChan)
			return
		}
		conn, err := ls.ln.Accept()

		if err != nil {
			newError("Cannot Accept socket from listener").Base(err).AtDebug().WriteToLog()
			//Guard against too many open file error
			if err == lasterror {
				errortolerance--
				if errortolerance == 0 {
					newError("unix domain socket melt down as the error is repeating").Base(err).AtError().WriteToLog()
					ls.cancal()
				}
				newError("unix domain socket listener is throttling accept as the error is repeating").Base(err).AtError().WriteToLog()
				time.Sleep(time.Second * 5)
			}
			lasterror = err
		}

		ls.listenerChan <- conn
	}
}

func isUnixDomainSocketFileSystemBased(path string) bool {
	//No Branching
	return path[0] != 0
}

func acquireLock(lockfilepath string) (*os.File, error) {
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
	return nil, err
}

func giveupLock(locker *os.File) error {
	err := syscall.Flock(int(locker.Fd()), syscall.LOCK_UN)
	if err != nil {
		closeerr := locker.Close()
		if err != nil {
			if ___DEBUG_PANIC_WHEN_ERROR_UNPROPAGATEABLE {
				panic(closeerr)
			}
			newError(closeerr).AtDebug().WriteToLog()
		}
		newError(err).AtDebug().WriteToLog()
		return err
	}
	closeerr := locker.Close()
	if closeerr != nil {
		newError(closeerr).AtDebug().WriteToLog()
		return closeerr
	}
	return closeerr
}

func cleansePath(path string) error {
	_, err := os.Stat(path)
	if err != nil {
		return nil
	}
	err = os.Remove(path)
	return err
}

//DEBUG CONSTS
const ___DEBUG_IGNORE_FLOCK = false
const ___DEBUG_PANIC_WHEN_ERROR_UNPROPAGATEABLE = false
const ___DEBUG_PANIC_WHEN_ENCOUNTED_IMPOSSIBLE_ERROR = false
