package lexer

import (
	"unicode"
	"unicode/utf8"
)

type Lexer struct {
	Text []rune
	pos  int
}

type TokenType uint32

const (
	Error TokenType = iota
	Ident
	Function
	AtKeyword
	Hash
	String
	BadString
	Url
	Number
	Dimension
	Percentage
	Whitespace
	LeftParenthesis
	RightParenthesis
	LeftBrace
	RightBrace
	Colon
	Semicolon
	Comma
	Comment
	At
	CDO
	CDC
	Unmatched
	EOF
)

type Token struct {
	Type TokenType
	Val  []rune
}

func (t TokenType) String() string {
	switch t {
	case Error:
		return "Error"
	case Ident:
		return "Ident"
	case Function:
		return "Function"
	case AtKeyword:
		return "AtKeyword"
	case Hash:
		return "Hash"
	case String:
		return "String"
	case BadString:
		return "BadString"
	case Url:
		return "Url"
	case Number:
		return "Number"
	case Dimension:
		return "Dimension"
	case Percentage:
		return "Percentage"
	case Whitespace:
		return "Whitespace"
	case LeftParenthesis:
		return "LeftParenthesis"
	case RightParenthesis:
		return "RightParenthesis"
	case LeftBrace:
		return "LeftBrace"
	case RightBrace:
		return "RightBrace"
	case Colon:
		return "Colon"
	case Semicolon:
		return "Semicolon"
	case Comma:
		return "Comma"
	case Comment:
		return "Comment"
	case At:
		return "At"
	case CDO:
		return "CDO"
	case CDC:
		return "CDC"
	case Unmatched:
		return "Unmatched"
	case EOF:
		return "EOF"
	default:
		return ""
	}
}

func (lex *Lexer) next() {
	if lex.pos >= len(lex.Text) {
		return
	}

	lex.pos++
}

func (lex *Lexer) peek(c int) rune {
	pos := lex.pos + c
	if pos >= len(lex.Text) {
		return -1
	}
	return lex.Text[pos]
}

func (lex *Lexer) NextToken() Token {
	var r rune
	switch r = lex.peek(0); {
	case r <= 0:
		return Token{Type: EOF, Val: []rune{}}
	case unicode.IsSpace(r):
		return Token{Type: Whitespace, Val: lex.scanWhitespace()}
	case r == '"', r == '\'':
		t, v := lex.scanString(r)
		return Token{Type: t, Val: v}
	case r == '{':
		lex.next()
		return Token{Type: LeftBrace, Val: []rune{r}}
	case r == '}':
		lex.next()
		return Token{Type: RightBrace, Val: []rune{r}}
	case r == '(':
		lex.next()
		return Token{Type: LeftParenthesis, Val: []rune{r}}
	case r == ')':
		lex.next()
		return Token{Type: RightParenthesis, Val: []rune{r}}
	case r == ':':
		lex.next()
		return Token{Type: Colon, Val: []rune{r}}
	case r == ';':
		lex.next()
		return Token{Type: Semicolon, Val: []rune{r}}
	case r == ',':
		lex.next()
		return Token{Type: Comma, Val: []rune{r}}
	case r == '#':
		if isIdentStart(lex.peek(1)) {
			lex.next()
			return Token{Type: Hash, Val: append([]rune{'#'}, lex.scanIdent()...)}
		}
	case r == '@':
		if isIdentStart(lex.peek(1)) {
			lex.next()
			return Token{Type: AtKeyword, Val: append([]rune{'@'}, lex.scanIdent()...)}
		}
	case isIdentStart(r):
		val := lex.scanIdent()
		if lex.peek(1) == '(' {
			if len(val) == 3 && matchASCIICaseInsensitive(val[0], 'u') && matchASCIICaseInsensitive(val[1], 'r') && matchASCIICaseInsensitive(val[2], 'l') {
				return Token{Type: Url, Val: val}
			} else {
				return Token{Type: Function, Val: val}
			}
		} else {
			return Token{Type: Ident, Val: val}
		}
	case unicode.IsDigit(r):
		return Token{Type: Number, Val: lex.scanNumber()}
	}

	lex.next()
	return Token{Type: Unmatched, Val: []rune{r}}
}

func (lex *Lexer) scanIdent() []rune {
	startPos := lex.pos
	for isIdent(lex.peek(0)) {
		lex.next()
	}

	return lex.Text[startPos:lex.pos]
}

func (lex *Lexer) scanNumber() []rune {
	startPos := lex.pos
	for unicode.IsDigit(lex.peek(0)) {
		lex.next()
	}

	return lex.Text[startPos:lex.pos]
}

func (lex *Lexer) scanWhitespace() []rune {
	startPos := lex.pos
	for unicode.IsSpace(lex.peek(0)) {
		lex.next()
	}

	return lex.Text[startPos:lex.pos]
}

func (lex *Lexer) scanEscapedChars() {

}

func (lex *Lexer) scanString(endRune rune) (TokenType, []rune) {
	startPos := lex.pos
	for {
		r := lex.peek(0)
		if r == endRune {
			lex.next()
			break
		} else if r <= 0 {
			break
		} else if isNewline(r) {
			lex.next()
			return BadString, lex.Text[startPos:lex.pos]
		} else if r == '\\' {
			r1 := lex.peek(0)
			if r1 <= 0 {
				break
			} else if isNewline(r1) {
				lex.next()
				continue
			} else {
				lex.scanEscapedChars()
			}
		} else {
			lex.next()
		}
	}

	return String, lex.Text[startPos:lex.pos]
}

func readRune(text string, pos int) (r rune, size int) {
	if r >= utf8.RuneSelf {
		r, size = utf8.DecodeRuneInString(text[pos:])
	} else {
		r, size = rune(text[pos]), 1
	}

	return r, size
}

func isIdent(r rune) bool {
	return isIdentStart(r) || isDigit(r) || r == '-'
}

func isNonASCII(r rune) bool {
	return r > unicode.MaxASCII
}

func matchASCIICaseInsensitive(r1 rune, r2 rune) bool {
	return r1 == r2 || (r1 >= 'A' && r1 <= 'Z' && r1+('a'-'A') == r2)
}

func isLetter(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z')
}

func isDigit(r rune) bool {
	return r >= '0' && r <= '9'
}

func isIdentStart(r rune) bool {
	return isLetter(r) || isNonASCII(r) || r == '_'
}

func isNewline(r rune) bool {
	return r == '\n' || r == '\r' || r == '\f'
}
