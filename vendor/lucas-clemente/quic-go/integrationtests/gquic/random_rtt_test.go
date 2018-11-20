package gquic_test

import (
	"bytes"
	"fmt"
	"math/rand"
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

// get a random duration between min and max
func getRandomDuration(min, max time.Duration) time.Duration {
	return min + time.Duration(rand.Int63n(int64(max-min)))
}

var _ = Describe("Random Duration Generator", func() {
	It("gets a random RTT", func() {
		var min time.Duration = time.Hour
		var max time.Duration

		var sum time.Duration
		rep := 10000
		for i := 0; i < rep; i++ {
			val := getRandomDuration(100*time.Millisecond, 500*time.Millisecond)
			sum += val
			if val < min {
				min = val
			}
			if val > max {
				max = val
			}
		}
		avg := sum / time.Duration(rep)
		Expect(avg).To(BeNumerically("~", 300*time.Millisecond, 5*time.Millisecond))
		Expect(min).To(BeNumerically(">=", 100*time.Millisecond))
		Expect(min).To(BeNumerically("<", 105*time.Millisecond))
		Expect(max).To(BeNumerically(">", 495*time.Millisecond))
		Expect(max).To(BeNumerically("<=", 500*time.Millisecond))
	})
})

var _ = Describe("Random RTT", func() {
	var proxy *quicproxy.QuicProxy

	runRTTTest := func(minRtt, maxRtt time.Duration, version protocol.VersionNumber) {
		var err error
		proxy, err = quicproxy.NewQuicProxy("localhost:", version, &quicproxy.Opts{
			RemoteAddr: "localhost:" + testserver.Port(),
			DelayPacket: func(_ quicproxy.Direction, _ uint64) time.Duration {
				return getRandomDuration(minRtt, maxRtt)
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
			It("gets a file a random RTT between 10ms and 30ms", func() {
				runRTTTest(10*time.Millisecond, 30*time.Millisecond, version)
			})
		})
	}
})
