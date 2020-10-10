package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"v2ray.com/core/common"
)

func main() {
	pwd, err := os.Getwd()
	if err != nil {
		fmt.Println("can not get current working directory")
		os.Exit(1)
	}
	pkg := filepath.Base(pwd)
	if pkg == "v2ray-core" {
		pkg = "core"
	}

	moduleName, gmnErr := common.GetModuleName(pwd)
	if gmnErr != nil {
		fmt.Println("can not get module path", gmnErr)
		os.Exit(1)
	}

	file, err := os.OpenFile("errors.generated.go", os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		log.Fatalf("Failed to generate errors.generated.go: %v", err)
		os.Exit(1)
	}
	defer file.Close()

	fmt.Fprintln(file, "package", pkg)
	fmt.Fprintln(file, "")
	fmt.Fprintln(file, "import \""+moduleName+"/common/errors\"")
	fmt.Fprintln(file, "")
	fmt.Fprintln(file, "type errPathObjHolder struct{}")
	fmt.Fprintln(file, "")
	fmt.Fprintln(file, "func newError(values ...interface{}) *errors.Error {")
	fmt.Fprintln(file, "	return errors.New(values...).WithPathObj(errPathObjHolder{})")
	fmt.Fprintln(file, "}")
}
