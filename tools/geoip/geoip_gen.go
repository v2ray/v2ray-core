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
)

const (
	apnicFile = "http://ftp.apnic.net/apnic/stats/apnic/delegated-apnic-latest"
)

type IPEntry struct {
	IP   []byte
	Bits uint32
}

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

	ips := make([]IPEntry, 0, 8192)
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
		ips = append(ips, IPEntry{
			IP:   []byte(ipBytes),
			Bits: mask,
		})
	}

	file, err := os.OpenFile("geoip_data.go", os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		log.Fatalf("Failed to generate geoip_data.go: %v", err)
	}
	defer file.Close()

	fmt.Fprintln(file, "package geoip")
	fmt.Fprintln(file, "import \"v2ray.com/core/app/router\"")

	fmt.Fprintln(file, "var ChinaIPs []*router.IP")

	fmt.Fprintln(file, "func init() {")

	fmt.Fprintln(file, "ChinaIPs = []*router.IP {")
	for _, ip := range ips {
		fmt.Fprintln(file, "&router.IP{", formatArray(ip.IP[12:16]), ",", ip.Bits, "},")
	}
	fmt.Fprintln(file, "}")
	fmt.Fprintln(file, "}")
}

func formatArray(a []byte) string {
	r := "[]byte{"
	for idx, v := range a {
		if idx > 0 {
			r += ","
		}
		r += fmt.Sprintf("%d", v)
	}
	r += "}"
	return r
}
