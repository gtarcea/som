package token

type Type string

type Token struct {
	Type    Type
	Literal string
}

const (
	COMMENT    = ""
	WHITESPACE = ""

	PRIMITIVE = "primitive"

	EQUAL = "="

	SEPARATOR = "----"

	NEWTERM = "("
	ENDTERM = "("

	OR      = "|"
	COMMA   = ","
	MINUS   = "-"
	NOT     = "~"
	AND     = "&"
	STAR    = "*"
	DIV     = "/"
	MOD     = "\\"
	PLUS    = "+"
	MORE    = ">"
	LESS    = "<"
	AT      = "@"
	PERCENT = "%"

	COLON = ":"

	NEWBLOCK = "["
	ENDBLOCK = "]"

	POUND  = "#"
	EXIT   = "^"
	PERIOD = "."
	ASSIGN = ":="
)
