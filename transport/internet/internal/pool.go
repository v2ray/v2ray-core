package internal

import (
	"net"
	"sync"
	"time"
	v2net "v2ray.com/core/common/net"
	"v2ray.com/core/common/signal"
)

type ConnectionId struct {
	Local      v2net.Address
	Remote     v2net.Address
	RemotePort v2net.Port
}

func NewConnectionId(source v2net.Address, dest v2net.Destination) ConnectionId {
	return ConnectionId{
		Local:      source,
		Remote:     dest.Address,
		RemotePort: dest.Port,
	}
}

type ExpiringConnection struct {
	conn   net.Conn
	expire time.Time
}

func (o *ExpiringConnection) Expired() bool {
	return o.expire.Before(time.Now())
}

type Pool struct {
	sync.Mutex
	connsByDest map[ConnectionId][]*ExpiringConnection
	cleanupOnce signal.Once
}

func NewConnectionPool() *Pool {
	return &Pool{
		connsByDest: make(map[ConnectionId][]*ExpiringConnection),
	}
}

func (o *Pool) Get(id ConnectionId) net.Conn {
	o.Lock()
	defer o.Unlock()

	list, found := o.connsByDest[id]
	if !found {
		return nil
	}
	connIdx := -1
	for idx, conn := range list {
		if !conn.Expired() {
			connIdx = idx
			break
		}
	}
	if connIdx == -1 {
		return nil
	}
	listLen := len(list)
	conn := list[connIdx]
	if connIdx != listLen-1 {
		list[connIdx] = list[listLen-1]
	}
	list = list[:listLen-1]
	o.connsByDest[id] = list
	return conn.conn
}

func (o *Pool) Cleanup() {
	defer o.cleanupOnce.Reset()

	for len(o.connsByDest) > 0 {
		time.Sleep(time.Second * 5)
		expiredConns := make([]net.Conn, 0, 16)
		o.Lock()
		for dest, list := range o.connsByDest {
			validConns := make([]*ExpiringConnection, 0, len(list))
			for _, conn := range list {
				if conn.Expired() {
					expiredConns = append(expiredConns, conn.conn)
				} else {
					validConns = append(validConns, conn)
				}
			}
			if len(validConns) != len(list) {
				o.connsByDest[dest] = validConns
			}
		}
		o.Unlock()
		for _, conn := range expiredConns {
			conn.Close()
		}
	}
}

func (o *Pool) Put(id ConnectionId, conn net.Conn) {
	expiringConn := &ExpiringConnection{
		conn:   conn,
		expire: time.Now().Add(time.Second * 4),
	}

	o.Lock()
	defer o.Unlock()

	list, found := o.connsByDest[id]
	if !found {
		list = []*ExpiringConnection{expiringConn}
	} else {
		list = append(list, expiringConn)
	}
	o.connsByDest[id] = list

	o.cleanupOnce.Do(func() {
		go o.Cleanup()
	})
}
