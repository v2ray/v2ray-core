package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"

	v2net "github.com/v2ray/v2ray-core/common/net"
)

func main() {
	GOPATH := os.Getenv("GOPATH")
	src := filepath.Join(GOPATH, "src", "github.com", "v2ray", "v2ray-core", "tools", "chinaip", "ipv4.txt")
	reader, err := os.Open(src)
	if err != nil {
		panic(err)
	}
	ipNet := v2net.NewIPNet()

	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)
		if len(line) == 0 {
			break
		}
		_, t, err := net.ParseCIDR(line)
		if err != nil {
			panic(err)
		}
		ipNet.Add(t)
	}
	dump := ipNet.Serialize()
	fmt.Println("map[uint32]byte {")
	for i := 0; i < len(dump); i += 2 {
		fmt.Println(dump[i], ": ", dump[i+1], ",")
	}
	fmt.Println("}")
}
