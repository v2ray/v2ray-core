package self_test

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net"

	quic "github.com/lucas-clemente/quic-go"
	"github.com/lucas-clemente/quic-go/integrationtests/tools/testserver"
	"github.com/lucas-clemente/quic-go/internal/protocol"
	"github.com/lucas-clemente/quic-go/internal/testdata"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Connection ID lengths tests", func() {
	randomConnIDLen := func() int {
		return 4 + int(rand.Int31n(15))
	}

	runServer := func(conf *quic.Config) quic.Listener {
		GinkgoWriter.Write([]byte(fmt.Sprintf("Using %d byte connection ID for the server\n", conf.ConnectionIDLength)))
		ln, err := quic.ListenAddr("localhost:0", testdata.GetTLSConfig(), conf)
		Expect(err).ToNot(HaveOccurred())
		go func() {
			defer GinkgoRecover()
			for {
				sess, err := ln.Accept()
				if err != nil {
					return
				}
				go func() {
					defer GinkgoRecover()
					str, err := sess.OpenStream()
					Expect(err).ToNot(HaveOccurred())
					defer str.Close()
					_, err = str.Write(testserver.PRData)
					Expect(err).ToNot(HaveOccurred())
				}()
			}
		}()
		return ln
	}

	runClient := func(addr net.Addr, conf *quic.Config) {
		GinkgoWriter.Write([]byte(fmt.Sprintf("Using %d byte connection ID for the client\n", conf.ConnectionIDLength)))
		cl, err := quic.DialAddr(
			fmt.Sprintf("quic.clemente.io:%d", addr.(*net.UDPAddr).Port),
			&tls.Config{InsecureSkipVerify: true},
			conf,
		)
		Expect(err).ToNot(HaveOccurred())
		defer cl.Close()
		str, err := cl.AcceptStream()
		Expect(err).ToNot(HaveOccurred())
		data, err := ioutil.ReadAll(str)
		Expect(err).ToNot(HaveOccurred())
		Expect(data).To(Equal(testserver.PRData))
	}

	Context("IETF QUIC", func() {
		It("downloads a file using a 0-byte connection ID for the client", func() {
			serverConf := &quic.Config{
				ConnectionIDLength: randomConnIDLen(),
				Versions:           []protocol.VersionNumber{protocol.VersionTLS},
			}
			clientConf := &quic.Config{
				Versions: []protocol.VersionNumber{protocol.VersionTLS},
			}

			ln := runServer(serverConf)
			defer ln.Close()
			runClient(ln.Addr(), clientConf)
		})

		It("downloads a file when both client and server use a random connection ID length", func() {
			serverConf := &quic.Config{
				ConnectionIDLength: randomConnIDLen(),
				Versions:           []protocol.VersionNumber{protocol.VersionTLS},
			}
			clientConf := &quic.Config{
				ConnectionIDLength: randomConnIDLen(),
				Versions:           []protocol.VersionNumber{protocol.VersionTLS},
			}

			ln := runServer(serverConf)
			defer ln.Close()
			runClient(ln.Addr(), clientConf)
		})
	})

	Context("gQUIC", func() {
		It("downloads a file using a 0-byte connection ID for the client", func() {
			ln := runServer(&quic.Config{})
			defer ln.Close()
			runClient(ln.Addr(), &quic.Config{RequestConnectionIDOmission: true})
		})
	})
})
