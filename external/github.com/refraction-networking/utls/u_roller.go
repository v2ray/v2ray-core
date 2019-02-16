package tls

import (
	"net"
	"sync"
	"time"
)

type Roller struct {
	HelloIDs            []ClientHelloID
	HelloIDMu           sync.Mutex
	WorkingHelloID      *ClientHelloID
	TcpDialTimeout      time.Duration
	TlsHandshakeTimeout time.Duration
}

// NewRoller creates Roller object with default range of HelloIDs to cycle through until a
// working/unblocked one is found.
func NewRoller() (*Roller, error) {
	tcpDialTimeoutInc, err := getRandInt(14)
	if err != nil {
		return nil, err
	}
	tcpDialTimeoutInc = 7 + tcpDialTimeoutInc

	tlsHandshakeTimeoutInc, err := getRandInt(20)
	if err != nil {
		return nil, err
	}
	tlsHandshakeTimeoutInc = 11 + tlsHandshakeTimeoutInc

	return &Roller{
		HelloIDs: []ClientHelloID{
			HelloChrome_Auto,
			HelloFirefox_Auto,
			HelloIOS_Auto,
			HelloRandomized,
		},
		TcpDialTimeout:      time.Second * time.Duration(tcpDialTimeoutInc),
		TlsHandshakeTimeout: time.Second * time.Duration(tlsHandshakeTimeoutInc),
	}, nil
}

// Dial attempts to establish connection to given address using different HelloIDs.
// If a working HelloID is found, it is used again for subsequent Dials.
// If tcp connection fails or all HelloIDs are tried, returns with last error.
//
// Usage examples:
//    Dial("tcp4", "google.com:443", "google.com")
//    Dial("tcp", "10.23.144.22:443", "mywebserver.org")
func (c *Roller) Dial(network, addr, serverName string) (*UConn, error) {
	helloIDs, err := shuffleClientHelloIDs(c.HelloIDs)
	if err != nil {
		return nil, err
	}

	c.HelloIDMu.Lock()
	workingHelloId := c.WorkingHelloID // keep using same helloID, if it works
	c.HelloIDMu.Unlock()
	if workingHelloId != nil {
		for i, ID := range helloIDs {
			if ID == *workingHelloId {
				helloIDs[i] = helloIDs[0]
				helloIDs[0] = *workingHelloId // push working hello ID first
				break
			}
		}
	}

	var tcpConn net.Conn
	for _, helloID := range helloIDs {
		tcpConn, err = net.DialTimeout(network, addr, c.TcpDialTimeout)
		if err != nil {
			return nil, err // on tcp Dial failure return with error right away
		}

		client := UClient(tcpConn, nil, helloID)
		client.SetSNI(serverName)
		client.SetDeadline(time.Now().Add(c.TlsHandshakeTimeout))
		err = client.Handshake()
		client.SetDeadline(time.Time{}) // unset timeout
		if err != nil {
			continue // on tls Dial error keep trying HelloIDs
		}

		c.HelloIDMu.Lock()
		c.WorkingHelloID = &helloID
		c.HelloIDMu.Unlock()
		return client, err
	}
	return nil, err
}

// returns a shuffled copy of input
func shuffleClientHelloIDs(helloIDs []ClientHelloID) ([]ClientHelloID, error) {
	perm, err := getRandPerm(len(helloIDs))
	if err != nil {
		return nil, err
	}

	shuffled := make([]ClientHelloID, len(helloIDs))
	for i, randI := range perm {
		shuffled[i] = helloIDs[randI]
	}
	return shuffled, nil
}
