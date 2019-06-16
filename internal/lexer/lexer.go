package lexer

import (
	"strings"

	"github.com/gtarcea/som/internal/token"
)

type Lexer struct {
	input           string
	currentPosition int
	readPosition    int
	char            byte
}

func NewLexer(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

func (l *Lexer) NextToken() token.Token {
	var t token.Token

	l.skipWhitespace()

	switch {
	case l.charIs('='):
		t = l.newTokenFromChar(token.EQUAL)
	case l.charIs(':'):
		t = l.lexColon()
	case l.charIs('\''):
		t = l.lexString()
	case l.charIs('-'):
		t = l.lexMinus()
	case l.charIs('|'):
		t = l.newTokenFromChar(token.OR)
	case l.charIs('~'):
		t = l.newTokenFromChar(token.NOT)
	case l.charIs('&'):
		t = l.newTokenFromChar(token.AND)
	case l.charIs('*'):
		t = l.newTokenFromChar(token.MULT)
	case l.charIs('/'):
		t = l.newTokenFromChar(token.DIV)
	case l.charIs('\\'):
		t = l.newTokenFromChar(token.MOD)
	case l.charIs('+'):
		t = l.newTokenFromChar(token.PLUS)
	case l.charIs('>'):
		t = l.newTokenFromChar(token.MORE)
	case l.charIs('<'):
		t = l.newTokenFromChar(token.LESS)
	case l.charIs('@'):
		t = l.newTokenFromChar(token.AT)
	case l.charIs('%'):
		t = l.newTokenFromChar(token.PERCENT)
	case l.charIs(','):
		t = l.newTokenFromChar(token.COMMA)
	case l.charIs('['):
		t = l.newTokenFromChar(token.NEWBLOCK)
	case l.charIs(']'):
		t = l.newTokenFromChar(token.ENDBLOCK)
	case l.charIs('('):
		t = l.newTokenFromChar(token.NEWTERM)
	case l.charIs(')'):
		t = l.newTokenFromChar(token.ENDTERM)
	case l.charIs('#'):
		t = l.newTokenFromChar(token.POUND)
	case l.charIs('^'):
		t = l.newTokenFromChar(token.EXIT)
	case l.charIs('.'):
		t = l.newTokenFromChar(token.PERIOD)
	case isLetter(l.char):
		t = l.lexIdentifierOrPrimitive()
	case isDigit(l.char):
		t = l.lexDigit()
	}

	l.readChar()
	return t
}

func (l *Lexer) charIs(c byte) bool {
	return c == l.char
}

func (l *Lexer) lexColon() token.Token {
	if l.peekChar() == '=' {
		char := l.char
		l.readChar()
		literal := string(char) + string(l.char)
		return token.Token{Type: token.ASSIGN, Literal: literal}
	}

	return newToken(token.COLON, l.char)
}

func (l *Lexer) lexMinus() token.Token {
	t := token.Token{Type: token.MINUS}
	if strings.HasPrefix(l.input[l.currentPosition:], "----") {
		t.Type = token.SEPARATOR
		var b strings.Builder
		b.WriteByte(l.char)
		for {
			char := l.peekChar()
			if char != '-' {
				break
			}
			l.readChar()
			b.WriteByte(l.char)
		}
		t.Literal = b.String()
	} else {
		t.Literal = "-"
	}

	return t
}

func (l *Lexer) lexString() token.Token {
	var b strings.Builder

	b.WriteString("'")
	for {
		l.readChar()
		if l.char == '\'' {
			break
		}
		l.lexStringChar(&b)
	}
	b.WriteString("'")
	return token.Token{Type: token.STRING, Literal: b.String()}
}

func (l *Lexer) lexStringChar(b *strings.Builder) {
	if l.char == '\\' {
		l.lexEscapeChar(b)
	} else {
		b.WriteByte(l.char)
	}
}

func (l *Lexer) lexEscapeChar(b *strings.Builder) {
	l.readChar()
	switch l.char {
	case 't':
		b.WriteString("\t")
	case 'b':
		b.WriteString("\b")
	case 'n':
		b.WriteString("\n")
	case 'r':
		b.WriteString("\r")
	case 'f':
		b.WriteString("\f")
	case '\'':
		b.WriteString("\\'")
	case '\\':
		b.WriteString("\\")
	}
}

func (l *Lexer) lexIdentifierOrPrimitive() token.Token {
	if l.char == 'p' {
		// could be the word primitive
		if strings.HasPrefix(l.input[l.currentPosition:], "primitive") {
			for l.char != 'e' {
				l.readChar()
			}

			return token.Token{Type: token.PRIMITIVE, Literal: "primitive"}
		}
	}

	// If not primitive then we are reading an identifier

	var (
		b strings.Builder
	)
	tokenType := token.IDENTIFIER

	b.WriteByte(l.char)
	for {
		char := l.peekChar()
		if !isIdentifierChar(char) {
			break
		}

		l.readChar()
		b.WriteByte(l.char)
	}

	if l.peekChar() == ':' {
		tokenType = token.KEYWORD
		l.readChar()
		b.WriteByte(l.char)
		char := l.peekChar()
		if isLetter(char) {
			tokenType = token.KEYWORD_SEQUENCE
			l.readChar()
			b.WriteByte(l.char)
			for {
				char := l.peekChar()
				if !isIdentifierChar(char) && char != ':' {
					break
				}
				l.readChar()
				b.WriteByte(l.char)
			}
		}
	}

	return token.Token{Type: token.Type(tokenType), Literal: b.String()}
}

func (l *Lexer) lexDigit() token.Token {
	var b strings.Builder
	b.WriteByte(l.char)

	t := token.Token{Type: token.INTEGER}
	sawPeriod := false

Loop:
	for {
		char := l.peekChar()
		switch {
		case char == '.' && !sawPeriod:
			char2 := l.peek2Char()
			if isDigit(char2) {
				// At this point we have seen a string that looks as follows:
				//   111.111
				// That is we've seen 1 or more numbers and then the peek shows us
				// a period, and peek2 shows a number after the period, so we are
				// reading a double. Also !sawPeriod ensures we haven't yet seen
				// a period. (The next time we see a period we know we are no longer
				// lexing a number and will exit the loop.)

				// advance lexer so that l.char == '.'
				l.readChar()
				// Now that we've seen a period make sure we don't drop back into this block
				sawPeriod = true
				// We are reading a double at this point
				t.Type = token.DOUBLE
			}
		case isDigit(char):
			// peek is a digit so advance lexer so that l.char is that digit
			l.readChar()
		default:
			// Either we are seeing a second period or a non-numeric character - either way exit loop
			break Loop
		}

		b.WriteByte(l.char)
	}

	t.Literal = b.String()

	return t
}

func (l *Lexer) skipWhitespace() {
	for l.char == ' ' || l.char == '\t' || l.char == '\n' || l.char == '\r' {
		l.readChar()
	}
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.char = 0
	} else {
		l.char = l.input[l.readPosition]
	}

	l.currentPosition = l.readPosition
	l.readPosition += 1
}

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}

	return l.input[l.readPosition]
}

