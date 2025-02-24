package kamailio_cfg

// Analyzer is a struct that holds the components necessary for analyzing Kamailio configurations.
// It contains a builder for constructing the AST and a reference to the root AST node.
type Analyzer struct {
	builder *KamailioASTBuilder
	ast     *ASTNode
}

// NewAnalyzer creates and returns a new instance of Analyzer.
// It initializes the builder field with a new KamailioASTBuilder.
func NewAnalyzer() *Analyzer {
	return &Analyzer{
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
	a.ast = a.builder.BuildAST(content)
}

// GetAST returns the root AST (Abstract Syntax Tree) node that was built by the analyzer.
//
// Returns:
//
//	*ASTNode - The root node of the AST.
func (a *Analyzer) GetAST() *ASTNode {
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
