package json

import (
	"io"
)

type State byte

const (
	StateContent              State = 0
	StateEscape               State = 1
	StateDoubleQuote          State = 2
	StateSingleQuote          State = 3
	StateComment              State = 4
	StateSlash                State = 5
	StateMultilineComment     State = 6
	StateMultilineCommentStar State = 7
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
			p = append(p, x)
			v.state = StateContent
		case StateDoubleQuote:
			if x == '"' {
				v.state = StateContent
			}
			p = append(p, x)
		case StateSingleQuote:
			if x == '\'' {
				v.state = StateContent
			}
			p = append(p, x)
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
