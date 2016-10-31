package http

import (
	"v2ray.com/core/common/dice"
)

func pickString(arr []string) string {
	n := len(arr)
	if n == 0 {
		return ""
	}
	return arr[dice.Roll(n)]
}

func (this *RequestConfig) PickUri() string {
	return pickString(this.Uri)
}

func (this *RequestConfig) PickHeaders() []string {
	n := len(this.Header)
	if n == 0 {
		return nil
	}
	headers := make([]string, n)
	for idx, headerConfig := range this.Header {
		headerName := headerConfig.Name
		headerValue := pickString(headerConfig.Value)
		headers[idx] = headerName + ": " + headerValue
	}
	return headers
}

func (this *RequestConfig) GetVersion() string {
	return "HTTP/" + this.Version
}

func (this *ResponseConfig) PickHeaders() []string {
	n := len(this.Header)
	if n == 0 {
		return nil
	}
	headers := make([]string, n)
	for idx, headerConfig := range this.Header {
		headerName := headerConfig.Name
		headerValue := pickString(headerConfig.Value)
		headers[idx] = headerName + ": " + headerValue
	}
	return headers
}

func (this *ResponseConfig) GetVersion() string {
	return "HTTP/" + this.Version
}
