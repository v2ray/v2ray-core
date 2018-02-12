package main

import (
	"context"
	"io"
	"os"
	"os/exec"

	"v2ray.com/core"
	"v2ray.com/core/common"
	"v2ray.com/core/common/platform"
	"v2ray.com/core/common/signal"
)

type logWriter struct{}

func (*logWriter) Write(b []byte) (int, error) {
	n, err := os.Stderr.Write(b)
	if err == nil {
		os.Stderr.WriteString(platform.LineSeparator())
	}
	return n, err
}

func jsonToProto(input io.Reader) (*core.Config, error) {
	v2ctl := platform.GetToolLocation("v2ctl")
	if _, err := os.Stat(v2ctl); err != nil {
		return nil, err
	}
	cmd := exec.Command(v2ctl, "config")
	cmd.Stdin = input
	cmd.Stderr = &logWriter{}
	cmd.SysProcAttr = getSysProcAttr()

	stdoutReader, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	defer stdoutReader.Close()

	if err := cmd.Start(); err != nil {
		return nil, err
	}

	var config *core.Config

	loadTask := signal.ExecuteAsync(func() error {
		c, err := core.LoadConfig(core.ConfigFormat_Protobuf, stdoutReader)
		if err != nil {
			return err
		}
		config = c
		return nil
	})

	waitTask := signal.ExecuteAsync(func() error {
		return cmd.Wait()
	})

	if err := signal.ErrorOrFinish2(context.Background(), loadTask, waitTask); err != nil {
		return nil, err
	}

	return config, nil
}

func init() {
	common.Must(core.RegisterConfigLoader(core.ConfigFormat_JSON, func(input io.Reader) (*core.Config, error) {
		config, err := jsonToProto(input)
		if err != nil {
			return nil, newError("failed to execute v2ctl to convert config file.").Base(err).AtWarning()
		}
		return config, nil
	}))
}
