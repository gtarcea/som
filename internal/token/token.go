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

	STRING = "string"

	EQUAL = "="

	SEPARATOR = "----"

	NEWTERM = "("
	ENDTERM = ")"

	MINUS = "-"

	// Operators
	OR      = "|"
	NOT     = "~"
	AND     = "&"
	MULT    = "*"
	DIV     = "/"
	MOD     = "\\"
	PLUS    = "+"
	MORE    = ">"
	LESS    = "<"
	AT      = "@"
	PERCENT = "%"
	COMMA   = ","

	SINGLE_QUOTE = "'"

	COLON = ":"

	NEWBLOCK = "["
	ENDBLOCK = "]"

	POUND  = "#"
	EXIT   = "^"
	PERIOD = "."
	ASSIGN = ":="

	IDENTIFIER = "identifier"
)
