package rsql

import (
	"unicode"
	"unicode/utf8"
)

var keywords = map[string]tokenType{
	"false": tokFalse,
	"true":  tokTrue,
	"null":  tokNull,
}

type lexer struct {
	src  string
	curr int
	next int
	ch   rune
}

func newLexer(source string) *lexer {
	l := &lexer{src: source}
	l.readRune()
	return l
}

func (l *lexer) readRune() {
	if l.next >= len(l.src) {
		l.ch = 0
		l.curr = l.next
		l.next = len(l.src)
		return
	}
	r, s := utf8.DecodeRuneInString(l.src[l.next:])
	l.ch = r
	l.curr = l.next
	l.next += s
}

func (l *lexer) peekRune() rune {
	if l.next >= len(l.src) {
		return 0
	}
	r, _ := utf8.DecodeRuneInString(l.src[l.next:])
	return r
}

func (l *lexer) Next() token {
	for unicode.IsSpace(l.ch) {
		l.readRune()
	}
	switch l.ch {
	case '=':
		if l.peekRune() == '=' {
			l.readRune() //=
			l.readRune() //=
			return token{tokEq, "=="}
		}
		l.readRune() //=
		return l.readFIQL()
	case '!':
		if l.peekRune() != '=' {
			l.readRune() //!
			l.readRune() //unexpected rune
			return token{tokInvalid, "!" + string(l.peekRune())}
		}
		l.readRune() //!
		l.readRune() //=
		return token{tokNe, "!="}
	case '<':
		l.readRune() //<
		if l.peekRune() == '=' {
			l.readRune() //=
			return token{tokLe, "<="}
		}
		return token{tokLt, "<"}
	case '>':
		l.readRune() //>
		if l.peekRune() == '=' {
			l.readRune() //=
			return token{tokGe, ">="}
		}
		return token{tokGt, ">"}
	case ';':
		l.readRune() //;
		return token{tokAnd, ";"}
	case ',':
		l.readRune() //,
		return token{tokOr, ","}
	case '(':
		l.readRune() //(
		return token{tokLparen, "("}
	case ')':
		l.readRune() //)
		return token{tokRparen, ")"}
	case 0:
		return token{tokEOF, ""}
	default:
		start := l.curr
		for l.ch != 0 && !isReserved(l.ch) {
			l.readRune()
		}
		lit := l.src[start:l.curr]
		switch {
		case isInteger(lit):
			return token{tokInteger, lit}
		case isFloat(lit):
			return token{tokFloat, lit}
		case isIdent(lit):
			return token{tokIdent, lit}
		default:
			if t, ok := keywords[lit]; ok {
				return token{t, lit}
			}
		}
		return token{tokString, lit}
	}
}

func isInteger(lit string) bool {
	if lit == "" {
		return false
	}
	runes := []rune(lit)
	i := 0
	if runes[0] == '+' || runes[0] == '-' {
		i = 1
		if len(runes) == 1 {
			return false
		}
	}
	hasDigit := false
	for ; i < len(runes); i++ {
		if !unicode.IsDigit(runes[i]) {
			return false
		}
		hasDigit = true
	}
	return hasDigit
}

func isFloat(lit string) bool {
	if lit == "" {
		return false
	}
	runes := []rune(lit)
	i := 0
	n := len(runes)

	if runes[0] == '+' || runes[0] == '-' {
		i++
		if i >= n {
			return false
		}
	}

	hasDot := false
	hasExp := false
	hasDigits := false

	for i < n && unicode.IsDigit(runes[i]) {
		i++
		hasDigits = true
	}

	if i < n && runes[i] == '.' {
		hasDot = true
		i++
		fracDigits := false
		for i < n && unicode.IsDigit(runes[i]) {
			i++
			fracDigits = true
		}
		if fracDigits {
			hasDigits = true
		}
	}

	if i < n && (runes[i] == 'e' || runes[i] == 'E') {
		hasExp = true
		i++
		if i < n && (runes[i] == '+' || runes[i] == '-') {
			i++
		}
		if i >= n || !unicode.IsDigit(runes[i]) {
			return false
		}
		for i < n && unicode.IsDigit(runes[i]) {
			i++
		}
	}

	return (hasDot || hasExp) && hasDigits && i == n
}

func isIdent(lit string) bool {
	if lit == "" {
		return false
	}
	runes := []rune(lit)
	if !unicode.IsLetter(runes[0]) && runes[0] != '_' {
		return false
	}
	for _, r := range runes[1:] {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' || r == '.' || r == '-' {
			continue
		}
		return false
	}
	return true
}

func isReserved(r rune) bool {
	switch r {
	case '=', '<', '>', '!', '(', ')', ';', ',':
		return true
	default:
		return false
	}
}

func (l *lexer) readFIQL() token {
	start := l.curr
	for unicode.IsLetter(l.ch) {
		l.readRune()
	}
	if l.ch != '=' {
		return token{tokInvalid, l.src[start:l.curr]}
	}
	switch l.src[start:l.curr] {
	case "in":
		l.readRune() //i
		l.readRune() //n
		l.readRune() //=
		return token{tokIn, "=in="}
	case "out":
		l.readRune() //o
		l.readRune() //u
		l.readRune() //t
		l.readRune() //=
		return token{tokOut, "=out="}
	case "lt":
		l.readRune() //l
		l.readRune() //t
		l.readRune() //=
		return token{tokLt, "=lt="}
	case "le":
		l.readRune() //l
		l.readRune() //e
		l.readRune() //=
		return token{tokLe, "=le="}
	case "gt":
		l.readRune() //g
		l.readRune() //t
		l.readRune() //=
		return token{tokGt, "=gt="}
	case "ge":
		l.readRune() //g
		l.readRune() //e
		l.readRune() //=
		return token{tokGe, "=ge="}
	default:
		return token{tokInvalid, l.src[start:l.curr]}
	}
}
