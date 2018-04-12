// +build !windows

package domainsocket

import (
	"context"
	gotls "crypto/tls"
	"os"
	"strings"
	"syscall"

	"v2ray.com/core/common"
	"v2ray.com/core/common/net"
	"v2ray.com/core/transport/internet"
	"v2ray.com/core/transport/internet/tls"
)

type Listener struct {
	addr      *net.UnixAddr
	ln        net.Listener
	tlsConfig *gotls.Config
	config    *Config
	addConn   internet.ConnHandler
	locker    *fileLocker
}

func Listen(ctx context.Context, address net.Address, port net.Port, handler internet.ConnHandler) (internet.Listener, error) {
	settings := getSettingsFromContext(ctx)
	if settings == nil {
		return nil, newError("domain socket settings not specified.")
	}

	addr, err := settings.GetUnixAddr()
	if err != nil {
		return nil, err
	}

	unixListener, err := net.ListenUnix("unix", addr)
	if err != nil {
		return nil, newError("failed to listen domain socket").Base(err).AtWarning()
	}

	ln := &Listener{
		addr:    addr,
		ln:      unixListener,
		config:  settings,
		addConn: handler,
	}

	if !settings.Abstract {
		ln.locker = &fileLocker{
			path: settings.Path + ".lock",
		}
		if err := ln.locker.Acquire(); err != nil {
			unixListener.Close()
			return nil, err
		}
	}

	if config := tls.ConfigFromContext(ctx); config != nil {
		ln.tlsConfig = config.GetTLSConfig()
	}

	go ln.run()

	return ln, nil
}

func (ln *Listener) Addr() net.Addr {
	return ln.addr
}

func (ln *Listener) Close() error {
	if ln.locker != nil {
		ln.locker.Release()
	}
	return ln.ln.Close()
}

func (ln *Listener) run() {
	for {
		conn, err := ln.ln.Accept()
		if err != nil {
			if strings.Contains(err.Error(), "closed") {
				break
			}
			newError("failed to accepted raw connections").Base(err).AtWarning().WriteToLog()
			continue
		}

		if ln.tlsConfig != nil {
			conn = tls.Server(conn, ln.tlsConfig)
		}

		ln.addConn(internet.Connection(conn))
	}
}

type fileLocker struct {
	path string
	file *os.File
}

func (fl *fileLocker) Acquire() error {
	f, err := os.Create(fl.path)
	if err != nil {
		return err
	}
	if err := syscall.Flock(int(f.Fd()), syscall.LOCK_EX); err != nil {
		f.Close()
		return newError("failed to lock file: ", fl.path).Base(err)
	}
	fl.file = f
	return nil
}

func (fl *fileLocker) Release() {
	if err := syscall.Flock(int(fl.file.Fd()), syscall.LOCK_UN); err != nil {
		newError("failed to unlock file: ", fl.path).Base(err).WriteToLog()
	}
	if err := fl.file.Close(); err != nil {
		newError("failed to close file: ", fl.path).Base(err).WriteToLog()
	}
	if err := os.Remove(fl.path); err != nil {
		newError("failed to remove file: ", fl.path).Base(err).WriteToLog()
	}
}

func init() {
	common.Must(internet.RegisterTransportListener(internet.TransportProtocol_DomainSocket, Listen))
}
