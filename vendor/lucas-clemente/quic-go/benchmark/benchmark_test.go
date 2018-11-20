package benchmark

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"math/rand"
	"net"

	quic "github.com/lucas-clemente/quic-go"
	_ "github.com/lucas-clemente/quic-go/integrationtests/tools/testlog"
	"github.com/lucas-clemente/quic-go/internal/protocol"
	"github.com/lucas-clemente/quic-go/internal/testdata"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func init() {
	var _ = Describe("Benchmarks", func() {
		dataLen := size * /* MB */ 1e6
		data := make([]byte, dataLen)
		rand.Seed(GinkgoRandomSeed())
		rand.Read(data) // no need to check for an error. math.Rand.Read never errors

		for i := range protocol.SupportedVersions {
			version := protocol.SupportedVersions[i]

			Context(fmt.Sprintf("with version %s", version), func() {
				Measure(fmt.Sprintf("transferring a %d MB file", size), func(b Benchmarker) {
					var ln quic.Listener
					serverAddr := make(chan net.Addr)
					handshakeChan := make(chan struct{})
					// start the server
					go func() {
						defer GinkgoRecover()
						var err error
						ln, err = quic.ListenAddr(
							"localhost:0",
							testdata.GetTLSConfig(),
							&quic.Config{Versions: []protocol.VersionNumber{version}},
						)
						Expect(err).ToNot(HaveOccurred())
						serverAddr <- ln.Addr()
						sess, err := ln.Accept()
						Expect(err).ToNot(HaveOccurred())
						// wait for the client to complete the handshake before sending the data
						// this should not be necessary, but due to timing issues on the CIs, this is necessary to avoid sending too many undecryptable packets
						<-handshakeChan
						str, err := sess.OpenStream()
						Expect(err).ToNot(HaveOccurred())
						_, err = str.Write(data)
						Expect(err).ToNot(HaveOccurred())
						err = str.Close()
						Expect(err).ToNot(HaveOccurred())
					}()

					// start the client
					addr := <-serverAddr
					sess, err := quic.DialAddr(
						addr.String(),
						&tls.Config{InsecureSkipVerify: true},
						&quic.Config{Versions: []protocol.VersionNumber{version}},
					)
					Expect(err).ToNot(HaveOccurred())
					close(handshakeChan)
					str, err := sess.AcceptStream()
					Expect(err).ToNot(HaveOccurred())

					buf := &bytes.Buffer{}
					// measure the time it takes to download the dataLen bytes
					// note we're measuring the time for the transfer, i.e. excluding the handshake
					runtime := b.Time("transfer time", func() {
						_, err := io.Copy(buf, str)
						Expect(err).NotTo(HaveOccurred())
					})
					Expect(buf.Bytes()).To(Equal(data))

					b.RecordValue("transfer rate [MB/s]", float64(dataLen)/1e6/runtime.Seconds())

					ln.Close()
					sess.Close()
				}, samples)
			})
		}
	})
}
