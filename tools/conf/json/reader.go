package json

import (
	"io"
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

type Reader struct {
	io.Reader
	state State
}

func (v *Reader) Read(b []byte) (int, error) {
	n, err := v.Reader.Read(b)
	if err != nil {
		return n, err
	}
	p := b[:0]
	for _, x := range b[:n] {
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
			if x == '*' {
				v.state = StateMultilineCommentStar
			}
		case StateMultilineCommentStar:
			switch x {
			case '/':
				v.state = StateContent
			case '*':
				// Stay
			default:
				v.state = StateMultilineComment
			}
		default:
			panic("Unknown state.")
		}
	}
	return len(p), nil
}
