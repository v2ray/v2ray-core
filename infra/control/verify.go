package control

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"golang.org/x/crypto/openpgp"

	"v2ray.com/core/common"
)

const (
	pubkey = `-----BEGIN PGP PUBLIC KEY BLOCK-----
Comment: GPGTools - https://gpgtools.org

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
=xv/V
-----END PGP PUBLIC KEY BLOCK-----
`
)

func firstIdentity(m map[string]*openpgp.Identity) string {
	for k := range m {
		return k
	}
	return ""
}

type VerifyCommand struct{}

func (c *VerifyCommand) Name() string {
	return "verify"
}

func (c *VerifyCommand) Description() Description {
	return Description{
		Short: "Verify if a binary is officially signed.",
		Usage: []string{
			"v2ctl verify [--sig=<sig-file>] file",
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
	if len(target) == 0 {
		return newError("empty file path.")
	}

	if len(*sigFile) == 0 {
		*sigFile = target + ".sig"
	}

	targetReader, err := os.Open(os.ExpandEnv(target))
	if err != nil {
		return newError("failed to open file: ", target).Base(err)
	}

	sigReader, err := os.Open(os.ExpandEnv(*sigFile))
	if err != nil {
		return newError("failed to open file ", *sigFile).Base(err)
	}

	keyring, err := openpgp.ReadArmoredKeyRing(strings.NewReader(pubkey))
	if err != nil {
		return newError("failed to create keyring").Base(err)
	}

	entity, err := openpgp.CheckDetachedSignature(keyring, targetReader, sigReader)
	if err != nil {
		return newError("failed to verify signature").Base(err)
	}

	fmt.Println("Signed by:", firstIdentity(entity.Identities))
	return nil
}

func init() {
	common.Must(RegisterCommand(&VerifyCommand{}))
}
