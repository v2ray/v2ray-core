package ws

import (
	"bufio"
	"io"
	"net"
	"sync"
	"time"

	"github.com/v2ray/v2ray-core/common/log"

	"github.com/gorilla/websocket"
)

type wsconn struct {
	wsc         *websocket.Conn
	readBuffer  *bufio.Reader
	connClosing bool
	reusable    bool
	retloc      *sync.Cond
	rlock       *sync.Mutex
	wlock       *sync.Mutex
}

func (ws *wsconn) Read(b []byte) (n int, err error) {

	//defer ws.rlock.Unlock()
	//ws.checkifRWAfterClosing()
	if ws.connClosing {

		return 0, io.EOF
	}
	getNewBuffer := func() error {
		_, r, err := ws.wsc.NextReader()
		if err != nil {
			log.Warning("WS transport: ws connection NewFrameReader return " + err.Error())
			ws.connClosing = true
			ws.Close()
			return err
		}
		ws.readBuffer = bufio.NewReader(r)
		return nil
	}

	readNext := func(b []byte) (n int, err error) {
		if ws.readBuffer == nil {
			err = getNewBuffer()
			if err != nil {
				//ws.Close()
				return 0, err
			}
		}

		n, err = ws.readBuffer.Read(b)

		if err == nil {
			return n, err
		}

		if err == io.EOF {
			ws.readBuffer = nil
			if n == 0 {
				return ws.Read(b)
			}
			return n, nil
		}
		//ws.Close()
		return n, err

	}
	n, err = readNext(b)

	return n, err

}

func (ws *wsconn) Write(b []byte) (n int, err error) {

	//defer
	//ws.checkifRWAfterClosing()
	if ws.connClosing {

		return 0, io.EOF
	}
	writeWs := func(b []byte) (n int, err error) {
		wr, err := ws.wsc.NextWriter(websocket.BinaryMessage)
		if err != nil {
			log.Warning("WS transport: ws connection NewFrameReader return " + err.Error())
			ws.connClosing = true
			ws.Close()
			return 0, err
		}
		n, err = wr.Write(b)
		if err != nil {
			//ws.Close()
			return 0, err
		}
		err = wr.Close()
		if err != nil {
			//ws.Close()
			return 0, err
		}
		return n, err
	}
	n, err = writeWs(b)
	return n, err
}
func (ws *wsconn) Close() error {
	ws.connClosing = true
	err := ws.wsc.Close()
	ws.retloc.Broadcast()
	return err
}
func (ws *wsconn) LocalAddr() net.Addr {
	return ws.wsc.LocalAddr()
}
func (ws *wsconn) RemoteAddr() net.Addr {
	return ws.wsc.RemoteAddr()
}
func (ws *wsconn) SetDeadline(t time.Time) error {
	return func() error {
		errr := ws.SetReadDeadline(t)
		errw := ws.SetWriteDeadline(t)
		if errr == nil || errw == nil {
			return nil
		}
		if errr != nil {
			return errr
		}

		return errw
	}()
}
func (ws *wsconn) SetReadDeadline(t time.Time) error {
	return ws.wsc.SetReadDeadline(t)
}
func (ws *wsconn) SetWriteDeadline(t time.Time) error {
	return ws.wsc.SetWriteDeadline(t)
}

func (ws *wsconn) checkifRWAfterClosing() {
	if ws.connClosing {
		log.Error("WS transport: Read or Write After Conn have been marked closing, this can be dangerous.")
		//panic("WS transport: Read or Write After Conn have been marked closing. Please report this crash to developer.")
	}
}

func (ws *wsconn) setup() {
	ws.connClosing = false

	ws.rlock = &sync.Mutex{}
	ws.wlock = &sync.Mutex{}

	initConnectedCond := func() {
		rsl := &sync.Mutex{}
		ws.retloc = sync.NewCond(rsl)
	}

	initConnectedCond()
	//ws.pingPong()
}

func (ws *wsconn) Reusable() bool {
	return ws.reusable && !ws.connClosing
}

func (ws *wsconn) SetReusable(reusable bool) {
	if !effectiveConfig.ConnectionReuse {
		return
	}
	ws.reusable = reusable
}

func (ws *wsconn) pingPong() {
	pongRcv := make(chan int, 0)
	ws.wsc.SetPongHandler(func(data string) error {
		pongRcv <- 0
		return nil
	})

	go func() {
		for !ws.connClosing {
			ws.wsc.WriteMessage(websocket.PingMessage, nil)
			tick := time.NewTicker(time.Second * 3)

			select {
			case <-pongRcv:
				break
			case <-tick.C:
				ws.Close()
			}
			<-tick.C
			tick.Stop()
		}

		return
	}()

}
