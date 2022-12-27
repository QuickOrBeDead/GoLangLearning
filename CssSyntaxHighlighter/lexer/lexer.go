package lexer

import (
	"unicode"
	"unicode/utf8"
)

type Lexer struct {
	Text  []rune
	pos   int
	start int
}

type TokenType uint32

const (
	ErrorToken TokenType = iota
	IdentToken
	FunctionToken
	AtKeywordToken
	HashToken
	StringToken
	BadStringToken
	UrlToken
	NumberToken
	DimensionToken
	PercentageToken
	WhitespaceToken
	LeftParenthesisToken
	RightParenthesisToken
	LeftBraceToken
	RightBraceToken
	ColonToken
	SemicolonToken
	CommaToken
	CommentToken
	AtToken
	CDOToken
	CDCToken
	UnmatchedToken
	EOF
)

type Token struct {
	Type TokenType
	Val  []rune
}

func (t TokenType) String() string {
	switch t {
	case ErrorToken:
		return "Error"
	case IdentToken:
		return "Ident"
	case FunctionToken:
		return "Function"
	case AtKeywordToken:
		return "AtKeyword"
	case HashToken:
		return "Hash"
	case StringToken:
		return "String"
	case BadStringToken:
		return "BadString"
	case UrlToken:
		return "Url"
	case NumberToken:
		return "Number"
	case DimensionToken:
		return "Dimension"
	case PercentageToken:
		return "Percentage"
	case WhitespaceToken:
		return "Whitespace"
	case LeftParenthesisToken:
		return "LeftParenthesis"
	case RightParenthesisToken:
		return "RightParenthesis"
	case LeftBraceToken:
		return "LeftBrace"
	case RightBraceToken:
		return "RightBrace"
	case ColonToken:
		return "Colon"
	case SemicolonToken:
		return "Semicolon"
	case CommaToken:
		return "Comma"
	case CommentToken:
		return "Comment"
	case AtToken:
		return "At"
	case CDOToken:
		return "CDO"
	case CDCToken:
		return "CDC"
	case UnmatchedToken:
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

func (lex *Lexer) setPos(p int) {
	lex.pos = p
}

func (lex *Lexer) peek(c int) rune {
	pos := lex.pos + c
	if pos >= len(lex.Text) {
		return -1
	}
	return lex.Text[pos]
}

func (lex *Lexer) shift() []rune {
	r := lex.Text[lex.start:lex.pos:lex.pos]
	lex.start = lex.pos
	return r
}

func (lex *Lexer) NextToken() Token {
	var r rune
	switch r = lex.peek(0); {
	case r <= 0:
		return Token{Type: EOF, Val: []rune{}}
	case lex.scanWhitespace():
		return Token{Type: WhitespaceToken, Val: lex.shift()}
	case r == '"', r == '\'':
		t, v := lex.scanString(r)
		return Token{Type: t, Val: v}
	case r == '{':
		lex.next()
		return Token{Type: LeftBraceToken, Val: lex.shift()}
	case r == '}':
		lex.next()
		return Token{Type: RightBraceToken, Val: lex.shift()}
	case r == '(':
		lex.next()
		return Token{Type: LeftParenthesisToken, Val: lex.shift()}
	case r == ')':
		lex.next()
		return Token{Type: RightParenthesisToken, Val: lex.shift()}
	case r == ':':
		lex.next()
		return Token{Type: ColonToken, Val: lex.shift()}
	case r == ';':
		lex.next()
		return Token{Type: SemicolonToken, Val: lex.shift()}
	case r == ',':
		lex.next()
		return Token{Type: CommaToken, Val: lex.shift()}
	case r == '#':
		if isIdentStart(lex.peek(1)) {
			lex.next()
			_ = lex.scanIdent()
			return Token{Type: HashToken, Val: lex.shift()}
		}
	case r == '@':
		if isIdentStart(lex.peek(1)) {
			lex.next()
			_ = lex.scanIdent()
			return Token{Type: AtKeywordToken, Val: lex.shift()}
		}
	case r == '.':
		t := lex.scanNumericToken()
		if t != ErrorToken {
			return Token{Type: t, Val: lex.shift()}
		}
	case isIdentStart(r):
		_ = lex.scanIdent()
		val := lex.shift()
		if lex.peek(1) == '(' {
			if len(val) == 3 && matchASCIICaseInsensitive(val[0], 'u') && matchASCIICaseInsensitive(val[1], 'r') && matchASCIICaseInsensitive(val[2], 'l') {
				return Token{Type: UrlToken, Val: val}
			} else {
				return Token{Type: FunctionToken, Val: val}
			}
		} else {
			return Token{Type: IdentToken, Val: val}
		}
	case unicode.IsDigit(r):
		t := lex.scanNumericToken()
		if t != ErrorToken {
			return Token{Type: t, Val: lex.shift()}
		}
	}

	lex.next()
	return Token{Type: UnmatchedToken, Val: lex.shift()}
}

func (lex *Lexer) scanIdent() bool {
	r := false
	for isIdent(lex.peek(0)) {
		lex.next()
		r = true
	}

	return r
}

func (lex *Lexer) scanNumericToken() TokenType {
	if lex.scanNumber() {
		if lex.peek(0) == '%' {
			lex.next()
			return PercentageToken
		} else if lex.scanIdent() {
			return DimensionToken
		}

		return NumberToken
	}

	return ErrorToken
}

func (lex *Lexer) scanNumber() bool {
	startPos := lex.pos
	r := lex.peek(0)
	isFirstDigit := unicode.IsDigit(r)
	if isFirstDigit {
		lex.next()
	}

	for r = lex.peek(0); unicode.IsDigit(r); r = lex.peek(0) {
		lex.next()
	}

	if r == '.' {
		lex.next()

		for unicode.IsDigit(lex.peek(0)) {
			lex.next()
		}
	} else if !isFirstDigit {
		lex.setPos(startPos)
		return false
	}

	return true
}

func (lex *Lexer) scanWhitespace() bool {
	r := false
	for unicode.IsSpace(lex.peek(0)) {
		lex.next()
		r = true
	}

	return r
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
			return BadStringToken, lex.Text[startPos:lex.pos]
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

	return StringToken, lex.Text[startPos:lex.pos]
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
