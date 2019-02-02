package platform_test

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"v2ray.com/core/common"
	. "v2ray.com/core/common/platform"
	. "v2ray.com/ext/assert"
)

func TestNormalizeEnvName(t *testing.T) {
	assert := With(t)

	cases := []struct {
		input  string
		output string
	}{
		{
			input:  "a",
			output: "A",
		},
		{
			input:  "a.a",
			output: "A_A",
		},
		{
			input:  "A.A.B",
			output: "A_A_B",
		},
	}
	for _, test := range cases {
		assert(NormalizeEnvName(test.input), Equals, test.output)
	}
}

func TestEnvFlag(t *testing.T) {
	assert := With(t)

	assert(EnvFlag{
		Name: "xxxxx.y",
	}.GetValueAsInt(10), Equals, 10)
}

func TestGetAssetLocation(t *testing.T) {
	assert := With(t)

	exec, err := os.Executable()
	common.Must(err)

	loc := GetAssetLocation("t")
	assert(filepath.Dir(loc), Equals, filepath.Dir(exec))

	os.Setenv("v2ray.location.asset", "/v2ray")
	if runtime.GOOS == "windows" {
		assert(GetAssetLocation("t"), Equals, "\\v2ray\\t")
	} else {
		assert(GetAssetLocation("t"), Equals, "/v2ray/t")
	}
}

func TestGetPluginLocation(t *testing.T) {
	assert := With(t)

	exec, err := os.Executable()
	common.Must(err)

	loc := GetPluginDirectory()
	assert(loc, Equals, filepath.Join(filepath.Dir(exec), "plugins"))

	os.Setenv("V2RAY_LOCATION_PLUGIN", "/v2ray")
	assert(GetPluginDirectory(), Equals, "/v2ray")
}
