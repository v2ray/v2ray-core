// +build generate

package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"

	v2net "v2ray.com/core/common/net"
)

const (
	apnicFile = "http://ftp.apnic.net/apnic/stats/apnic/delegated-apnic-latest"
)

func main() {
	resp, err := http.Get(apnicFile)
	if err != nil {
		panic(err)
	}
	if resp.StatusCode != 200 {
		panic(fmt.Errorf("Unexpected status %d", resp.StatusCode))
	}
	defer resp.Body.Close()
	scanner := bufio.NewScanner(resp.Body)

	ipNet := v2net.NewIPNet()
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)
		parts := strings.Split(line, "|")
		if len(parts) < 5 {
			continue
		}
		if strings.ToLower(parts[1]) != "cn" || strings.ToLower(parts[2]) != "ipv4" {
			continue
		}
		ip := parts[3]
		count, err := strconv.Atoi(parts[4])
		if err != nil {
			continue
		}
		mask := 32 - int(math.Floor(math.Log2(float64(count))+0.5))
		cidr := fmt.Sprintf("%s/%d", ip, mask)
		_, t, err := net.ParseCIDR(cidr)
		if err != nil {
			panic(err)
		}
		ipNet.Add(t)
	}
	dump := ipNet.Serialize()

	file, err := os.OpenFile("chinaip_init.go", os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		log.Fatalf("Failed to generate chinaip_init.go: %v", err)
	}
	defer file.Close()

	fmt.Fprintln(file, "package rules")
	fmt.Fprintln(file, "import (")
	fmt.Fprintln(file, "v2net \"v2ray.com/core/common/net\"")
	fmt.Fprintln(file, ")")

	fmt.Fprintln(file, "var (")
	fmt.Fprintln(file, "chinaIPNet *v2net.IPNet")
	fmt.Fprintln(file, ")")

	fmt.Fprintln(file, "func init() {")

	fmt.Fprintln(file, "chinaIPNet = v2net.NewIPNetInitialValue(map[uint32]byte {")
	for i := 0; i < len(dump); i += 2 {
		fmt.Fprintln(file, dump[i], ": ", dump[i+1], ",")
	}
	fmt.Fprintln(file, "})")
	fmt.Fprintln(file, "}")
}
