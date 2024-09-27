package parser

import (
	"fmt"
	"regexp"
	"strconv"
)

type SyntaxError struct {
	Msg string
}

func (e *SyntaxError) Error() string {
	return e.Msg
}

// TokenType represents the type of token
type TokenType string

const (
	QUOTE      byte = '"'
	SYMBOL     byte = '$'
	OPEN_PAREN byte = '('

	EQ  byte = '='
	NE  byte = '!'
	LT  byte = '<'
	GT  byte = '>'
	AND byte = '&'
	OR  byte = '|'
	NOT byte = '!'
	ADD byte = '+'
	SUB byte = '-'
	MUL byte = '*'
	DIV byte = '/'
	MOD byte = '%'

	HASH byte = '#'
	BANG byte = '!'
)

var keywords = []string{
	"return",
	"break",
	"case",
	"const",
	"continue",
	"default",
	"else",
	"if",
	"for",
	"switch",
	"while",
	"include_file",
	"modparam",
}

var defines = []string{
	"define",
	"redefine",
	"undef",
	"ifdef",
	"ifndef",
	"trydef",
	"substdefs",
	"substdef",
	"subst",
	"else",
	"endif",
	"include_file", // the preprocessor directive is optional for these two
	"import_file",
}

const (
	ILLEGAL TokenType = "ILLEGAL"
	EOF     TokenType = "EOF"
	IDENT   TokenType = "IDENT"
	INT     TokenType = "INT"
	STRING  TokenType = "STRING"
	// Keywords
	KEYWORD TokenType = "KEYWORD"
	PREPROC TokenType = "PREPROC"

	// cfg specific keywords
	REQUEST_ROUTE TokenType = "REQUEST_ROUTE"
	REPLY_ROUTE   TokenType = "REPLY_ROUTE"
	ROUTE         TokenType = "ROUTE"

	// Operators and symbols
	ELLIPSIS TokenType = "ELLIPSIS"
	// Core Variables
	CORE_VARIABLE TokenType = "CORE_VARIABLE"

	// File specific
	LOADMODULE TokenType = "LOADMODULE"
	COMMENT    TokenType = "COMMENT"
	NEWLINE    TokenType = "NEWLINE"

	// Operators
	ASSIGN TokenType = "ASSIGN" // =
	EQ_OP  TokenType = "EQ_OP"  // ==
	NE_OP  TokenType = "NE_OP"  // !=
	LT_OP  TokenType = "LT_OP"  // <
	LE_OP  TokenType = "LE_OP"  // <=
	GT_OP  TokenType = "GT_OP"  // >
	GE_OP  TokenType = "GE_OP"  // >=
	AND_OP TokenType = "AND_OP" // &&
	OR_OP  TokenType = "OR_OP"  // ||
	NOT_OP TokenType = "NOT_OP" // !
	ADD_OP TokenType = "ADD_OP" // +
	SUB_OP TokenType = "SUB_OP" // -
	MUL_OP TokenType = "MUL_OP" // *
	DIV_OP TokenType = "DIV_OP" // /
	MOD_OP TokenType = "MOD_OP" // %
	INC_OP TokenType = "INC_OP" // ++
	DEC_OP TokenType = "DEC_OP" // --

	// Punctuation
	COMMA     TokenType = "COMMA"     // ,
	SEMICOLON TokenType = "SEMICOLON" // ;
	LPAREN    TokenType = "LPAREN"    // (
	RPAREN    TokenType = "RPAREN"    // )
	LBRACE    TokenType = "LBRACE"    // {
	RBRACE    TokenType = "RBRACE"    // }
	LBRACKET  TokenType = "LBRACKET"  // [
	RBRACKET  TokenType = "RBRACKET"  // ]
	DOT       TokenType = "DOT"       // .
	COLON     TokenType = "COLON"     // :
	// TODO: verify if needed
	// ARROW    TokenType = "ARROW"     // ->
	// DOUBLECOLON TokenType = "DOUBLECOLON" // ::

)

// Token should be an interface
// with TokenType and Literal methods
type Token interface {
	Type() TokenType
	Literal() interface{}
}

// BasicToken represents a simple token with a type and literal value
type BasicToken struct {
	TypeVal    TokenType
	LiteralVal string
}

// Type returns the token type
func (t *BasicToken) Type() TokenType {
	return t.TypeVal
}

// Literal returns the literal value
func (t *BasicToken) Literal() interface{} {
	return t.LiteralVal
}

type CoreVariableToken struct {
	TypeVal      TokenType
	VariableType string
	VariableName string
}

