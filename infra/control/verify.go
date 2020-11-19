package control

import (
	"flag"
	"github.com/xiaokangwang/VSign/signerVerify"
	"os"
	"v2ray.com/core/common"
)

type VerifyCommand struct{}

func (c *VerifyCommand) Name() string {
	return "verify"
}

func (c *VerifyCommand) Description() Description {
	return Description{
		Short: "Verify if a binary is officially signed.",
		Usage: []string{
			"v2ctl verify --sig=<sig-file> file...",
			"Verify the file officially signed by V2Ray.",
		},
	}
}

func (c *VerifyCommand) Execute(args []string) error {
	fs := flag.NewFlagSet(c.Name(), flag.ContinueOnError)

	sigFile := fs.String("sig", "", "Path to the signature file")

	if err := fs.Parse(args); err != nil {
		return err
	}

	target := fs.Arg(0)
	if target == "" {
		return newError("empty file path.")
	}

	if *sigFile == "" {
		return newError("empty signature path.")
	}

	sigReader, err := os.Open(os.ExpandEnv(*sigFile))
	if err != nil {
		return newError("failed to open file ", *sigFile).Base(err)
	}

	files := fs.Args()

	err = signerVerify.OutputAndJudge(signerVerify.CheckSignaturesV2Fly(sigReader, files))

	if err == nil {
		return nil
	}

	return newError("file is not officially signed by V2Ray").Base(err)
}

func init() {
	common.Must(RegisterCommand(&VerifyCommand{}))
}
