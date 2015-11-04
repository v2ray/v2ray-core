package main

import (
	"os"
	"testing"

	"github.com/v2ray/v2ray-core/testing/unit"
)

func TestBuildMacOS(t *testing.T) {
	assert := unit.Assert(t)

	targetFile := os.ExpandEnv("$GOPATH/bin/v2ray-macos.zip")
	os.Remove(targetFile)

	*targetOS = "macos"
	*targetArch = "amd64"
	*archive = true
	main()

	_, err := os.Stat(targetFile)
	assert.Error(err).IsNil()
}
