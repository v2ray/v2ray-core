package unit

import (
	"bytes"
	"fmt"
	"runtime"
	"strings"
)

type Subject struct {
	assert *Assertion
	name   string
}

func NewSubject(assert *Assertion) *Subject {
	return &Subject{
		assert: assert,
		name:   "",
	}
}

// decorate prefixes the string with the file and line of the call site
// and inserts the final newline if needed and indentation tabs for formatting.
func decorate(s string) string {
	_, file, line, ok := runtime.Caller(3)
	if strings.Contains(file, "testing") {
		_, file, line, ok = runtime.Caller(4)
	}
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

func (subject *Subject) FailWithMessage(message string) {
	fmt.Println(decorate(message))
	subject.assert.t.Fail()
}

func (subject *Subject) Named(name string) {
	subject.name = name
}

func (subject *Subject) DisplayString(value string) string {
	if len(value) == 0 {
		value = "unknown"
	}
	if len(subject.name) == 0 {
		return "<" + value + ">"
	}
	return subject.name + "(<" + value + ">)"
}
