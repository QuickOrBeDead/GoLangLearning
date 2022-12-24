package lexer

import (
	"testing"
)

type TestDataToken struct {
	tokenType TokenType
	val       string
}

type TestData struct {
	css    string
	tokens []TestDataToken
}

func TestNextToken(t *testing.T) {
	values := []TestData{
		{" ", []TestDataToken{{WhitespaceToken, " "}}},
		{"5.2", []TestDataToken{{NumberToken, "5.2"}}},
		{"50.12", []TestDataToken{{NumberToken, "50.12"}}},
		{".25", []TestDataToken{{NumberToken, ".25"}}},
		{"5.103 1.02", []TestDataToken{{NumberToken, "5.103"}, {WhitespaceToken, " "}, {NumberToken, "1.02"}}},
		{"5.103 .02", []TestDataToken{{NumberToken, "5.103"}, {WhitespaceToken, " "}, {NumberToken, ".02"}}},
		{".103 .02", []TestDataToken{{NumberToken, ".103"}, {WhitespaceToken, " "}, {NumberToken, ".02"}}},
		{"color: red;", []TestDataToken{{IdentToken, "color"}, {ColonToken, ":"}, {WhitespaceToken, " "}, {IdentToken, "red"}, {SemicolonToken, ";"}}},
	}

	for _, v := range values {
		l := Lexer{Text: []rune(v.css)}
		tokens := []Token{}
		for token, i := l.NextToken(), 0; token.Type != EOF; token, i = l.NextToken(), i+1 {
			tokens = append(tokens, token)
		}

		if len(tokens) != len(v.tokens) {
			t.Fatalf("len(tokens) - %v != %v - len(v.tokens)", len(tokens), len(v.tokens))
		}

		for i, token := range tokens {
			actualTokenType := token.Type
			expectedTokenType := v.tokens[i].tokenType
			actualTokenValue := string(token.Val)
			expectedTokenValue := v.tokens[i].val

			if actualTokenValue != expectedTokenValue {
				t.Fatalf("%s %d. token value (expected) %v != %v (actual)", v.css, i, expectedTokenValue, actualTokenValue)
			}

			if actualTokenType != expectedTokenType {
				t.Fatalf("%s %d. token type (expected) %v != %v (actual)", v.css, i, expectedTokenType.String(), actualTokenType.String())
			}
		}
	}
}
