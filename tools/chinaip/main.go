package main

import (
	"bufio"
	"fmt"
	"math"
	"net"
	"net/http"
	"strconv"
	"strings"

	v2net "github.com/v2ray/v2ray-core/common/net"
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
	fmt.Println("map[uint32]byte {")
	for i := 0; i < len(dump); i += 2 {
		fmt.Println(dump[i], ": ", dump[i+1], ",")
	}
	fmt.Println("}")
}
