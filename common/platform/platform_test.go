package platform_test

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"v2ray.com/core/common"
	. "v2ray.com/core/common/platform"
)

func TestNormalizeEnvName(t *testing.T) {
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
		if v := NormalizeEnvName(test.input); v != test.output {
			t.Error("unexpected output: ", v, " want ", test.output)
		}
	}
}

func TestEnvFlag(t *testing.T) {
	if v := (EnvFlag{
		Name: "xxxxx.y",
	}.GetValueAsInt(10)); v != 10 {
		t.Error("env value: ", v)
	}
}

func TestGetAssetLocation(t *testing.T) {
	exec, err := os.Executable()
	common.Must(err)

	loc := GetAssetLocation("t")
	if filepath.Dir(loc) != filepath.Dir(exec) {
		t.Error("asset dir: ", loc, " not in ", exec)
	}

	os.Setenv("v2ray.location.asset", "/v2ray")
	if runtime.GOOS == "windows" {
		if v := GetAssetLocation("t"); v != "\\v2ray\\t" {
			t.Error("asset loc: ", v)
		}
	} else {
		if v := GetAssetLocation("t"); v != "/v2ray/t" {
			t.Error("asset loc: ", v)
		}
	}
}
