package lexer

import (
	"testing"
)

func TestNextToken(t *testing.T) {
	values := []struct {
		css         string
		tokenTypes  []TokenType
		tokenValues []string
	}{
		{"color: red;", []TokenType{Ident, Colon, Whitespace, Ident, Semicolon}, []string{"color", ":", " ", "red", ";"}},
	}

	for _, v := range values {
		l := Lexer{Text: []rune(v.css)}
		for token, i := l.NextToken(), 0; token.Type != EOF; token, i = l.NextToken(), i+1 {
			actualTokenType := token.Type
			expectedTokenType := v.tokenTypes[i]
			actualTokenValue := string(token.Val)
			expectedTokenValue := v.tokenValues[i]

			if actualTokenValue != expectedTokenValue {
				t.Fatalf("%s %d. token value (expected) %v != %v (actual)", v.css, i, expectedTokenValue, actualTokenValue)
			}

			if actualTokenType != expectedTokenType {
				t.Fatalf("%s %d. token type (expected) %v != %v (actual)", v.css, i, expectedTokenType.String(), actualTokenType.String())
			}
		}
	}
}