func (l *Lexer) peek2Char() byte {
	if l.readPosition+1 >= len(l.input) {
		return 0
	}

	return l.input[l.readPosition+1]
}

func (l *Lexer) readIdentifier() string {
	position := l.currentPosition
	for isLetter(l.char) {
		l.readChar()
	}

	return l.input[position:l.readPosition]
}

func (l *Lexer) readNumber() string {
	position := l.currentPosition
	for isDigit(l.char) {
		l.readChar()
	}

	return l.input[position:l.readPosition]
}

func (l *Lexer) readString() string {
	position := l.currentPosition + 1
	for {
		l.readChar()
		if l.char == '"' || l.char == 0 {
			break
		}
	}

	return l.input[position:l.currentPosition]
}

func isLetter(char byte) bool {
	return 'a' <= char && char <= 'z' || 'A' <= char && char <= 'Z' || char == '_'
}

func isDigit(char byte) bool {
	return '0' <= char && char <= '9'
}

func isAlphaNumeric(char byte) bool {
	return isLetter(char) || isDigit(char)
}

func isIdentifierChar(char byte) bool {
	return isAlphaNumeric(char) || char == '_'
}

func (l *Lexer) newTokenFromChar(tokenType token.Type) token.Token {
	return token.Token{Type: tokenType, Literal: string(l.char)}
}

func newToken(tokenType token.Type, char byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(char)}
}
