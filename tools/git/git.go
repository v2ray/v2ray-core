package git

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const (
	VersionUndefined = "undefined"
)

func getRepoRoot() string {
	GOPATH := os.Getenv("GOPATH")
	return filepath.Join(GOPATH, "src", "github.com", "v2ray", "v2ray-core")
}

func RevParse(args ...string) (string, error) {
	args = append([]string{"rev-parse"}, args...)
	cmd := exec.Command("git", args...)
	cmd.Dir = getRepoRoot()
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

func NameRev(args ...string) (string, error) {
	args = append([]string{"name-rev"}, args...)
	cmd := exec.Command("git", args...)
	cmd.Dir = getRepoRoot()
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

func RepoVersion(rev string) (string, error) {
	rev, err := RevParse(rev)
	if err != nil {
		return "", err
	}
	version, err := NameRev("name-rev", "--tags", "--name-only", rev)
	if err != nil {
		return "", err
	}
	if strings.HasSuffix(version, "^0") {
		version = version[:len(version)-2]
	}
	return version, nil
}

func RepoVersionHead() (string, error) {
	return RepoVersion("HEAD")
}
