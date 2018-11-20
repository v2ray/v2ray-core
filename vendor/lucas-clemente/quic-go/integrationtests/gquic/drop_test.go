package gquic_test

import (
	"bytes"
	"fmt"
	mrand "math/rand"
	"os/exec"
	"strconv"

	_ "github.com/lucas-clemente/quic-clients" // download clients
	"github.com/lucas-clemente/quic-go/integrationtests/tools/proxy"
	"github.com/lucas-clemente/quic-go/integrationtests/tools/testserver"
	"github.com/lucas-clemente/quic-go/internal/protocol"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	. "github.com/onsi/gomega/gexec"
)

var directions = []quicproxy.Direction{quicproxy.DirectionIncoming, quicproxy.DirectionOutgoing, quicproxy.DirectionBoth}

var _ = Describe("Drop tests", func() {
	var proxy *quicproxy.QuicProxy

	startProxy := func(dropCallback quicproxy.DropCallback, version protocol.VersionNumber) {
		var err error
		proxy, err = quicproxy.NewQuicProxy("localhost:0", version, &quicproxy.Opts{
			RemoteAddr: "localhost:" + testserver.Port(),
			DropPacket: dropCallback,
		})
		Expect(err).ToNot(HaveOccurred())
	}

	downloadFile := func(version protocol.VersionNumber) {
		command := exec.Command(
			clientPath,
			"--quic-version="+version.ToAltSvc(),
			"--host=127.0.0.1",
			"--port="+strconv.Itoa(proxy.LocalPort()),
			"https://quic.clemente.io/prdata",
		)
		session, err := Start(command, nil, GinkgoWriter)
		Expect(err).NotTo(HaveOccurred())
		defer session.Kill()
		Eventually(session, 20).Should(Exit(0))
		Expect(bytes.Contains(session.Out.Contents(), testserver.PRData)).To(BeTrue())
	}

	downloadHello := func(version protocol.VersionNumber) {
		command := exec.Command(
			clientPath,
			"--quic-version="+version.ToAltSvc(),
			"--host=127.0.0.1",
			"--port="+strconv.Itoa(proxy.LocalPort()),
			"https://quic.clemente.io/hello",
		)
		session, err := Start(command, nil, GinkgoWriter)
		Expect(err).NotTo(HaveOccurred())
		defer session.Kill()
		Eventually(session, 20).Should(Exit(0))
		Expect(session.Out).To(Say(":status 200"))
		Expect(session.Out).To(Say("body: Hello, World!\n"))
	}

	deterministicDropper := func(p, interval, dropInARow uint64) bool {
		return (p % interval) < dropInARow
	}

	stochasticDropper := func(freq int) bool {
		return mrand.Int63n(int64(freq)) == 0
	}

	AfterEach(func() {
		Expect(proxy.Close()).To(Succeed())
	})

	for _, v := range protocol.SupportedVersions {
		version := v

		Context(fmt.Sprintf("with QUIC version %s", version), func() {
			Context("during the crypto handshake", func() {
				for _, d := range directions {
					direction := d

					It(fmt.Sprintf("establishes a connection when the first packet is lost in %s direction", d), func() {
						startProxy(func(d quicproxy.Direction, p uint64) bool {
							return p == 1 && d.Is(direction)
						}, version)
						downloadHello(version)
					})

					It(fmt.Sprintf("establishes a connection when the second packet is lost in %s direction", d), func() {
						startProxy(func(d quicproxy.Direction, p uint64) bool {
							return p == 2 && d.Is(direction)
						}, version)
						downloadHello(version)
					})

					It(fmt.Sprintf("establishes a connection when 1/5 of the packets are lost in %s direction", d), func() {
						startProxy(func(d quicproxy.Direction, p uint64) bool {
							return d.Is(direction) && stochasticDropper(5)
						}, version)
						downloadHello(version)
					})
				}
			})

			Context("after the crypto handshake", func() {
				for _, d := range directions {
					direction := d

					It(fmt.Sprintf("downloads a file when every 5th packet is dropped in %s direction", d), func() {
						startProxy(func(d quicproxy.Direction, p uint64) bool {
							return p >= 10 && d.Is(direction) && deterministicDropper(p, 5, 1)
						}, version)
						downloadFile(version)
					})

					It(fmt.Sprintf("downloads a file when 1/5th of all packet are dropped randomly in %s direction", d), func() {
						startProxy(func(d quicproxy.Direction, p uint64) bool {
							return p >= 10 && d.Is(direction) && stochasticDropper(5)
						}, version)
						downloadFile(version)
					})

					It(fmt.Sprintf("downloads a file when 10 packets every 100 packet are dropped in %s direction", d), func() {
						startProxy(func(d quicproxy.Direction, p uint64) bool {
							return p >= 10 && d.Is(direction) && deterministicDropper(p, 100, 10)
						}, version)
						downloadFile(version)
					})
				}
			})
		})
	}
})
