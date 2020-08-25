package conf

import (
	"sort"

	"github.com/golang/protobuf/proto"

	"v2ray.com/core/transport/internet/headers/http"
	"v2ray.com/core/transport/internet/headers/noop"
	"v2ray.com/core/transport/internet/headers/srtp"
	"v2ray.com/core/transport/internet/headers/tls"
	"v2ray.com/core/transport/internet/headers/utp"
	"v2ray.com/core/transport/internet/headers/wechat"
	"v2ray.com/core/transport/internet/headers/wireguard"
)

type NoOpAuthenticator struct{}

func (NoOpAuthenticator) Build() (proto.Message, error) {
	return new(noop.Config), nil
}

type NoOpConnectionAuthenticator struct{}

func (NoOpConnectionAuthenticator) Build() (proto.Message, error) {
	return new(noop.ConnectionConfig), nil
}

type SRTPAuthenticator struct{}

func (SRTPAuthenticator) Build() (proto.Message, error) {
	return new(srtp.Config), nil
}

type UTPAuthenticator struct{}

func (UTPAuthenticator) Build() (proto.Message, error) {
	return new(utp.Config), nil
}

type WechatVideoAuthenticator struct{}

func (WechatVideoAuthenticator) Build() (proto.Message, error) {
	return new(wechat.VideoConfig), nil
}

type WireguardAuthenticator struct{}

func (WireguardAuthenticator) Build() (proto.Message, error) {
	return new(wireguard.WireguardConfig), nil
}

type DTLSAuthenticator struct{}

func (DTLSAuthenticator) Build() (proto.Message, error) {
	return new(tls.PacketConfig), nil
}

type HTTPAuthenticatorRequest struct {
	Version string                 `json:"version"`
	Method  string                 `json:"method"`
	Path    StringList             `json:"path"`
	Headers map[string]*StringList `json:"headers"`
}

func sortMapKeys(m map[string]*StringList) []string {
	var keys []string
	for key := range m {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

func (v *HTTPAuthenticatorRequest) Build() (*http.RequestConfig, error) {
	config := &http.RequestConfig{
		Uri: []string{"/"},
		Header: []*http.Header{
			{
				Name:  "Host",
				Value: []string{"www.baidu.com", "www.bing.com"},
			},
			{
				Name: "User-Agent",
				Value: []string{
					"Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/53.0.2785.143 Safari/537.36",
					"Mozilla/5.0 (iPhone; CPU iPhone OS 10_0_2 like Mac OS X) AppleWebKit/601.1 (KHTML, like Gecko) CriOS/53.0.2785.109 Mobile/14A456 Safari/601.1.46",
				},
			},
			{
				Name:  "Accept-Encoding",
				Value: []string{"gzip, deflate"},
			},
			{
				Name:  "Connection",
				Value: []string{"keep-alive"},
			},
			{
				Name:  "Pragma",
				Value: []string{"no-cache"},
			},
		},
	}

	if len(v.Version) > 0 {
		config.Version = &http.Version{Value: v.Version}
	}

	if len(v.Method) > 0 {
		config.Method = &http.Method{Value: v.Method}
	}

	if len(v.Path) > 0 {
		config.Uri = append([]string(nil), (v.Path)...)
	}

	if len(v.Headers) > 0 {
		config.Header = make([]*http.Header, 0, len(v.Headers))
		headerNames := sortMapKeys(v.Headers)
		for _, key := range headerNames {
			value := v.Headers[key]
			if value == nil {
				return nil, newError("empty HTTP header value: " + key).AtError()
			}
			config.Header = append(config.Header, &http.Header{
				Name:  key,
				Value: append([]string(nil), (*value)...),
			})
		}
	}

	return config, nil
}

type HTTPAuthenticatorResponse struct {
	Version string                 `json:"version"`
	Status  string                 `json:"status"`
	Reason  string                 `json:"reason"`
	Headers map[string]*StringList `json:"headers"`
}

func (v *HTTPAuthenticatorResponse) Build() (*http.ResponseConfig, error) {
	config := &http.ResponseConfig{
		Header: []*http.Header{
			{
				Name:  "Content-Type",
				Value: []string{"application/octet-stream", "video/mpeg"},
			},
			{
				Name:  "Transfer-Encoding",
				Value: []string{"chunked"},
			},
			{
				Name:  "Connection",
				Value: []string{"keep-alive"},
			},
			{
				Name:  "Pragma",
				Value: []string{"no-cache"},
			},
			{
				Name:  "Cache-Control",
				Value: []string{"private", "no-cache"},
			},
		},
	}

	if len(v.Version) > 0 {
		config.Version = &http.Version{Value: v.Version}
	}

	if len(v.Status) > 0 || len(v.Reason) > 0 {
		config.Status = &http.Status{
			Code:   "200",
			Reason: "OK",
		}
		if len(v.Status) > 0 {
			config.Status.Code = v.Status
		}
		if len(v.Reason) > 0 {
			config.Status.Reason = v.Reason
		}
	}

	if len(v.Headers) > 0 {
		config.Header = make([]*http.Header, 0, len(v.Headers))
		headerNames := sortMapKeys(v.Headers)
		for _, key := range headerNames {
			value := v.Headers[key]
			if value == nil {
				return nil, newError("empty HTTP header value: " + key).AtError()
			}
			config.Header = append(config.Header, &http.Header{
				Name:  key,
				Value: append([]string(nil), (*value)...),
			})
		}
	}

	return config, nil
}

type HTTPAuthenticator struct {
	Request  HTTPAuthenticatorRequest  `json:"request"`
	Response HTTPAuthenticatorResponse `json:"response"`
}

func (v *HTTPAuthenticator) Build() (proto.Message, error) {
	config := new(http.Config)
	requestConfig, err := v.Request.Build()
	if err != nil {
		return nil, err
	}
	config.Request = requestConfig

	responseConfig, err := v.Response.Build()
	if err != nil {
		return nil, err
	}
	config.Response = responseConfig

	return config, nil
}
