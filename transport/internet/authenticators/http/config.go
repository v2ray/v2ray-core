package http

import (
	"v2ray.com/core/common/dice"
)

func (this *Version) GetValue() string {
	if this == nil {
		return "1.1"
	}
	return this.Value
}

func (this *Method) GetValue() string {
	if this == nil {
		return "GET"
	}
	return this.Value
}

func (this *Status) GetCode() string {
	if this == nil {
		return "200"
	}
	return this.Code
}

func (this *Status) GetReason() string {
	if this == nil {
		return "OK"
	}
	return this.Reason
}

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

func (this *RequestConfig) GetFullVersion() string {
	return "HTTP/" + this.Version.GetValue()
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

func (this *ResponseConfig) GetFullVersion() string {
	return "HTTP/" + this.Version.GetValue()
}
