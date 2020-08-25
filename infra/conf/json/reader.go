package json

import (
	"io"

	"v2ray.com/core/common/buf"
)

// State is the internal state of parser.
type State byte

const (
	StateContent State = iota
	StateEscape
	StateDoubleQuote
	StateDoubleQuoteEscape
	StateSingleQuote
	StateSingleQuoteEscape
	StateComment
	StateSlash
	StateMultilineComment
	StateMultilineCommentStar
)

// Reader is a reader for filtering comments.
// It supports Java style single and multi line comment syntax, and Python style single line comment syntax.
type Reader struct {
	io.Reader

	state State
	br    *buf.BufferedReader
}

// Read implements io.Reader.Read(). Buffer must be at least 3 bytes.
func (v *Reader) Read(b []byte) (int, error) {
	if v.br == nil {
		v.br = &buf.BufferedReader{Reader: buf.NewReader(v.Reader)}
	}

	p := b[:0]
	for len(p) < len(b)-2 {
		x, err := v.br.ReadByte()
		if err != nil {
			if len(p) == 0 {
				return 0, err
			}
			return len(p), nil
		}
		switch v.state {
		case StateContent:
			switch x {
			case '"':
				v.state = StateDoubleQuote
				p = append(p, x)
			case '\'':
				v.state = StateSingleQuote
				p = append(p, x)
			case '\\':
				v.state = StateEscape
			case '#':
				v.state = StateComment
			case '/':
				v.state = StateSlash
			default:
				p = append(p, x)
			}
		case StateEscape:
			p = append(p, '\\', x)
			v.state = StateContent
		case StateDoubleQuote:
			switch x {
			case '"':
				v.state = StateContent
				p = append(p, x)
			case '\\':
				v.state = StateDoubleQuoteEscape
			default:
				p = append(p, x)
			}
		case StateDoubleQuoteEscape:
			p = append(p, '\\', x)
			v.state = StateDoubleQuote
		case StateSingleQuote:
			switch x {
			case '\'':
				v.state = StateContent
				p = append(p, x)
			case '\\':
				v.state = StateSingleQuoteEscape
			default:
				p = append(p, x)
			}
		case StateSingleQuoteEscape:
			p = append(p, '\\', x)
			v.state = StateSingleQuote
		case StateComment:
			if x == '\n' {
				v.state = StateContent
				p = append(p, '\n')
			}
		case StateSlash:
			switch x {
			case '/':
				v.state = StateComment
			case '*':
				v.state = StateMultilineComment
			default:
				p = append(p, '/', x)
			}
		case StateMultilineComment:
			switch x {
			case '*':
				v.state = StateMultilineCommentStar
			case '\n':
				p = append(p, '\n')
			}
		case StateMultilineCommentStar:
			switch x {
			case '/':
				v.state = StateContent
			case '*':
				// Stay
			case '\n':
				p = append(p, '\n')
			default:
				v.state = StateMultilineComment
			}
		default:
			panic("Unknown state.")
		}
	}
	return len(p), nil
}
