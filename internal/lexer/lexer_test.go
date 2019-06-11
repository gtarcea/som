package lexer

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/gtarcea/som/internal/token"
)

func TestNextToken(t *testing.T) {
	input := `=::= 'hello' 'hello\'' 123 123.3
`
	tests := []struct {
		expectedTokenType token.Type
		expectedLiteral   string
	}{
		{token.EQUAL, "="},
		{token.COLON, ":"},
		{token.ASSIGN, ":="},
		{token.STRING, "'hello'"},
		{token.STRING, "'hello\\''"},
		{token.INTEGER, "123"},
		{token.DOUBLE, "123.3"},
	}

	l := NewLexer(input)
	for _, test := range tests {
		tok := l.NextToken()
		t.Run(test.expectedLiteral, func(t *testing.T) {
			require.Equalf(t, tok.Type, test.expectedTokenType, "Expected Token Type %s got %s", tok.Type, test.expectedTokenType)
			require.Equalf(t, tok.Literal, test.expectedLiteral, "Expected Token Literal %s got %s", tok.Literal, test.expectedLiteral)
		})
	}
}
