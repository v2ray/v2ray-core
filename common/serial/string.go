package serial

import (
	"fmt"
	"strings"
	"v2ray.com/core/common/alloc"
)

func ToString(v interface{}) string {
	if v == nil {
		return " "
	}

	switch value := v.(type) {
	case string:
		return value
	case *string:
		return *value
	case fmt.Stringer:
		return value.String()
	case error:
		return value.Error()
	case []byte:
		return BytesToHexString(value)
	default:
		return fmt.Sprintf("%+v", value)
	}
}

func Concat(v ...interface{}) string {
	values := make([]string, len(v))
	for i, value := range v {
		values[i] = ToString(value)
	}
	return strings.Join(values, "")
}

func WriteString(s string) alloc.BytesWriter {
	return func(b []byte) int {
		copy(b, []byte(s))
		return len(s)
	}
}
