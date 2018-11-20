package benchmark

import (
	"flag"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestBenchmark(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Benchmark Suite")
}

var (
	size    int // file size in MB, will be read from flags
	samples int // number of samples for Measure, will be read from flags
)

func init() {
	flag.IntVar(&size, "size", 50, "data length (in MB)")
	flag.IntVar(&samples, "samples", 6, "number of samples")
	flag.Parse()
}
