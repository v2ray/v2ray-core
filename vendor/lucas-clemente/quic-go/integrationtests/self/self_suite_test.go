package self_test

import (
	"math/rand"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	_ "github.com/lucas-clemente/quic-go/integrationtests/tools/testlog"
)

func TestSelf(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Self integration tests")
}

var _ = BeforeSuite(func() {
	rand.Seed(GinkgoRandomSeed())
})