func (c *CoreVariableToken) Type() TokenType {
	return c.TypeVal
}

func (c *CoreVariableToken) Literal() interface{} {
	if c.VariableName == "" {
		return c.VariableType
	}
	return c.VariableType + "," + c.VariableName
}

// Numberliteral represents a number literal
// with TokenType as INT and Literal as the value
type NumberLiteral struct {
	Value int
}

func (n *NumberLiteral) Literal() interface{} {
	return n.Value
}

func (n *NumberLiteral) Type() TokenType {
	return INT
}

type EOFToken struct{}

func (e *EOFToken) Literal() interface{} {
	return nil
}

func (e *EOFToken) Type() TokenType {
	return EOF
}

type Lexer struct {
	// lexer fields
	input   []byte // input string
	pos     int    // current position in input (points to current char)
	readPos int    // next position to read
	ch      byte   // current char under examination
}

// NewLexer initializes a new Lexer with the given input
func NewLexer(input []byte) *Lexer {
	l := &Lexer{input: input}
	l.next()
	return l
}

// next reads the next character in the input and advances the positions
func (l *Lexer) next() {
	if l.readPos >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPos]
	}
	l.pos = l.readPos
	l.readPos++
}

func (l *Lexer) peek() byte {
	if l.readPos >= len(l.input) {
		return 0
	}
	return l.input[l.readPos]
}

func (l *Lexer) advanceN(n int) {
	for i := 0; i < n; i++ {
		l.next()
	}
}

// isLetter checks if a character is a letter
func isLetter(ch byte) bool {
	return ('a' <= ch && ch <= 'z') || ('A' <= ch && ch <= 'Z')
}

type TokenHandler func(string) Token

// Handler for INT tokens (returns a NumberLiteral)
func handleInt(match string) Token {
	intValue, err := strconv.Atoi(match)
	if err != nil {
		return &NumberLiteral{Value: 0} // Handle invalid integer
	}
	return &NumberLiteral{Value: intValue}
}

func handleStringToken(tokenType TokenType) TokenHandler {
	return func(match string) Token {
		return &BasicToken{TypeVal: tokenType, LiteralVal: match[1 : len(match)-1]}
	}
}

// Handler for IDENT and other basic tokens
func handleBasicToken(tokenType TokenType) TokenHandler {
	return func(match string) Token {
		return &BasicToken{TypeVal: tokenType, LiteralVal: match}
	}
}

type TokenRegexHandler struct {
	TokenType    TokenType
	TokenHandler TokenHandler
}

var mapRegexHandler = map[*regexp.Regexp]TokenRegexHandler{
	regexp.MustCompile(`^\d+`): {
		TokenType:    INT,
		TokenHandler: handleInt,
	},
	regexp.MustCompile(`^[_a-zA-Z][_a-zA-Z0-9]*`): {
		TokenType:    IDENT,
		TokenHandler: handleBasicToken(IDENT),
	},
	regexp.MustCompile(`^"[^"]*"`): {
		TokenType:    STRING,
		TokenHandler: handleStringToken(STRING),
	},
	// punctuation
	regexp.MustCompile(`^,`): {
		TokenType:    COMMA,
		TokenHandler: handleBasicToken(COMMA),
	},
	regexp.MustCompile(`^;`): {
		TokenType:    SEMICOLON,
		TokenHandler: handleBasicToken(SEMICOLON),
	},
	regexp.MustCompile(`^\(`): {
		TokenType:    LPAREN,
		TokenHandler: handleBasicToken(LPAREN),
	},
	regexp.MustCompile(`^\)`): {
		TokenType:    RPAREN,
		TokenHandler: handleBasicToken(RPAREN),
	},
	regexp.MustCompile(`^\{`): {
		TokenType:    LBRACE,
		TokenHandler: handleBasicToken(LBRACE),
	},
	regexp.MustCompile(`^\}`): {
		TokenType:    RBRACE,
		TokenHandler: handleBasicToken(RBRACE),
	},
	regexp.MustCompile(`^\[`): {
		TokenType:    LBRACKET,
		TokenHandler: handleBasicToken(LBRACKET),
	},
	regexp.MustCompile(`^\]`): {
		TokenType:    RBRACKET,
		TokenHandler: handleBasicToken(RBRACKET),
	},
	regexp.MustCompile(`^\.`): {
		TokenType:    DOT,
		TokenHandler: handleBasicToken(DOT),
	},
	regexp.MustCompile(`^:`): {
		TokenType:    COLON,
		TokenHandler: handleBasicToken(COLON),
	},
}

