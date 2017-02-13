package internal

import (
	"net"
	"sync"
	"time"
)

// ConnectionRecyler is the interface for recycling connections.
type ConnectionRecyler interface {
	// Put returns a connection back to a connection pool.
	Put(ConnectionID, net.Conn)
}

type NoOpConnectionRecyler struct{}

func (NoOpConnectionRecyler) Put(ConnectionID, net.Conn) {}

// ExpiringConnection is a connection that will expire in certain time.
type ExpiringConnection struct {
	conn   net.Conn
	expire time.Time
}

// Expired returns true if the connection has expired.
func (ec *ExpiringConnection) Expired() bool {
	return ec.expire.Before(time.Now())
}

// Pool is a connection pool.
type Pool struct {
	sync.Mutex
	connsByDest  map[ConnectionID][]*ExpiringConnection
	cleanupToken chan bool
}

// NewConnectionPool creates a new Pool.
func NewConnectionPool() *Pool {
	p := &Pool{
		connsByDest:  make(map[ConnectionID][]*ExpiringConnection),
		cleanupToken: make(chan bool, 1),
	}
	p.cleanupToken <- true
	return p
}

// Get returns a connection with matching connection ID. Nil if not found.
func (p *Pool) Get(id ConnectionID) net.Conn {
	p.Lock()
	defer p.Unlock()

	list, found := p.connsByDest[id]
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
	p.connsByDest[id] = list
	return conn.conn
}

func (p *Pool) cleanup() {
	defer func() {
		p.cleanupToken <- true
	}()

	for len(p.connsByDest) > 0 {
		time.Sleep(time.Second * 5)
		expiredConns := make([]net.Conn, 0, 16)
		p.Lock()
		for dest, list := range p.connsByDest {
			validConns := make([]*ExpiringConnection, 0, len(list))
			for _, conn := range list {
				if conn.Expired() {
					expiredConns = append(expiredConns, conn.conn)
				} else {
					validConns = append(validConns, conn)
				}
			}
			if len(validConns) != len(list) {
				p.connsByDest[dest] = validConns
			}
		}
		p.Unlock()
		for _, conn := range expiredConns {
			conn.Close()
		}
	}
}

// Put implements ConnectionRecyler.Put().
func (p *Pool) Put(id ConnectionID, conn net.Conn) {
	expiringConn := &ExpiringConnection{
		conn:   conn,
		expire: time.Now().Add(time.Second * 4),
	}

	p.Lock()
	defer p.Unlock()

	list, found := p.connsByDest[id]
	if !found {
		list = []*ExpiringConnection{expiringConn}
	} else {
		list = append(list, expiringConn)
	}
	p.connsByDest[id] = list

	select {
	case <-p.cleanupToken:
		go p.cleanup()
	default:
	}
}
