package serial

import (
	"fmt"
	"strings"
)

// ToString serialize an arbitrary value into string.
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

// Concat concatenates all input into a single string.
func Concat(v ...interface{}) string {
	builder := strings.Builder{}
	for _, value := range v {
		builder.WriteString(ToString(value))
	}
	return builder.String()
}
