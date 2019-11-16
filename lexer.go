package main

// regex to match numbers
// var number = regexp.MustCompile("^[-+]?[0-9]+.?[[0-9]*]?$")

// Token is a lexical component the Cog programming language.
type Token int

// List of the tokens which make up Cog
const (
	ILLIGAL Token = iota
	EOF

	WORD  // main
	QUOTE // \main
	STR   // "hello"
	NUM   // 234
)

type Lexer struct {
	input        string
	position     int  // current position in input (points to current char)
	readPosition int  // current reading position in input (after current char)
	ch           byte // current char under examination
}

func NewLexer() *Lexer {
	l := &Lexer{}
	return l
}

type Lexeme struct {
	Type    Token
	Literal string
}

func (l *Lexer) SetInput(input string) {
	l.input = input
	l.position = 0
	l.readPosition = 0
	l.readChar()
	return
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}

	l.position = l.readPosition
	l.readPosition++
}

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return '0'
	}
	return l.input[l.readPosition]
}

func (l *Lexer) NextLexeme() Lexeme {
	var lex Lexeme

	l.skipWhitespace()

	switch l.ch {
	case '"':
		return l.readString()
	case 92: //bashslash
		return l.readQuote()
	case 0:
		lex.Literal = ""
		lex.Type = EOF
	default:
		return l.readWordOrNumber()
	}

	l.readChar()
	return lex
}

func (l *Lexer) readWordOrNumber() Lexeme {
	var lex Lexeme
	pos := l.position
	for !isWhitespace(l.ch) {
		l.readChar()
	}
	lex.Literal = l.input[pos:l.position]
	if number.Match([]byte(lex.Literal)) {
		lex.Type = NUM
	} else {
		lex.Type = WORD
	}
	return lex
}

func (l *Lexer) readQuote() Lexeme {
	var lex Lexeme
	pos := l.position + 1 //skip the backslash
	for !isWhitespace(l.ch) {
		l.readChar()
	}
	lex.Type = QUOTE
	lex.Literal = l.input[pos:l.position]
	return lex
}

// func (l *Lexer) readVar() Lexeme {
// 	var lex Lexeme
// 	pos := l.position + 1
// 	for !isWhitespace(l.ch) && !isGrouping(l.ch) {
// 		l.readChar()
// 	}
// 	lex.Type = VAL
// 	lex.Literal = l.input[pos:l.position]
// 	return lex
// }

// func (l *Lexer) readNumber() Lexeme {
// 	pos := l.position
// 	if l.ch == '+' || l.ch == '-' {
// 		l.readChar()
// 	}
// 	for isDigit(l.ch) || l.ch == '.' {
// 		l.readChar()
// 	}
// 	return Lexeme{
// 		NUM,
// 		l.input[pos:l.position],
// 	}
// }

func (l *Lexer) skipWhitespace() {
	for isWhitespace(l.ch) {
		l.readChar()
	}
}

func (l *Lexer) readString() Lexeme {
	pos := l.position + 1
	for {
		l.readChar()
		if l.ch == '"' {
			break
		}
	}
	lex := Lexeme{
		STR,
		l.input[pos:l.position],
	}
	l.readChar() // skip passed the " we stopped on
	return lex
}

func isWhitespace(ch byte) bool {
	return ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r'
}

// func isDigit(ch byte) bool {
// 	return '0' <= ch && ch <= '9'
// }

// func isGrouping(ch byte) bool {
// 	return ch == '[' || ch == '{' || ch == '(' || ch == ']' || ch == '}' || ch == ')'
// }
