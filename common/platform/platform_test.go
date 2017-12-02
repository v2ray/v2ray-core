package platform_test

import (
	"os"
	"path/filepath"
	"testing"

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
	assert(err, IsNil)

	loc := GetAssetLocation("t")
	assert(filepath.Dir(loc), Equals, filepath.Dir(exec))

	os.Setenv("v2ray.location.asset", "/v2ray")
	assert(GetAssetLocation("t"), Equals, "/v2ray/t")
}
