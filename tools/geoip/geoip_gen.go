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

	"v2ray.com/core/app/router"
	"v2ray.com/core/common/errors"
	"v2ray.com/core/tools/geoip"

	"github.com/golang/protobuf/proto"
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
		panic(errors.New("unexpected status ", resp.StatusCode))
	}
	defer resp.Body.Close()
	scanner := bufio.NewScanner(resp.Body)

	ips := &geoip.CountryIPRange{
		Ips: make([]*router.CIDR, 0, 8192),
	}
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
		mask := uint32(math.Floor(math.Log2(float64(count)) + 0.5))
		ipBytes := net.ParseIP(ip)
		if len(ipBytes) == 0 {
			panic("Invalid IP " + ip)
		}
		ips.Ips = append(ips.Ips, &router.CIDR{
			Ip:     []byte(ipBytes)[12:16],
			Prefix: 32 - mask,
		})
	}

	ipbytes, err := proto.Marshal(ips)
	if err != nil {
		log.Fatalf("Failed to marshal country IPs: %v", err)
	}

	file, err := os.OpenFile("geoip.generated.go", os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		log.Fatalf("Failed to generate geoip_data.go: %v", err)
	}
	defer file.Close()

	fmt.Fprintln(file, "package geoip")

	fmt.Fprintln(file, "var ChinaIPs = "+formatArray(ipbytes))
}

func formatArray(a []byte) string {
	r := "[]byte{"
	for idx, val := range a {
		if idx > 0 {
			r += ","
		}
		r += fmt.Sprintf("%d", val)
	}
	r += "}"
	return r
}
