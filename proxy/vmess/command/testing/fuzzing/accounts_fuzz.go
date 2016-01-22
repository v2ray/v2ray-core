// +build gofuzz

package fuzzing

import (
	. "github.com/v2ray/v2ray-core/proxy/vmess/command"
)

func Fuzz(data []byte) int {
	cmd := new(SwitchAccount)
	if err := cmd.Unmarshal(data); err != nil {
		return 0
	}
	return 1
}
