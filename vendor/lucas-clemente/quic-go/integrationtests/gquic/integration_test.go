package gquic_test

import (
	"bytes"
	"fmt"
	"os/exec"
	"sync"

	"github.com/lucas-clemente/quic-go/integrationtests/tools/testserver"
	"github.com/lucas-clemente/quic-go/internal/protocol"

	_ "github.com/lucas-clemente/quic-clients" // download clients

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	. "github.com/onsi/gomega/gexec"
)

var _ = Describe("Integration tests", func() {
	for i := range protocol.SupportedVersions {
		version := protocol.SupportedVersions[i]

		Context(fmt.Sprintf("with QUIC version %s", version), func() {
			It("gets a simple file", func() {
				command := exec.Command(
					clientPath,
					"--quic-version="+version.ToAltSvc(),
					"--host=127.0.0.1",
					"--port="+testserver.Port(),
					"https://quic.clemente.io/hello",
				)
				session, err := Start(command, nil, GinkgoWriter)
				Expect(err).NotTo(HaveOccurred())
				defer session.Kill()
				Eventually(session, 5).Should(Exit(0))
				Expect(session.Out).To(Say(":status 200"))
				Expect(session.Out).To(Say("body: Hello, World!\n"))
			})

			It("posts and reads a body", func() {
				command := exec.Command(
					clientPath,
					"--quic-version="+version.ToAltSvc(),
					"--host=127.0.0.1",
					"--port="+testserver.Port(),
					"--body=foo",
					"https://quic.clemente.io/echo",
				)
				session, err := Start(command, nil, GinkgoWriter)
				Expect(err).NotTo(HaveOccurred())
				defer session.Kill()
				Eventually(session, 5).Should(Exit(0))
				Expect(session.Out).To(Say(":status 200"))
				Expect(session.Out).To(Say("body: foo\n"))
			})

			It("gets a file", func() {
				command := exec.Command(
					clientPath,
					"--quic-version="+version.ToAltSvc(),
					"--host=127.0.0.1",
					"--port="+testserver.Port(),
					"https://quic.clemente.io/prdata",
				)
				session, err := Start(command, nil, GinkgoWriter)
				Expect(err).NotTo(HaveOccurred())
				defer session.Kill()
				Eventually(session, 10).Should(Exit(0))
				Expect(bytes.Contains(session.Out.Contents(), testserver.PRData)).To(BeTrue())
			})

			It("gets many copies of a file in parallel", func() {
				wg := sync.WaitGroup{}
				for i := 0; i < 10; i++ {
					wg.Add(1)
					go func() {
						defer wg.Done()
						defer GinkgoRecover()
						command := exec.Command(
							clientPath,
							"--quic-version="+version.ToAltSvc(),
							"--host=127.0.0.1",
							"--port="+testserver.Port(),
							"https://quic.clemente.io/prdata",
						)
						session, err := Start(command, nil, GinkgoWriter)
						Expect(err).NotTo(HaveOccurred())
						defer session.Kill()
						Eventually(session, 20).Should(Exit(0))
						Expect(bytes.Contains(session.Out.Contents(), testserver.PRData)).To(BeTrue())
					}()
				}
				wg.Wait()
			})
		})
	}
})
