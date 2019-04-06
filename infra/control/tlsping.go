package control

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"net"

	"v2ray.com/core/common"
)

type TlsPingCommand struct{}

func (c *TlsPingCommand) Name() string {
	return "tlsping"
}

func (c *TlsPingCommand) Description() Description {
	return Description{
		Short: "Ping the domain with TLS handshake",
		Usage: []string{"v2ctl tlsping <domain> --ip <ip>"},
	}
}

func printCertificates(certs []*x509.Certificate) {
	for _, cert := range certs {
		if len(cert.DNSNames) == 0 {
			continue
		}
		fmt.Println("Allowed domains: ", cert.DNSNames)
	}
}

func (c *TlsPingCommand) Execute(args []string) error {
	fs := flag.NewFlagSet(c.Name(), flag.ContinueOnError)
	ipStr := fs.String("ip", "", "IP address of the domain")

	if err := fs.Parse(args); err != nil {
		return newError("flag parsing").Base(err)
	}

	if fs.NArg() < 1 {
		return newError("domain not specified")
	}

	domain := fs.Arg(0)
	fmt.Println("Tls ping: ", domain)

	var ip net.IP
	if len(*ipStr) > 0 {
		v := net.ParseIP(*ipStr)
		if v == nil {
			return newError("invalid IP: ", *ipStr)
		}
		ip = v
	} else {
		v, err := net.ResolveIPAddr("ip", domain)
		if err != nil {
			return newError("resolve IP").Base(err)
		}
		ip = v.IP
	}
	fmt.Println("Using IP: ", ip.String())

	fmt.Println("-------------------")
	fmt.Println("Pinging without SNI")
	{
		tcpConn, err := net.DialTCP("tcp", nil, &net.TCPAddr{IP: ip, Port: 443})
		if err != nil {
			return newError("dial tcp").Base(err)
		}
		tlsConn := tls.Client(tcpConn, &tls.Config{
			InsecureSkipVerify: true,
			NextProtos:         []string{"http/1.1"},
			MaxVersion:         tls.VersionTLS12,
			MinVersion:         tls.VersionTLS12,
		})
		err = tlsConn.Handshake()
		if err != nil {
			fmt.Println("Handshake failure: ", err)
		} else {
			fmt.Println("Handshake succeeded")
			printCertificates(tlsConn.ConnectionState().PeerCertificates)
		}
		tlsConn.Close()
	}

	fmt.Println("-------------------")
	fmt.Println("Pinging with SNI")
	{
		tcpConn, err := net.DialTCP("tcp", nil, &net.TCPAddr{IP: ip, Port: 443})
		if err != nil {
			return newError("dial tcp").Base(err)
		}
		tlsConn := tls.Client(tcpConn, &tls.Config{
			ServerName: domain,
			NextProtos: []string{"http/1.1"},
			MaxVersion: tls.VersionTLS12,
			MinVersion: tls.VersionTLS12,
		})
		err = tlsConn.Handshake()
		if err != nil {
			fmt.Println("handshake failure: ", err)
		} else {
			fmt.Println("handshake succeeded")
			printCertificates(tlsConn.ConnectionState().PeerCertificates)
		}
		tlsConn.Close()
	}

	fmt.Println("Tls ping finished")

	return nil
}

func init() {
	common.Must(RegisterCommand(&TlsPingCommand{}))
}
