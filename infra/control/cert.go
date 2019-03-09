package control

import (
	"context"
	"crypto/x509"
	"encoding/json"
	"flag"
	"os"
	"strings"
	"time"

	"v2ray.com/core/common"
	"v2ray.com/core/common/protocol/tls/cert"
	"v2ray.com/core/common/task"
)

type stringList []string

func (l *stringList) String() string {
	return "String list"
}

func (l *stringList) Set(v string) error {
	if len(v) == 0 {
		return newError("empty value")
	}
	*l = append(*l, v)
	return nil
}

type jsonCert struct {
	Certificate []string `json:"certificate"`
	Key         []string `json:"key"`
}

type CertificateCommand struct {
}

func (c *CertificateCommand) Name() string {
	return "cert"
}

func (c *CertificateCommand) Description() Description {
	return Description{
		Short: "Generate TLS certificates.",
		Usage: []string{
			"v2ctl cert [--ca] [--domain=v2ray.com] [--expire=240h]",
			"Generate new TLS certificate",
			"--ca The new certificate is a CA certificate",
			"--domain Common name for the certificate",
			"--exipre Time until certificate expires. 240h = 10 days.",
		},
	}
}

func (c *CertificateCommand) printJson(certificate *cert.Certificate) {
	certPEM, keyPEM := certificate.ToPEM()
	jCert := &jsonCert{
		Certificate: strings.Split(strings.TrimSpace(string(certPEM)), "\n"),
		Key:         strings.Split(strings.TrimSpace(string(keyPEM)), "\n"),
	}
	content, err := json.MarshalIndent(jCert, "", "  ")
	common.Must(err)
	os.Stdout.Write(content)
	os.Stdout.WriteString("\n")
}

func (c *CertificateCommand) writeFile(content []byte, name string) error {
	f, err := os.Create(name)
	if err != nil {
		return err
	}
	defer f.Close()

	return common.Error2(f.Write(content))
}

func (c *CertificateCommand) printFile(certificate *cert.Certificate, name string) error {
	certPEM, keyPEM := certificate.ToPEM()
	return task.Run(context.Background(), func() error {
		return c.writeFile(certPEM, name+"_cert.pem")
	}, func() error {
		return c.writeFile(keyPEM, name+"_key.pem")
	})
}

func (c *CertificateCommand) Execute(args []string) error {
	fs := flag.NewFlagSet(c.Name(), flag.ContinueOnError)

	var domainNames stringList
	fs.Var(&domainNames, "domain", "Domain name for the certificate")

	commonName := fs.String("name", "V2Ray Inc", "The common name of this certificate")
	organization := fs.String("org", "V2Ray Inc", "Organization of the certificate")

	isCA := fs.Bool("ca", false, "Whether this certificate is a CA")
	jsonOutput := fs.Bool("json", true, "Print certificate in JSON format")
	fileOutput := fs.String("file", "", "Save certificate in file.")

	expire := fs.Duration("expire", time.Hour*24*90 /* 90 days */, "Time until the certificate expires. Default value 3 months.")

	if err := fs.Parse(args); err != nil {
		return err
	}

	var opts []cert.Option
	if *isCA {
		opts = append(opts, cert.Authority(*isCA))
		opts = append(opts, cert.KeyUsage(x509.KeyUsageCertSign|x509.KeyUsageKeyEncipherment|x509.KeyUsageDigitalSignature))
	}

	opts = append(opts, cert.NotAfter(time.Now().Add(*expire)))
	opts = append(opts, cert.CommonName(*commonName))
	if len(domainNames) > 0 {
		opts = append(opts, cert.DNSNames(domainNames...))
	}
	opts = append(opts, cert.Organization(*organization))

	cert, err := cert.Generate(nil, opts...)
	if err != nil {
		return newError("failed to generate TLS certificate").Base(err)
	}

	if *jsonOutput {
		c.printJson(cert)
	}

	if len(*fileOutput) > 0 {
		if err := c.printFile(cert, *fileOutput); err != nil {
			return err
		}
	}

	return nil
}

func init() {
	common.Must(RegisterCommand(&CertificateCommand{}))
}
