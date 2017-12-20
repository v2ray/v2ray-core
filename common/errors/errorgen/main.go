// +build generate

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"v2ray.com/core/common"
)

var (
	pkg  = flag.String("pkg", "", "Target package")
	path = flag.String("path", "", "Path")
)

func main() {
	flag.Parse()

	if len(*pkg) == 0 {
		panic("Package is not specified.")
	}

	if len(*path) == 0 {
		panic("Path is not specified.")
	}

	paths := strings.Split(*path, ",")
	for i := range paths {
		paths[i] = "\"" + paths[i] + "\""
	}
	pathStr := strings.Join(paths, ", ")

	file, err := os.OpenFile("errors.generated.go", os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		log.Fatalf("Failed to generate errors.generated.go: %v", err)
	}

	common.Must2(fmt.Fprintln(file, "package", *pkg))
	common.Must2(fmt.Fprintln(file, ""))
	common.Must2(fmt.Fprintln(file, "import \"v2ray.com/core/common/errors\""))
	common.Must2(fmt.Fprintln(file, ""))
	common.Must2(fmt.Fprintln(file, "func newError(values ...interface{}) *errors.Error { return errors.New(values...).Path("+pathStr+") }"))

	common.Must(file.Close())
}
