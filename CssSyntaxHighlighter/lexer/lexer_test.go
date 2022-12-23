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
		{" ", []TokenType{WhitespaceToken}, []string{" "}},
		{"5.2", []TokenType{NumberToken}, []string{"5.2"}},
		{"50.12", []TokenType{NumberToken}, []string{"50.12"}},
		{"5.103 1.02", []TokenType{NumberToken, WhitespaceToken, NumberToken}, []string{"5.103", " ", "1.02"}},
		{"color: red;", []TokenType{IdentToken, ColonToken, WhitespaceToken, IdentToken, SemicolonToken}, []string{"color", ":", " ", "red", ";"}},
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
