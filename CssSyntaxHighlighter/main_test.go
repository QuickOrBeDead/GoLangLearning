package main

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
		l := Lexer{text: []rune(v.css)}
		for token, i := l.nextToken(), 0; token.Type != EOF; token, i = l.nextToken(), i+1 {
			actualTokenType := token.Type
			expectedTokenType := v.tokenTypes[i]
			actualTokenValue := string(token.Val)
			expectedTokenValue := v.tokenValues[i]

			if actualTokenValue != expectedTokenValue {
				t.Fatalf("%s %d. token value %v != %v", v.css, i, expectedTokenValue, actualTokenValue)
			}

			if actualTokenType != expectedTokenType {
				t.Fatalf("%s %d. token type %v != %v", v.css, i, expectedTokenType, actualTokenType)
			}
		}
	}
}
