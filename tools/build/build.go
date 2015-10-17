package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

var (
	targetOS   = flag.String("os", runtime.GOOS, "Target OS of this build.")
	targetArch = flag.String("arch", runtime.GOARCH, "Target CPU arch of this build.")
	archive    = flag.Bool("zip", false, "Whether to make an archive of files or not.")

	GOPATH   string
	repoRoot string
)

type OS string

const (
	Windows = OS("windows")
	MacOS   = OS("darwin")
	Linux   = OS("linux")
)

type Architecture string

const (
	X86   = Architecture("386")
	Amd64 = Architecture("amd64")
	Arm   = Architecture("arm")
	Arm64 = Architecture("arm64")
)

func getRepoVersion() (string, error) {
	revHead, err := exec.Command("git", "rev-parse", "HEAD").Output()
	if err != nil {
		return "", err
	}
	version, err := exec.Command("git", "name-rev", "--tags", "--name-only", strings.TrimSpace(string(revHead))).Output()
	if err != nil {
		return "", err
	}
	verStr := strings.TrimSpace(string(version))
	if strings.HasSuffix(verStr, "^0") {
		verStr = verStr[:len(verStr)-2]
	}
	if verStr == "undefined" {
		verStr = "custom"
	}
	return verStr, nil
}

func getOS() OS {
	if targetOS == nil {
		panic("OS is not specified.")
	}
	osStr := strings.ToLower(*targetOS)
	if osStr == "windows" || osStr == "win" {
		return Windows
	}
	if osStr == "darwin" || osStr == "mac" || osStr == "macos" || osStr == "osx" {
		return MacOS
	}
	if osStr == "linux" || osStr == "debian" || osStr == "ubuntu" || osStr == "redhat" || osStr == "centos" {
		return Linux
	}
	panic("Unknown OS " + *targetOS)
}

func getArch() Architecture {
	if targetArch == nil {
		panic("Arch is not specified.")
	}
	archStr := strings.ToLower(*targetArch)
	if archStr == "x86" || archStr == "386" || archStr == "i386" {
		return X86
	}
	if archStr == "amd64" || archStr == "x86-64" || archStr == "x64" {
		return Amd64
	}
	if archStr == "arm" {
		return Arm
	}
	if archStr == "arm64" {
		return Arm64
	}
	panic("Unknown Arch " + *targetArch)
}

func getSuffix(vos OS, arch Architecture) string {
	suffix := "-custom"
	switch vos {
	case Windows:
		switch arch {
		case X86:
			suffix = "-windows-32"
		case Amd64:
			suffix = "-windows-64"
		}
	case MacOS:
		suffix = "-macos"
	case Linux:
		switch arch {
		case X86:
			suffix = "-linux-32"
		case Amd64:
			suffix = "-linux-64"
		case Arm:
			suffix = "-linux-arm"
		case Arm64:
			suffix = "-linux-arm64"
		}

	}
	return suffix
}

func createTargetDirectory(version string, vos OS, arch Architecture) (string, error) {
	suffix := "-custom"
	if version != "custom" {
		suffix = getSuffix(vos, arch)
	}
	targetDir := filepath.Join(GOPATH, "bin", "v2ray"+suffix)
	if version != "custom" {
		os.RemoveAll(targetDir)
	}
	err := os.MkdirAll(targetDir, os.ModeDir|0777)
	return targetDir, err
}

func getTargetFile(vos OS, arch Architecture) string {
	suffix := getSuffix(vos, arch)
	if vos == "Windows" {
		suffix += ".exe"
	}
	return "v2ray" + suffix
}

func normalizedContent(content []byte, vos OS) []byte {
	str := strings.Replace(string(content), "\r\n", "\n", -1)
	if vos == Windows {
		str = strings.Replace(str, "\n", "\r\n", -1)
	}
	return []byte(str)
}

func copyConfigFiles(dir string, version string, vos OS, arch Architecture) error {
	clientConfig, err := ioutil.ReadFile(filepath.Join(repoRoot, "release", "config", "vpoint_socks_vmess.json"))
	if err != nil {
		return err
	}
	clientConfig = normalizedContent(clientConfig, vos)
	clientConfigName := "vpoint_socks_vmess.json"
	if vos == Windows || vos == MacOS {
		clientConfigName = "config.json"
	}
	err = ioutil.WriteFile(filepath.Join(dir, clientConfigName), clientConfig, 0777)
	if err != nil {
		return err
	}

	if vos == Windows || vos == MacOS {
		return nil
	}

	serverConfig, err := ioutil.ReadFile(filepath.Join(repoRoot, "release", "config", "vpoint_vmess_freedom.json"))
	if err != nil {
		return err
	}
	serverConfig = normalizedContent(serverConfig, vos)
	err = ioutil.WriteFile(filepath.Join(dir, "vpoint_vmess_freedom.json"), serverConfig, 0777)
	if err != nil {
		return err
	}
	return nil
}

func buildV2Ray(targetDir, targetFile string, version string, vos OS, arch Architecture) error {
	ldFlags := "-s"
	if version != "custom" {
		year, month, day := time.Now().UTC().Date()
		today := fmt.Sprintf("%04d%02d%02d", year, int(month), day)
		ldFlags = ldFlags + " -X github.com/v2ray/v2ray-core.version=" + version + " -X github.com/v2ray/v2ray-core.build=" + today
	}
	target := filepath.Join(targetDir, targetFile)
	fmt.Println("Building to " + target)
	cmd := exec.Command("go", "build", "-o", target, "-compiler", "gc", "-ldflags", "\""+ldFlags+"\"", "github.com/v2ray/v2ray-core/release/server")
	cmd.Env = append(cmd.Env, "GOOS="+string(vos), "GOARCH="+string(arch))
	cmd.Env = append(cmd.Env, os.Environ()...)
	_, err := cmd.Output()
	return err
}

func main() {
	flag.Parse()

	v2rayOS := getOS()
	v2rayArch := getArch()

	if err := os.Chdir(repoRoot); err != nil {
		fmt.Println("Unable to switch to V2Ray repo: " + err.Error())
		return
	}
	version, err := getRepoVersion()
	if err != nil {
		fmt.Println("Unable to detect V2Ray version: " + err.Error())
		return
	}
	fmt.Printf("Building V2Ray (%s) for %s %s\n", version, v2rayOS, v2rayArch)

	targetDir, err := createTargetDirectory(version, v2rayOS, v2rayArch)
	if err != nil {
		fmt.Println("Unable to create directory " + targetDir + ": " + err.Error())
	}

	targetFile := getTargetFile(v2rayOS, v2rayArch)
	err = buildV2Ray(targetDir, targetFile, version, v2rayOS, v2rayArch)
	if err != nil {
		fmt.Println("Unable to build V2Ray: " + err.Error())
	}

	err = copyConfigFiles(targetDir, version, v2rayOS, v2rayArch)
	if err != nil {
		fmt.Println("Unable to copy config files: " + err.Error())
	}
}

func init() {
	GOPATH = os.Getenv("GOPATH")
	repoRoot = filepath.Join(GOPATH, "src", "github.com", "v2ray", "v2ray-core")
}
