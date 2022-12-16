package main

import (
	"html/template"
	"net/http"
	"os"
	"strings"
	"unicode"
	"unicode/utf8"
)

// https://www.w3.org/TR/css-syntax-3/#tokenizing-and-parsing
func main() {
	colors := make(map[TokenType]string)
	colors[Ident] = "blue"
	colors[Function] = "blue"
	colors[AtKeyword] = "blue"
	colors[Hash] = "red"
	colors[String] = "green"
	colors[BadString] = "green"
	colors[Url] = "blue"
	colors[Number] = "gray"
	colors[Dimension] = "blue"
	colors[Percentage] = "blue"
	colors[Whitespace] = "blue"
	colors[LeftParenthesis] = "blue"
	colors[RightParenthesis] = "blue"
	colors[LeftBrace] = "orange"
	colors[RightBrace] = "orange"
	colors[Colon] = "yellow"
	colors[Semicolon] = "yellow"
	colors[Comma] = "white"
	colors[Comment] = "blue"
	colors[At] = "blue"
	colors[CDO] = "blue"
	colors[CDC] = "blue"

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		page, err := loadFile("./Pages/Index.html")
		if page == nil || err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		t, err := template.New("t").
			Funcs(template.FuncMap{"html": func(s string) template.HTML { return template.HTML(s) }}).
			Parse(string(page))

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		data := make(map[string]interface{})

		if r.Method == http.MethodPost {
			var sb strings.Builder
			css := r.FormValue("cssText")
			for _, v := range getTokens(css) {
				if v.Type == Whitespace {
					sb.WriteString(string(v.Val))
				} else {
					color, ok := colors[v.Type]
					if !ok {
						color = "white"
					}

					sb.WriteString("<span style=\"color:")
					sb.WriteString(color)
					sb.WriteString("\">")
					sb.WriteString(string(v.Val))
					sb.WriteString("</span>")
				}
			}

			data["Result"] = sb.String()
			data["CssText"] = css
		}

		t.Execute(w, data)
	})
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("./Pages/css"))))
	http.ListenAndServe(":8080", nil)
}

func loadFile(path string) ([]byte, error) {
	if _, err := os.Stat(path); err == nil {
		return os.ReadFile(path)
	} else if err == os.ErrNotExist {
		return nil, nil
	} else {
		return nil, err
	}
}

type Lexer struct {
	text []rune
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

func (lex *Lexer) next() {
	if lex.pos >= len(lex.text) {
		return
	}

	lex.pos++
}

func (lex *Lexer) peek(c int) rune {
	pos := lex.pos + c
	if pos >= len(lex.text) {
		return -1
	}
	return lex.text[pos]
}

func (lex *Lexer) nextToken() Token {
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

	return lex.text[startPos:lex.pos]
}

func (lex *Lexer) scanNumber() []rune {
	startPos := lex.pos
	for unicode.IsDigit(lex.peek(0)) {
		lex.next()
	}

	return lex.text[startPos:lex.pos]
}

func (lex *Lexer) scanWhitespace() []rune {
	startPos := lex.pos
	for unicode.IsSpace(lex.peek(0)) {
		lex.next()
	}

	return lex.text[startPos:lex.pos]
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
			return BadString, lex.text[startPos:lex.pos]
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

	return String, lex.text[startPos:lex.pos]
}

func readRune(text string, pos int) (r rune, size int) {
	if r >= utf8.RuneSelf {
		r, size = utf8.DecodeRuneInString(text[pos:])
	} else {
		r, size = rune(text[pos]), 1
	}

	return r, size
}

func getTokens(text string) []Token {
	l := Lexer{text: []rune(text)}

	tokens := make([]Token, 0, 100)
	for t := l.nextToken(); t.Type != EOF; t = l.nextToken() {
		tokens = append(tokens, t)
	}

	return tokens
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
