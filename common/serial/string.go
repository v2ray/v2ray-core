package serial

import (
	"fmt"
	"strings"
)

// ToString serialize an abitrary value into string.
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
	default:
		return fmt.Sprintf("%+v", value)
	}
}

func Concat(v ...interface{}) string {
	builder := strings.Builder{}
	for _, value := range v {
		builder.WriteString(ToString(value))
	}
	return builder.String()
}

func WriteString(s string) func([]byte) (int, error) {
	return func(b []byte) (int, error) {
		return copy(b, []byte(s)), nil
	}
}
