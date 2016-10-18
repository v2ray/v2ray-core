package conf

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"

	"strings"
	"v2ray.com/core/common/loader"
	v2net "v2ray.com/core/common/net"
	"v2ray.com/core/transport/internet"
	"v2ray.com/core/transport/internet/kcp"
	"v2ray.com/core/transport/internet/tcp"
	"v2ray.com/core/transport/internet/tls"
	"v2ray.com/core/transport/internet/ws"
)

var (
	kcpHeaderLoader = NewJSONConfigLoader(ConfigCreatorCache{
		"none": func() interface{} { return new(NoOpAuthenticator) },
		"srtp": func() interface{} { return new(SRTPAuthenticator) },
		"utp":  func() interface{} { return new(UTPAuthenticator) },
	}, "type", "")
)

type KCPConfig struct {
	Mtu             *uint32         `json:"mtu"`
	Tti             *uint32         `json:"tti"`
	UpCap           *uint32         `json:"uplinkCapacity"`
	DownCap         *uint32         `json:"downlinkCapacity"`
	Congestion      *bool           `json:"congestion"`
	ReadBufferSize  *uint32         `json:"readBufferSize"`
	WriteBufferSize *uint32         `json:"writeBufferSize"`
	HeaderConfig    json.RawMessage `json:"header"`
}

func (this *KCPConfig) Build() (*loader.TypedSettings, error) {
	config := new(kcp.Config)

	if this.Mtu != nil {
		mtu := *this.Mtu
		if mtu < 576 || mtu > 1460 {
			return nil, fmt.Errorf("KCP|Config: Invalid MTU size: %d", mtu)
		}
		config.Mtu = &kcp.MTU{Value: mtu}
	}
	if this.Tti != nil {
		tti := *this.Tti
		if tti < 10 || tti > 100 {
			return nil, fmt.Errorf("KCP|Config: Invalid TTI: %d", tti)
		}
		config.Tti = &kcp.TTI{Value: tti}
	}
	if this.UpCap != nil {
		config.UplinkCapacity = &kcp.UplinkCapacity{Value: *this.UpCap}
	}
	if this.DownCap != nil {
		config.DownlinkCapacity = &kcp.DownlinkCapacity{Value: *this.DownCap}
	}
	if this.Congestion != nil {
		config.Congestion = *this.Congestion
	}
	if this.ReadBufferSize != nil {
		size := *this.ReadBufferSize
		if size > 0 {
			config.ReadBuffer = &kcp.ReadBuffer{Size: size * 1024 * 1024}
		} else {
			config.ReadBuffer = &kcp.ReadBuffer{Size: 512 * 1024}
		}
	}
	if this.WriteBufferSize != nil {
		size := *this.WriteBufferSize
		if size > 0 {
			config.WriteBuffer = &kcp.WriteBuffer{Size: size * 1024 * 1024}
		} else {
			config.WriteBuffer = &kcp.WriteBuffer{Size: 512 * 1024}
		}
	}
	if len(this.HeaderConfig) > 0 {
		headerConfig, _, err := kcpHeaderLoader.Load(this.HeaderConfig)
		if err != nil {
			return nil, errors.New("KCP|Config: Failed to parse header config: " + err.Error())
		}
		ts, err := headerConfig.(Buildable).Build()
		if err != nil {
			return nil, errors.New("Failed to get KCP authenticator config: " + err.Error())
		}
		config.HeaderConfig = ts
	}

	return loader.NewTypedSettings(config), nil
}

type TCPConfig struct {
	ConnectionReuse *bool `json:"connectionReuse"`
}

func (this *TCPConfig) Build() (*loader.TypedSettings, error) {
	config := new(tcp.Config)
	if this.ConnectionReuse != nil {
		config.ConnectionReuse = &tcp.ConnectionReuse{
			Enable: *this.ConnectionReuse,
		}
	}
	return loader.NewTypedSettings(config), nil
}

type WebSocketConfig struct {
	ConnectionReuse *bool  `json:"connectionReuse"`
	Path            string `json:"Path"`
}

func (this *WebSocketConfig) Build() (*loader.TypedSettings, error) {
	config := &ws.Config{
		Path: this.Path,
	}
	if this.ConnectionReuse != nil {
		config.ConnectionReuse = &ws.ConnectionReuse{
			Enable: *this.ConnectionReuse,
		}
	}
	return loader.NewTypedSettings(config), nil
}

type TLSCertConfig struct {
	CertFile string `json:"certificateFile"`
	KeyFile  string `json:"keyFile"`
}
type TLSConfig struct {
	Insecure bool             `json:"allowInsecure"`
	Certs    []*TLSCertConfig `json:"certificates"`
}

func (this *TLSConfig) Build() (*loader.TypedSettings, error) {
	config := new(tls.Config)
	config.Certificate = make([]*tls.Certificate, len(this.Certs))
	for idx, certConf := range this.Certs {
		cert, err := ioutil.ReadFile(certConf.CertFile)
		if err != nil {
			return nil, errors.New("TLS: Failed to load certificate file: " + err.Error())
		}
		key, err := ioutil.ReadFile(certConf.KeyFile)
		if err != nil {
			return nil, errors.New("TLS: Failed to load key file: " + err.Error())
		}
		config.Certificate[idx] = &tls.Certificate{
			Key:         key,
			Certificate: cert,
		}
	}
	config.AllowInsecure = this.Insecure
	return loader.NewTypedSettings(config), nil
}

type StreamConfig struct {
	Network     *Network         `json:"network"`
	Security    string           `json:"security"`
	TLSSettings *TLSConfig       `json:"tlsSettings"`
	TCPSettings *TCPConfig       `json:"tcpSettings"`
	KCPSettings *KCPConfig       `json:"kcpSettings"`
	WSSettings  *WebSocketConfig `json:"wsSettings"`
}

func (this *StreamConfig) Build() (*internet.StreamConfig, error) {
	config := &internet.StreamConfig{
		Network: v2net.Network_RawTCP,
	}
	if this.Network != nil {
		config.Network = (*this.Network).Build()
	}
	if strings.ToLower(this.Security) == "tls" {
		tlsSettings := this.TLSSettings
		if tlsSettings == nil {
			tlsSettings = &TLSConfig{}
		}
		ts, err := tlsSettings.Build()
		if err != nil {
			return nil, errors.New("Failed to build TLS config: " + err.Error())
		}
		config.SecuritySettings = append(config.SecuritySettings, ts)
	}
	if this.TCPSettings != nil {
		ts, err := this.TCPSettings.Build()
		if err != nil {
			return nil, errors.New("Failed to build TCP config: " + err.Error())
		}
		config.NetworkSettings = append(config.NetworkSettings, &internet.NetworkSettings{
			Network:  v2net.Network_TCP,
			Settings: ts,
		})
	}
	if this.KCPSettings != nil {
		ts, err := this.KCPSettings.Build()
		if err != nil {
			return nil, errors.New("Failed to build KCP config: " + err.Error())
		}
		config.NetworkSettings = append(config.NetworkSettings, &internet.NetworkSettings{
			Network:  v2net.Network_KCP,
			Settings: ts,
		})
	}
	if this.WSSettings != nil {
		ts, err := this.WSSettings.Build()
		if err != nil {
			return nil, errors.New("Failed to build WebSocket config: " + err.Error())
		}
		config.NetworkSettings = append(config.NetworkSettings, &internet.NetworkSettings{
			Network:  v2net.Network_WebSocket,
			Settings: ts,
		})
	}
	return config, nil
}
