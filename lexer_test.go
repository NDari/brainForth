package main

import (
	"testing"
)

func TestNextToken(t *testing.T) {
	input := `
		1    +2 -3          add "things and such" * \dup  -3 22 :
		{ } ( ) [ [ 2 ] ]
	`
	tests := []struct {
		testNumber      int
		expectedType    Token
		expectedLiteral string
	}{
		{1, NUM, "1"},
		{2, NUM, "+2"},
		{2, NUM, "-3"},
		{3, WORD, "add"},
		{4, STR, "things and such"},
		{5, WORD, "*"},
		{6, QUOTE, "dup"},
		{7, NUM, "-3"},
		{7, NUM, "22"},
		{8, WORD, ":"},
		{9, WORD, "{"},
		{10, WORD, "}"},
		{11, WORD, "("},
		{12, WORD, ")"},
		{13, WORD, "["},
		{14, WORD, "["},
		{15, NUM, "2"},
		{16, WORD, "]"},
		{17, WORD, "]"},
		{15, EOF, ""},
	}
	l := NewLexer()
	l.SetInput(input)

	for _, tt := range tests {
		lex := l.NextLexeme()
		t.Log(lex)

		if lex.Type != tt.expectedType {
			t.Fatalf("test %d - tokentype wrong. expected=%d, got=%d", tt.testNumber, tt.expectedType, lex.Type)
		}

		if lex.Literal != tt.expectedLiteral {
			t.Fatalf("test %d - literal wrong. expected=%s, got=%s", tt.testNumber, tt.expectedLiteral, lex.Literal)
		}
	}
}
