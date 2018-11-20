package gquic_test

import (
	"fmt"
	"math/rand"
	"path/filepath"
	"runtime"

	_ "github.com/lucas-clemente/quic-go/integrationtests/tools/testlog"
	"github.com/lucas-clemente/quic-go/integrationtests/tools/testserver"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

var (
	clientPath string
	serverPath string
)

func TestIntegration(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "GQuic Tests Suite")
}

var _ = BeforeSuite(func() {
	rand.Seed(GinkgoRandomSeed())
})

var _ = JustBeforeEach(func() {
	testserver.StartQuicServer(nil)
})

var _ = AfterEach(testserver.StopQuicServer)

func init() {
	_, thisfile, _, ok := runtime.Caller(0)
	if !ok {
		panic("Failed to get current path")
	}
	clientPath = filepath.Join(thisfile, fmt.Sprintf("../../../../quic-clients/client-%s-debug", runtime.GOOS))
	serverPath = filepath.Join(thisfile, fmt.Sprintf("../../../../quic-clients/server-%s-debug", runtime.GOOS))
}
