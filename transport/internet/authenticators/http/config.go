package http

import (
	"strings"
	"v2ray.com/core/common/dice"
)

func (v *Version) GetValue() string {
	if v == nil {
		return "1.1"
	}
	return v.Value
}

func (v *Method) GetValue() string {
	if v == nil {
		return "GET"
	}
	return v.Value
}

func (v *Status) GetCode() string {
	if v == nil {
		return "200"
	}
	return v.Code
}

func (v *Status) GetReason() string {
	if v == nil {
		return "OK"
	}
	return v.Reason
}

func pickString(arr []string) string {
	n := len(arr)
	switch n {
	case 0:
		return ""
	case 1:
		return arr[0]
	default:
		return arr[dice.Roll(n)]
	}
}

func (v *RequestConfig) PickUri() string {
	return pickString(v.Uri)
}

func (v *RequestConfig) PickHeaders() []string {
	n := len(v.Header)
	if n == 0 {
		return nil
	}
	headers := make([]string, n)
	for idx, headerConfig := range v.Header {
		headerName := headerConfig.Name
		headerValue := pickString(headerConfig.Value)
		headers[idx] = headerName + ": " + headerValue
	}
	return headers
}

func (v *RequestConfig) GetFullVersion() string {
	return "HTTP/" + v.Version.GetValue()
}

func (v *ResponseConfig) HasHeader(header string) bool {
	cHeader := strings.ToLower(header)
	for _, tHeader := range v.Header {
		if strings.ToLower(tHeader.Name) == cHeader {
			return true
		}
	}
	return false
}

func (v *ResponseConfig) PickHeaders() []string {
	n := len(v.Header)
	if n == 0 {
		return nil
	}
	headers := make([]string, n)
	for idx, headerConfig := range v.Header {
		headerName := headerConfig.Name
		headerValue := pickString(headerConfig.Value)
		headers[idx] = headerName + ": " + headerValue
	}
	return headers
}

func (v *ResponseConfig) GetFullVersion() string {
	return "HTTP/" + v.Version.GetValue()
}
