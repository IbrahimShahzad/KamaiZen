package parser

// struct for the parser
type Parser struct {
	lexer *Lexer
}

// Create a new parser
func NewParser(lexer *Lexer) *Parser {
	return &Parser{lexer: lexer}
}

// Parse the input
func (p *Parser) Parse() {
	// ...
}
