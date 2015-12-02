package unit

import (
	"bytes"
	"fmt"
	"runtime"
	"strings"
	"testing"
)

var tGlobal *testing.T

func Current(t *testing.T) {
	tGlobal = t
}

func getCaller() (string, int) {
	stackLevel := 3
	for {
		_, file, line, ok := runtime.Caller(stackLevel)
		if strings.Contains(file, "assert") {
			stackLevel++
		} else {
			if ok {
				// Truncate file name at last file name separator.
				if index := strings.LastIndex(file, "/"); index >= 0 {
					file = file[index+1:]
				} else if index = strings.LastIndex(file, "\\"); index >= 0 {
					file = file[index+1:]
				}
			} else {
				file = "???"
				line = 1
			}
			return file, line
		}
	}
}

// decorate prefixes the string with the file and line of the call site
// and inserts the final newline if needed and indentation tabs for formatting.
func decorate(s string) string {
	file, line := getCaller()
	buf := new(bytes.Buffer)
	// Every line is indented at least one tab.
	buf.WriteString("  ")
	fmt.Fprintf(buf, "%s:%d: ", file, line)
	lines := strings.Split(s, "\n")
	if l := len(lines); l > 1 && lines[l-1] == "" {
		lines = lines[:l-1]
	}
	for i, line := range lines {
		if i > 0 {
			// Second and subsequent lines are indented an extra tab.
			buf.WriteString("\n\t\t")
		}
		buf.WriteString(line)
	}
	buf.WriteByte('\n')
	return buf.String()
}

func Fail(message string) {
	fmt.Println(decorate(message))
	tGlobal.Fail()
}