var ignoredRegex = map[*regexp.Regexp]struct{}{
	regexp.MustCompile(`^\s+`):                           {}, // Whitespace
	regexp.MustCompile(`^/\*[^*]*\*+([^\/][^*]*\*+)*\/`): {}, // Multi-line comments
	regexp.MustCompile(`^//.*`):                          {}, // Single-line comments
	regexp.MustCompile(`^#[^!].*`):                       {}, // Single-line comments
}

func (l *Lexer) readOperator() Token {
	// for operators that are two characters long
	switch l.ch + l.peek() {
	case (EQ + EQ):
		l.next()
		l.next()
		return &BasicToken{
			TypeVal:    EQ_OP,
			LiteralVal: string(EQ) + string(EQ),
		}
	case (NOT + EQ):
		l.next()
		l.next()
		return &BasicToken{
			TypeVal:    NE_OP,
			LiteralVal: string(NE) + string(EQ),
		}
	case (LT + EQ):
		l.next()
		l.next()
		return &BasicToken{
			TypeVal:    LE_OP,
			LiteralVal: string(LT) + string(EQ),
		}
	case (GT + EQ):
		l.next()
		l.next()
		return &BasicToken{
			TypeVal:    GE_OP,
			LiteralVal: string(GT) + string(EQ),
		}
	case (AND + AND):
		l.next()
		l.next()
		return &BasicToken{
			TypeVal:    AND_OP,
			LiteralVal: string(AND) + string(AND),
		}
	case (OR + OR):
		l.next()
		l.next()
		return &BasicToken{
			TypeVal:    OR_OP,
			LiteralVal: string(OR) + string(OR),
		}
	case (ADD + ADD):
		l.next()
		l.next()
		return &BasicToken{
			TypeVal:    INC_OP,
			LiteralVal: string(ADD) + string(ADD),
		}
	case (SUB + SUB):
		l.next()
		l.next()
		return &BasicToken{
			TypeVal:    DEC_OP,
			LiteralVal: string(SUB) + string(SUB),
		}
	}
	// for operators that are one character long
	switch l.ch {
	case EQ:
		l.next()
		return &BasicToken{
			TypeVal:    ASSIGN,
			LiteralVal: string(EQ),
		}
	case NOT:
		l.next()
		return &BasicToken{
			TypeVal:    NE_OP,
			LiteralVal: string(NE),
		}
	case LT:
		l.next()
		return &BasicToken{
			TypeVal:    LT_OP,
			LiteralVal: string(LT),
		}
	case GT:
		l.next()
		return &BasicToken{
			TypeVal:    GT_OP,
			LiteralVal: string(GT),
		}
	// TODO: Bitwise operators
	// case AND:
	// 	l.next()
	// 	return &BasicToken{
	// 		TypeVal:    AND_OP,
	// 		LiteralVal: string(AND),
	// 	}
	case ADD:
		l.next()
		return &BasicToken{
			TypeVal:    ADD_OP,
			LiteralVal: string(ADD),
		}
	case SUB:
		l.next()
		return &BasicToken{
			TypeVal:    SUB_OP,
			LiteralVal: string(SUB),
		}
	case MUL:
		l.next()
		return &BasicToken{
			TypeVal:    MUL_OP,
			LiteralVal: string(MUL),
		}
	case DIV:
		l.next()
		return &BasicToken{
			TypeVal:    DIV_OP,
			LiteralVal: string(DIV),
		}
	}
	return nil
}

func (l *Lexer) readMacros() Token {
	switch l.ch + l.peek() {
	case (HASH + BANG):
		fallthrough
	case (BANG + BANG):
		l.next()
		l.next()
		re := regexp.MustCompile(`^[_a-zA-Z][_a-zA-Z0-9]*`)
		m := re.Find(l.input[l.pos:])
		l.advanceN(len(m))
		return &BasicToken{
			TypeVal:    PREPROC,
			LiteralVal: string(m),
		}
	}
	return nil
}

func (l *Lexer) readString() string {
	// Skip the opening quote
	// multi line strings span multiple lines
	// each line with " at the start and end
	// example: "multi-line"\n\t\t\t"string"
	pos := l.pos + 1
	for {
		l.next()
		if l.ch == '"' || l.ch == 0 { // End of string
			break
		}
	}
	l.next()
	return string(l.input[pos : l.pos-1])
}

func (l *Lexer) handleMultiLineString() Token {
	token := &BasicToken{
		TypeVal:    STRING,
		LiteralVal: l.readString(),
	}
	return token
}

