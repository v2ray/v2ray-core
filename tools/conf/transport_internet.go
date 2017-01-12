package conf

import (
	"encoding/json"
	"io/ioutil"
	"strings"

	"v2ray.com/core/common/errors"
	"v2ray.com/core/common/serial"
	"v2ray.com/core/transport/internet"
	"v2ray.com/core/transport/internet/kcp"
	"v2ray.com/core/transport/internet/tcp"
	"v2ray.com/core/transport/internet/tls"
	"v2ray.com/core/transport/internet/websocket"
)

var (
	kcpHeaderLoader = NewJSONConfigLoader(ConfigCreatorCache{
		"none":         func() interface{} { return new(NoOpAuthenticator) },
		"srtp":         func() interface{} { return new(SRTPAuthenticator) },
		"utp":          func() interface{} { return new(UTPAuthenticator) },
		"wechat-video": func() interface{} { return new(WechatVideoAuthenticator) },
	}, "type", "")

	tcpHeaderLoader = NewJSONConfigLoader(ConfigCreatorCache{
		"none": func() interface{} { return new(NoOpConnectionAuthenticator) },
		"http": func() interface{} { return new(HTTPAuthenticator) },
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

func (v *KCPConfig) Build() (*serial.TypedMessage, error) {
	config := new(kcp.Config)

	if v.Mtu != nil {
		mtu := *v.Mtu
		if mtu < 576 || mtu > 1460 {
			return nil, errors.Format("KCP|Config: Invalid MTU size: %d", mtu)
		}
		config.Mtu = &kcp.MTU{Value: mtu}
	}
	if v.Tti != nil {
		tti := *v.Tti
		if tti < 10 || tti > 100 {
			return nil, errors.Format("KCP|Config: Invalid TTI: %d", tti)
		}
		config.Tti = &kcp.TTI{Value: tti}
	}
	if v.UpCap != nil {
		config.UplinkCapacity = &kcp.UplinkCapacity{Value: *v.UpCap}
	}
	if v.DownCap != nil {
		config.DownlinkCapacity = &kcp.DownlinkCapacity{Value: *v.DownCap}
	}
	if v.Congestion != nil {
		config.Congestion = *v.Congestion
	}
	if v.ReadBufferSize != nil {
		size := *v.ReadBufferSize
		if size > 0 {
			config.ReadBuffer = &kcp.ReadBuffer{Size: size * 1024 * 1024}
		} else {
			config.ReadBuffer = &kcp.ReadBuffer{Size: 512 * 1024}
		}
	}
	if v.WriteBufferSize != nil {
		size := *v.WriteBufferSize
		if size > 0 {
			config.WriteBuffer = &kcp.WriteBuffer{Size: size * 1024 * 1024}
		} else {
			config.WriteBuffer = &kcp.WriteBuffer{Size: 512 * 1024}
		}
	}
	if len(v.HeaderConfig) > 0 {
		headerConfig, _, err := kcpHeaderLoader.Load(v.HeaderConfig)
		if err != nil {
			return nil, errors.Base(err).Message("Invalid mKCP header config.")
		}
		ts, err := headerConfig.(Buildable).Build()
		if err != nil {
			return nil, errors.Base(err).Message("Invalid mKCP header config.")
		}
		config.HeaderConfig = ts
	}

	return serial.ToTypedMessage(config), nil
}

type TCPConfig struct {
	ConnectionReuse *bool           `json:"connectionReuse"`
	HeaderConfig    json.RawMessage `json:"header"`
}

func (v *TCPConfig) Build() (*serial.TypedMessage, error) {
	config := new(tcp.Config)
	if v.ConnectionReuse != nil {
		config.ConnectionReuse = &tcp.ConnectionReuse{
			Enable: *v.ConnectionReuse,
		}
	}
	if len(v.HeaderConfig) > 0 {
		headerConfig, _, err := tcpHeaderLoader.Load(v.HeaderConfig)
		if err != nil {
			return nil, errors.Base(err).Message("Invalid TCP header config.")
		}
		ts, err := headerConfig.(Buildable).Build()
		if err != nil {
			return nil, errors.Base(err).Message("Invalid TCP header config.")
		}
		config.HeaderSettings = ts
	}

	return serial.ToTypedMessage(config), nil
}

type WebSocketConfig struct {
	ConnectionReuse *bool  `json:"connectionReuse"`
	Path            string `json:"Path"`
}

func (v *WebSocketConfig) Build() (*serial.TypedMessage, error) {
	config := &websocket.Config{
		Path: v.Path,
	}
	if v.ConnectionReuse != nil {
		config.ConnectionReuse = &websocket.ConnectionReuse{
			Enable: *v.ConnectionReuse,
		}
	}
	return serial.ToTypedMessage(config), nil
}

type TLSCertConfig struct {
	CertFile string `json:"certificateFile"`
	KeyFile  string `json:"keyFile"`
}
type TLSConfig struct {
	Insecure   bool             `json:"allowInsecure"`
	Certs      []*TLSCertConfig `json:"certificates"`
	ServerName string           `json:"serverName"`
}

func (v *TLSConfig) Build() (*serial.TypedMessage, error) {
	config := new(tls.Config)
	config.Certificate = make([]*tls.Certificate, len(v.Certs))
	for idx, certConf := range v.Certs {
		cert, err := ioutil.ReadFile(certConf.CertFile)
		if err != nil {
			return nil, errors.Base(err).Message("Failed to load TLS certificate file: ", certConf.CertFile)
		}
		key, err := ioutil.ReadFile(certConf.KeyFile)
		if err != nil {
			return nil, errors.Base(err).Message("Failed to load TLS key file: ", certConf.KeyFile)
		}
		config.Certificate[idx] = &tls.Certificate{
			Key:         key,
			Certificate: cert,
		}
	}
	config.AllowInsecure = v.Insecure
	if len(v.ServerName) > 0 {
		config.ServerName = v.ServerName
	}
	return serial.ToTypedMessage(config), nil
}

type TransportProtocol string

func (p TransportProtocol) Build() (internet.TransportProtocol, error) {
	switch strings.ToLower(string(p)) {
	case "tcp":
		return internet.TransportProtocol_TCP, nil
	case "kcp", "mkcp":
		return internet.TransportProtocol_MKCP, nil
	case "ws", "websocket":
		return internet.TransportProtocol_WebSocket, nil
	default:
		return internet.TransportProtocol_TCP, errors.New("Config: unknown transport protocol: ", p)
	}
}

type StreamConfig struct {
	Network     *TransportProtocol `json:"network"`
	Security    string             `json:"security"`
	TLSSettings *TLSConfig         `json:"tlsSettings"`
	TCPSettings *TCPConfig         `json:"tcpSettings"`
	KCPSettings *KCPConfig         `json:"kcpSettings"`
	WSSettings  *WebSocketConfig   `json:"wsSettings"`
}

func (v *StreamConfig) Build() (*internet.StreamConfig, error) {
	config := &internet.StreamConfig{
		Protocol: internet.TransportProtocol_TCP,
	}
	if v.Network != nil {
		protocol, err := (*v.Network).Build()
		if err != nil {
			return nil, err
		}
		config.Protocol = protocol
	}
	if strings.ToLower(v.Security) == "tls" {
		tlsSettings := v.TLSSettings
		if tlsSettings == nil {
			tlsSettings = &TLSConfig{}
		}
		ts, err := tlsSettings.Build()
		if err != nil {
			return nil, errors.Base(err).Message("Failed to build TLS config.")
		}
		config.SecuritySettings = append(config.SecuritySettings, ts)
		config.SecurityType = ts.Type
	}
	if v.TCPSettings != nil {
		ts, err := v.TCPSettings.Build()
		if err != nil {
			return nil, errors.Base(err).Message("Failed to build TCP config.")
		}
		config.TransportSettings = append(config.TransportSettings, &internet.TransportSettings{
			Protocol: internet.TransportProtocol_TCP,
			Settings: ts,
		})
	}
	if v.KCPSettings != nil {
		ts, err := v.KCPSettings.Build()
		if err != nil {
			return nil, errors.Base(err).Message("Failed to build mKCP config.")
		}
		config.TransportSettings = append(config.TransportSettings, &internet.TransportSettings{
			Protocol: internet.TransportProtocol_MKCP,
			Settings: ts,
		})
	}
	if v.WSSettings != nil {
		ts, err := v.WSSettings.Build()
		if err != nil {
			return nil, errors.Base(err).Message("Failed to build WebSocket config.")
		}
		config.TransportSettings = append(config.TransportSettings, &internet.TransportSettings{
			Protocol: internet.TransportProtocol_WebSocket,
			Settings: ts,
		})
	}
	return config, nil
}

type ProxyConfig struct {
	Tag string `json:"tag"`
}

func (v *ProxyConfig) Build() (*internet.ProxyConfig, error) {
	if len(v.Tag) == 0 {
		return nil, errors.New("Proxy tag is not set.")
	}
	return &internet.ProxyConfig{
		Tag: v.Tag,
	}, nil
}
