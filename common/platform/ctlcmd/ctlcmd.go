package ctlcmd

import (
	"context"
	"io"
	"os"
	"os/exec"

	"v2ray.com/core/common/buf"
	"v2ray.com/core/common/platform"
	"v2ray.com/core/common/signal"
)

//go:generate go run $GOPATH/src/v2ray.com/core/common/errors/errorgen/main.go -pkg ctlcmd -path Command,Platform,CtlCmd

func Run(args []string, input io.Reader) (buf.MultiBuffer, error) {
	v2ctl := platform.GetToolLocation("v2ctl")
	if _, err := os.Stat(v2ctl); err != nil {
		return nil, newError("v2ctl doesn't exist").Base(err)
	}

	errBuffer := &buf.MultiBuffer{}

	cmd := exec.Command(v2ctl, args...)
	cmd.Stderr = errBuffer
	cmd.SysProcAttr = getSysProcAttr()
	if input != nil {
		cmd.Stdin = input
	}

	stdoutReader, err := cmd.StdoutPipe()
	if err != nil {
		return nil, newError("failed to get stdout from v2ctl").Base(err)
	}
	defer stdoutReader.Close()

	if err := cmd.Start(); err != nil {
		return nil, newError("failed to start v2ctl").Base(err)
	}

	var content buf.MultiBuffer
	loadTask := func() error {
		c, err := buf.ReadAllToMultiBuffer(stdoutReader)
		if err != nil {
			return newError("failed to read config").Base(err)
		}
		content = c
		return nil
	}

	waitTask := func() error {
		if err := cmd.Wait(); err != nil {
			msg := "failed to execute v2ctl"
			if errBuffer.Len() > 0 {
				msg += ": " + errBuffer.String()
			}
			return newError(msg).Base(err)
		}
		return nil
	}

	if err := signal.ExecuteParallel(context.Background(), loadTask, waitTask); err != nil {
		return nil, err
	}

	return content, nil
}
