package gquic_test

import (
	"bytes"
	"fmt"
	"os/exec"
	"strconv"
	"time"

	_ "github.com/lucas-clemente/quic-clients" // download clients
	"github.com/lucas-clemente/quic-go/integrationtests/tools/proxy"
	"github.com/lucas-clemente/quic-go/integrationtests/tools/testserver"
	"github.com/lucas-clemente/quic-go/internal/protocol"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gexec"
)

var _ = Describe("non-zero RTT", func() {
	var proxy *quicproxy.QuicProxy

	runRTTTest := func(rtt time.Duration, version protocol.VersionNumber) {
		var err error
		proxy, err = quicproxy.NewQuicProxy("localhost:", version, &quicproxy.Opts{
			RemoteAddr: "localhost:" + testserver.Port(),
			DelayPacket: func(_ quicproxy.Direction, _ uint64) time.Duration {
				return rtt / 2
			},
		})
		Expect(err).ToNot(HaveOccurred())

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

	AfterEach(func() {
		err := proxy.Close()
		Expect(err).ToNot(HaveOccurred())
		time.Sleep(time.Millisecond)
	})

	for i := range protocol.SupportedVersions {
		version := protocol.SupportedVersions[i]

		Context(fmt.Sprintf("with QUIC version %s", version), func() {
			roundTrips := [...]int{10, 50, 100, 200}
			for _, rtt := range roundTrips {
				It(fmt.Sprintf("gets a 500kB file with %dms RTT", rtt), func() {
					runRTTTest(time.Duration(rtt)*time.Millisecond, version)
				})
			}
		})
	}
})
