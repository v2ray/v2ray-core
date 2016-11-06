package conf

import (
	"errors"

	"v2ray.com/core/common/loader"
	"v2ray.com/core/transport/internet/authenticators/http"
	"v2ray.com/core/transport/internet/authenticators/noop"
	"v2ray.com/core/transport/internet/authenticators/srtp"
	"v2ray.com/core/transport/internet/authenticators/utp"
)

type NoOpAuthenticator struct{}

func (NoOpAuthenticator) Build() (*loader.TypedSettings, error) {
	return loader.NewTypedSettings(new(noop.Config)), nil
}

type NoOpConnectionAuthenticator struct{}

func (NoOpConnectionAuthenticator) Build() (*loader.TypedSettings, error) {
	return loader.NewTypedSettings(new(noop.Config)), nil
}

type SRTPAuthenticator struct{}

func (SRTPAuthenticator) Build() (*loader.TypedSettings, error) {
	return loader.NewTypedSettings(new(srtp.Config)), nil
}

type UTPAuthenticator struct{}

func (UTPAuthenticator) Build() (*loader.TypedSettings, error) {
	return loader.NewTypedSettings(new(utp.Config)), nil
}

type HTTPAuthenticatorRequest struct {
	Version *string                `json:"version"`
	Method  *string                `json:"method"`
	Path    *StringList            `json:"path"`
	Headers map[string]*StringList `json:"headers"`
}

func (this *HTTPAuthenticatorRequest) Build() (*http.RequestConfig, error) {
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

	if this.Version != nil {
		config.Version = &http.Version{Value: *this.Version}
	}

	if this.Method != nil {
		config.Method = &http.Method{Value: *this.Method}
	}

	if this.Path != nil && this.Path.Len() > 0 {
		config.Uri = append([]string(nil), (*this.Path)...)
	}

	if len(this.Headers) > 0 {
		config.Header = make([]*http.Header, 0, len(this.Headers))
		for key, value := range this.Headers {
			config.Header = append(config.Header, &http.Header{
				Name:  key,
				Value: append([]string(nil), (*value)...),
			})
		}
	}

	return config, nil
}

type HTTPAuthenticatorResponse struct {
	Version *string                `json:"version"`
	Status  *string                `json:"status"`
	Reason  *string                `json:"reason"`
	Headers map[string]*StringList `json:"headers"`
}

func (this *HTTPAuthenticatorResponse) Build() (*http.ResponseConfig, error) {
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
		},
	}

	if this.Version != nil {
		config.Version = &http.Version{Value: *this.Version}
	}

	if this.Status != nil || this.Reason != nil {
		config.Status = &http.Status{
			Code:   "200",
			Reason: "OK",
		}
		if this.Status != nil {
			config.Status.Code = *this.Status
		}
		if this.Reason != nil {
			config.Status.Reason = *this.Reason
		}
	}

	if len(this.Headers) > 0 {
		config.Header = make([]*http.Header, 0, len(this.Headers))
		for key, value := range this.Headers {
			config.Header = append(config.Header, &http.Header{
				Name:  key,
				Value: append([]string(nil), (*value)...),
			})
		}
	}

	return config, nil
}

type HTTPAuthenticator struct {
	Request  *HTTPAuthenticatorRequest  `json:"request"`
	Response *HTTPAuthenticatorResponse `json:"response"`
}

func (this *HTTPAuthenticator) Build() (*loader.TypedSettings, error) {
	config := new(http.Config)
	if this.Request == nil {
		return nil, errors.New("HTTP request settings not set.")
	}
	requestConfig, err := this.Request.Build()
	if err != nil {
		return nil, err
	}
	config.Request = requestConfig

	if this.Response != nil {
		responseConfig, err := this.Response.Build()
		if err != nil {
			return nil, err
		}
		config.Response = responseConfig
	}

	return loader.NewTypedSettings(config), nil
}
