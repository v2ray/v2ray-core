package self_test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"time"

	quic "github.com/lucas-clemente/quic-go"
	"github.com/lucas-clemente/quic-go/h2quic"
	"github.com/lucas-clemente/quic-go/integrationtests/tools/testserver"
	"github.com/lucas-clemente/quic-go/internal/protocol"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
)

var _ = Describe("Client tests", func() {
	var client *http.Client

	// also run some tests with the TLS handshake
	versions := append(protocol.SupportedVersions, protocol.VersionTLS)

	BeforeEach(func() {
		err := os.Setenv("HOSTALIASES", "quic.clemente.io 127.0.0.1")
		Expect(err).ToNot(HaveOccurred())
		addr, err := net.ResolveUDPAddr("udp4", "quic.clemente.io:0")
		Expect(err).ToNot(HaveOccurred())
		if addr.String() != "127.0.0.1:0" {
			Fail("quic.clemente.io does not resolve to 127.0.0.1. Consider adding it to /etc/hosts.")
		}
		testserver.StartQuicServer(versions)
	})

	AfterEach(func() {
		testserver.StopQuicServer()
	})

	for _, v := range versions {
		version := v

		Context(fmt.Sprintf("with QUIC version %s", version), func() {
			BeforeEach(func() {
				client = &http.Client{
					Transport: &h2quic.RoundTripper{
						QuicConfig: &quic.Config{
							Versions: []protocol.VersionNumber{version},
						},
					},
				}
			})

			It("downloads a hello", func() {
				resp, err := client.Get("https://quic.clemente.io:" + testserver.Port() + "/hello")
				Expect(err).ToNot(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(200))
				body, err := ioutil.ReadAll(gbytes.TimeoutReader(resp.Body, 3*time.Second))
				Expect(err).ToNot(HaveOccurred())
				Expect(string(body)).To(Equal("Hello, World!\n"))
			})

			It("downloads a small file", func() {
				resp, err := client.Get("https://quic.clemente.io:" + testserver.Port() + "/prdata")
				Expect(err).ToNot(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(200))
				body, err := ioutil.ReadAll(gbytes.TimeoutReader(resp.Body, 5*time.Second))
				Expect(err).ToNot(HaveOccurred())
				Expect(body).To(Equal(testserver.PRData))
			})

			It("downloads a large file", func() {
				resp, err := client.Get("https://quic.clemente.io:" + testserver.Port() + "/prdatalong")
				Expect(err).ToNot(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(200))
				body, err := ioutil.ReadAll(gbytes.TimeoutReader(resp.Body, 20*time.Second))
				Expect(err).ToNot(HaveOccurred())
				Expect(body).To(Equal(testserver.PRDataLong))
			})

			It("uploads a file", func() {
				resp, err := client.Post(
					"https://quic.clemente.io:"+testserver.Port()+"/echo",
					"text/plain",
					bytes.NewReader(testserver.PRData),
				)
				Expect(err).ToNot(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(200))
				body, err := ioutil.ReadAll(gbytes.TimeoutReader(resp.Body, 5*time.Second))
				Expect(err).ToNot(HaveOccurred())
				Expect(bytes.Equal(body, testserver.PRData)).To(BeTrue())
			})
		})
	}
})
