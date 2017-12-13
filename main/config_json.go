package main

import (
	"io"
	"os"
	"os/exec"

	"v2ray.com/core"
	"v2ray.com/core/common/platform"
)

func jsonToProto(input io.Reader) (*core.Config, error) {
	v2ctl := platform.GetToolLocation("v2ctl")
	_, err := os.Stat(v2ctl)
	if err != nil {
		return nil, err
	}
	cmd := exec.Command(v2ctl, "config")
	cmd.Stdin = input
	cmd.Stderr = os.Stderr
	cmd.SysProcAttr = getSysProcAttr()

	stdoutReader, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	defer stdoutReader.Close()

	if err := cmd.Start(); err != nil {
		return nil, err
	}

	config, err := core.LoadConfig(core.ConfigFormat_Protobuf, stdoutReader)

	cmd.Wait()

	return config, err
}

func init() {
	core.RegisterConfigLoader(core.ConfigFormat_JSON, func(input io.Reader) (*core.Config, error) {
		config, err := jsonToProto(input)
		if err != nil {
			return nil, newError("failed to execute v2ctl to convert config file.").Base(err)
		}
		return config, nil
	})
}
