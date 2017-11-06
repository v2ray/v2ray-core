package platform

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type EnvFlag struct {
	Name    string
	AltName string
}

func (f EnvFlag) GetValue(defaultValue string) string {
	if v, found := os.LookupEnv(f.Name); found {
		return v
	}
	if len(f.AltName) > 0 {
		if v, found := os.LookupEnv(f.AltName); found {
			return v
		}
	}

	return defaultValue
}

func (f EnvFlag) GetValueAsInt(defaultValue int) int {
	const PlaceHolder = "xxxxxx"
	s := f.GetValue(PlaceHolder)
	if s == PlaceHolder {
		return defaultValue
	}
	v, err := strconv.ParseInt(s, 10, 32)
	if err != nil {
		return defaultValue
	}
	return int(v)
}

func NormalizeEnvName(name string) string {
	return strings.Replace(strings.ToUpper(strings.TrimSpace(name)), ".", "_", -1)
}

var assetPath = "/"

func loadAssetLocation() {
	defAssetLocation, err := os.Executable()
	if err == nil {
		defAssetLocation = filepath.Dir(defAssetLocation)
		assetPath = (EnvFlag{
			Name: "v2ray.location.asset",
		}).GetValue(defAssetLocation)
	}
}

func init() {
	loadAssetLocation()
}

/*ForceReevaluate Force V2Ray to reevaluate environment Var*/
func ForceReevaluate() {
	loadAssetLocation()
}

func GetAssetLocation(file string) string {
	return filepath.Join(assetPath, file)
}
