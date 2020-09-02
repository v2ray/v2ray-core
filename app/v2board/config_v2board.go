package v2board

//go:generate errorgen

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strings"

	"os"

	"v2ray.com/core"
	"v2ray.com/core/common"
)

type V2BoardConfig struct {
	Server string `json:"server"`
	Node   int    `json:"node"`
	Token  string `json:"token"`
}

type V2BoardConfigRsp struct {
	Msg  string      `json:"msg"`
	Data core.Config `json:"data"`
}

func (v *V2Board) ConfigLoader(input interface{}) (*core.Config, error) {
	newError("Start ConfigLoader").AtInfo().WriteToLog()

	absPath, _ := filepath.Abs(filepath.Dir(os.Args[0]))

	buffer, err := ioutil.ReadFile(filepath.Join(absPath, "v2board.json"))
	common.Must(err)

	common.Must(json.Unmarshal(buffer, &v.config))

	uri := v.ConfigUri()
	// fmt.Println("config uri:", uri)

	resp, err := http.Get(uri)
	common.Must(err)
	defer resp.Body.Close()

	newError("Response status:", resp.Status).AtDebug().WriteToLog()
	buf := new(strings.Builder)
	io.Copy(buf, resp.Body)
	// check errors
	newError(buf.String()).AtDebug().WriteToLog()
	conf, err := core.LoadConfig("json", "", strings.NewReader(buf.String()))
	return conf, err
}
