package parser

import (
	"github.com/gtarcea/som/internal/lexer"
	"github.com/gtarcea/som/internal/token"
	"github.com/hashicorp/go-multierror"
)

type Parser struct {
	l            *lexer.Lexer
	errors       multierror.Error
	currentToken token.Token
	peekToken    token.Token
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l: l,
	}

	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) nextToken() {
	p.currentToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) Parse() {
	return
}