func sanitizeTokens(tokens []Token) []Token {
	var sanitized []Token
	for _, token := range tokens {
		if token.Type() == NEWLINE {
			continue
		}

		// if two consecutive strings, merge them
		if len(sanitized) > 0 {
			if sanitized[len(sanitized)-1].Type() == STRING && token.Type() == STRING {
				sanitized[len(sanitized)-1].(*BasicToken).LiteralVal += token.Literal().(string)
				continue
			}
		}

		sanitized = append(sanitized, token)
	}

	// check identifier for keywords
	for _, token := range sanitized {
		if token.Type() == IDENT {
			if isOneOfKeywords(token.Literal().(string)) {
				token.(*BasicToken).TypeVal = KEYWORD
				continue
			}
			if token.(*BasicToken).isOneOfMany(defines...) {
				token.(*BasicToken).TypeVal = PREPROC
			}
		}
	}
	return sanitized
}

func (t *BasicToken) isOneOfMany(s ...string) bool {
	for _, keyword := range s {
		if t.LiteralVal == keyword {
			return true
		}
	}
	return false
}

func isOneOfKeywords(s string) bool {
	for _, keyword := range keywords {
		if s == keyword {
			return true
		}
	}
	return false
}

func (l *Lexer) handleCoreVariable() Token {
	// corevariable: $var(variable_name)
	// corevariable: $avp(variable_name)
	// corevariable: $xav(variable_name)
	// skip the $
	var variableType string
	var variableName string
	var isCoreVariable bool
	l.next()
	// it could be either an expression or a core variable
	// if it's a variable it should be followed by a letter
	if isLetter(l.ch) {
		isCoreVariable = true
		// it's a core variable
		// save the name of the variable
		r := regexp.MustCompile(`^[_a-zA-Z][_a-zA-Z0-9]*`)
		m := r.Find(l.input[l.pos:])
		l.advanceN(len(m))
		variableType = string(m)
		if l.ch == OPEN_PAREN {
			l.next()
			m = r.Find(l.input[l.pos:])
			l.advanceN(len(m))
			variableName = string(m)
			l.next()

		}
	}
	if isCoreVariable {
		return &CoreVariableToken{
			TypeVal:      CORE_VARIABLE,
			VariableType: variableType,
			VariableName: variableName,
		}
	}
	return nil
}

// Tokenise tokenizes the input and returns a slice of tokens
func (l *Lexer) Tokenise() []Token {
	var tokens []Token
	var isStringParsing bool

	// 0 is the null character
	for l.ch != 0 {
		matched := false

		// If currently parsing a string, consume until the closing quote
		if isStringParsing {
			if l.ch == QUOTE {
				tokens = append(tokens, l.handleMultiLineString())
				isStringParsing = false // End string parsing mode
				matched = true
				continue
			}
		}

		// ignore pattern
		for re := range ignoredRegex {
			if re.Match(l.input[l.pos:]) {
				matched = true
				// skip over it
				m := re.Find(l.input[l.pos:])
				l.advanceN(len(m))
				continue
			}
		}

		if !matched {
			for re, handler := range mapRegexHandler {
				if re.Match(l.input[l.pos:]) {
					m := re.Find(l.input[l.pos:])
					token := handler.TokenHandler(string(m))
					tokens = append(tokens, token)
					l.advanceN(len(m))
					if token.Type() == STRING {
						isStringParsing = true
					}
					matched = true
					break
				}
			}
		}

		if l.ch == SYMBOL {
			x := l.handleCoreVariable()
			if x != nil {
				tokens = append(tokens, x)
				matched = true
			}
		}

		// NOTE:
		// apply a list of functions unless one of
		// them returns a valid token
		// think of a clever way to do this!

		if !matched {
			if token := l.readOperator(); token != nil {
				tokens = append(tokens, token)
				matched = true
			}
		}

		if !matched {
			if token := l.readMacros(); token != nil {
				tokens = append(tokens, token)
				matched = true
			}
		}

		// If nothing matched, it's an ILLEGAL token
		if !matched {
			fmt.Printf("Illegal character: %c , at %s\n", l.ch, l.input[l.pos:])
			tokens = append(tokens, &BasicToken{TypeVal: ILLEGAL, LiteralVal: string(l.ch)})
			// l.next()
			// TODO: remove panic and handle it properly
			panic(&SyntaxError{Msg: "Syntax error:" + string(l.ch)})
		}
	}

	return append(sanitizeTokens(tokens), &EOFToken{})
}
