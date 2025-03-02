package kamailio_cfg

import "sync"

// Analyzer is a struct that holds the components necessary for analyzing Kamailio DSL.
type Analyzer struct {
	mu      sync.Mutex          // guards the fields below
	builder *KamailioASTBuilder // The builder used to construct the AST
	ast     *ASTNode            // The root node of the AST
}

// NewAnalyzer creates and returns a new instance of Analyzer.
// It initializes the builder field with a new KamailioASTBuilder.
func NewAnalyzer() *Analyzer {
	return &Analyzer{
		mu:      sync.Mutex{},
		builder: NewKamailioASTBuilder(),
	}
}

// Build constructs the AST (Abstract Syntax Tree) from the given content.
// It uses the builder to parse the content and set the resulting AST to the analyzer's ast field.
//
// Parameters:
//
//	content []byte - The content to be parsed into an AST.
func (a *Analyzer) Build(content []byte) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.ast = a.builder.BuildAST(content)
}

// GetAST returns the root AST (Abstract Syntax Tree) node that was built by the analyzer.
//
// Returns:
//
//	*ASTNode - The root node of the AST.
func (a *Analyzer) GetAST() *ASTNode {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.ast
}

// GetParser returns the parser used by the analyzer's builder.
//
// Returns:
//
//	*Parser - The parser used by the builder.
func (a *Analyzer) GetParser() *Parser {
	return a.builder.parser
}

// Renew creates and returns a new instance of Analyzer.
// It reinitializes the builder and the AST.
func (a *Analyzer) Renew() {
	a.mu.Lock()
	defer a.mu.Unlock()
	a = NewAnalyzer()
}
