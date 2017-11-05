package net_test

import (
	"net"
	"os"
	"path/filepath"
	"testing"

	proto "github.com/golang/protobuf/proto"
	"v2ray.com/core/app/router"
	"v2ray.com/core/common/platform"

	"v2ray.com/ext/sysio"

	"v2ray.com/core/common"
	. "v2ray.com/core/common/net"
	. "v2ray.com/ext/assert"
)

func parseCIDR(str string) *net.IPNet {
	_, ipNet, err := net.ParseCIDR(str)
	common.Must(err)
	return ipNet
}

func TestIPNet(t *testing.T) {
	assert := With(t)

	ipNet := NewIPNetTable()
	ipNet.Add(parseCIDR(("0.0.0.0/8")))
	ipNet.Add(parseCIDR(("10.0.0.0/8")))
	ipNet.Add(parseCIDR(("100.64.0.0/10")))
	ipNet.Add(parseCIDR(("127.0.0.0/8")))
	ipNet.Add(parseCIDR(("169.254.0.0/16")))
	ipNet.Add(parseCIDR(("172.16.0.0/12")))
	ipNet.Add(parseCIDR(("192.0.0.0/24")))
	ipNet.Add(parseCIDR(("192.0.2.0/24")))
	ipNet.Add(parseCIDR(("192.168.0.0/16")))
	ipNet.Add(parseCIDR(("198.18.0.0/15")))
	ipNet.Add(parseCIDR(("198.51.100.0/24")))
	ipNet.Add(parseCIDR(("203.0.113.0/24")))
	ipNet.Add(parseCIDR(("8.8.8.8/32")))
	ipNet.AddIP(net.ParseIP("91.108.4.0"), 16)
	assert(ipNet.Contains(ParseIP("192.168.1.1")), IsTrue)
	assert(ipNet.Contains(ParseIP("192.0.0.0")), IsTrue)
	assert(ipNet.Contains(ParseIP("192.0.1.0")), IsFalse)
	assert(ipNet.Contains(ParseIP("0.1.0.0")), IsTrue)
	assert(ipNet.Contains(ParseIP("1.0.0.1")), IsFalse)
	assert(ipNet.Contains(ParseIP("8.8.8.7")), IsFalse)
	assert(ipNet.Contains(ParseIP("8.8.8.8")), IsTrue)
	assert(ipNet.Contains(ParseIP("2001:cdba::3257:9652")), IsFalse)
	assert(ipNet.Contains(ParseIP("91.108.255.254")), IsTrue)
}

func loadGeoIP(country string) ([]*router.CIDR, error) {
	geoipBytes, err := sysio.ReadAsset("geoip.dat")
	if err != nil {
		return nil, err
	}
	var geoipList router.GeoIPList
	if err := proto.Unmarshal(geoipBytes, &geoipList); err != nil {
		return nil, err
	}

	for _, geoip := range geoipList.Entry {
		if geoip.CountryCode == country {
			return geoip.Cidr, nil
		}
	}

	panic("country not found: " + country)
}

func BenchmarkIPNetQuery(b *testing.B) {
	common.Must(sysio.CopyFile(platform.GetAssetLocation("geoip.dat"), filepath.Join(os.Getenv("GOPATH"), "src", "v2ray.com", "core", "tools", "release", "config", "geoip.dat")))

	ips, err := loadGeoIP("CN")
	common.Must(err)

	ipNet := NewIPNetTable()
	for _, ip := range ips {
		ipNet.AddIP(ip.Ip, byte(ip.Prefix))
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		ipNet.Contains([]byte{8, 8, 8, 8})
	}
}

func BenchmarkCIDRQuery(b *testing.B) {
	common.Must(sysio.CopyFile(platform.GetAssetLocation("geoip.dat"), filepath.Join(os.Getenv("GOPATH"), "src", "v2ray.com", "core", "tools", "release", "config", "geoip.dat")))

	ips, err := loadGeoIP("CN")
	common.Must(err)

	ipNet := make([]*net.IPNet, 0, 1024)
	for _, ip := range ips {
		if len(ip.Ip) != 4 {
			continue
		}
		ipNet = append(ipNet, &net.IPNet{
			IP:   net.IP(ip.Ip),
			Mask: net.CIDRMask(int(ip.Prefix), 32),
		})
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		for _, n := range ipNet {
			if n.Contains([]byte{8, 8, 8, 8}) {
				break
			}
		}
	}
}
