package conf

import (
	"encoding/json"

	"errors"
	"strings"
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

type HTTPAuthenticatorHeader struct {
	Name  string      `json:"name"`
	Value *StringList `json:"value"`
}

type HTTPAuthenticatorRequest struct {
	Version *string                   `json:"version"`
	Method  *string                   `json:"method"`
	Path    *StringList               `json:"path"`
	Headers []HTTPAuthenticatorHeader `json:"headers"`
}

func (this *HTTPAuthenticatorRequest) Build() (*http.RequestConfig, error) {
	config := &http.RequestConfig{
		Uri: []string{"/"},
		Header: []*http.Header{
			{
				Name:  "Host",
				Value: []string{"www.baidu.com", "www.bing.com"},
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
		config.Header = make([]*http.Header, len(this.Headers))
		for idx, header := range this.Headers {
			config.Header[idx] = &http.Header{
				Name:  header.Name,
				Value: append([]string(nil), (*header.Value)...),
			}
		}
	}

	return config, nil
}

type HTTPAuthenticatorResponse struct {
	Version *string                   `json:"version"`
	Status  *string                   `json:"status"`
	Reason  *string                   `json:"reason"`
	Headers []HTTPAuthenticatorHeader `json:"headers"`
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
		config.Header = make([]*http.Header, len(this.Headers))
		for idx, header := range this.Headers {
			config.Header[idx] = &http.Header{
				Name:  header.Name,
				Value: append([]string(nil), (*header.Value)...),
			}
		}
	}

	return config, nil
}

type HTTPAuthenticator struct {
	Request  *HTTPAuthenticatorRequest `json:"request"`
	Response json.RawMessage           `json:"response"`
}

func (this *HTTPAuthenticator) Build() (*loader.TypedSettings, error) {
	config := new(http.Config)
	if this.Request != nil {
		requestConfig, err := this.Request.Build()
		if err != nil {
			return nil, err
		}
		config.Request = requestConfig
	}

	if len(this.Response) > 0 {
		var text string
		parsed := false
		if err := json.Unmarshal(this.Response, &text); err == nil {
			if strings.ToLower(text) != "disabled" {
				return nil, errors.New("Unknown HTTP header settings: " + text)
			}
			parsed = true
		}

		if !parsed {
			var response HTTPAuthenticatorResponse
			if err := json.Unmarshal(this.Response, &response); err != nil {
				return nil, errors.New("Failed to parse HTTP header response.")
			}
			responseConfig, err := response.Build()
			if err != nil {
				return nil, err
			}
			config.Response = responseConfig
		}
	}

	return loader.NewTypedSettings(config), nil
}
