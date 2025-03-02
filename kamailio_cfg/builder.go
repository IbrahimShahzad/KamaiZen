package kamailio_cfg

import (
	sitter "github.com/smacker/go-tree-sitter"
)

// ASTNode represents a node in the Abstract Syntax Tree (AST).
// It contains a reference to a sitter.Node, which is a node in the tree-sitter parsing library.
type ASTNode struct {
	Node *sitter.Node
}

// IsError checks if the ASTNode represents an error node.
//
// Returns:
//
//	bool - True if the node is an error node, false otherwise.
func (a *ASTNode) IsError() bool {
	return a.Node.IsError()
}

// StartPoint returns the node starting point.
//
// Returns:
//
//	sitter.Point - The starting point of the node
func (a *ASTNode) StartPoint() sitter.Point {
	return a.Node.StartPoint()
}

// EndPoint returns the node end point.
//
// Returns:
//
//	sitter.Point - The end point of the node
func (a *ASTNode) EndPoint() sitter.Point {
	return a.Node.EndPoint()
}

// ChildCount returns the number of child nodes the ASTNode has
//
// Returns:
//
//	int - The number of child nodes
func (a *ASTNode) ChildCount() int {
	return int(a.Node.ChildCount())
}

// NamedChildCount returns the number of named child nodes the ASTNode has
//
// Returns:
//
//	int - The number of named child nodes
func (a *ASTNode) NamedChildCount() int {
	return int(a.Node.NamedChildCount())
}

func (a *ASTNode) Accept(v ASTVisitor, analyzer *Analyzer) {
	analyzer.mu.Lock()
	v.Visit(a, analyzer)
	analyzer.mu.Unlock()
}

// KamailioASTBuilder is responsible for building the Abstract Syntax Tree (AST) for Kamailio configurations.
// It contains a parser that is used to parse the configuration content.
type KamailioASTBuilder struct {
	parser *Parser
}

// NewKamailioASTBuilder creates and returns a new instance of KamailioASTBuilder.
// It initializes the builder with a new parser.
//
// Returns:
//
//	*KamailioASTBuilder - A new instance of KamailioASTBuilder.
func NewKamailioASTBuilder() *KamailioASTBuilder {
	return &KamailioASTBuilder{
		parser: NewParser(),
	}
}

// BuildAST parses the given source code and constructs the Abstract Syntax Tree (AST).
// It uses the parser to parse the source code and returns the root AST node.
//
// Parameters:
//
//	sourceCode []byte - The source code to be parsed into an AST.
//
// Returns:
//
//	*ASTNode - The root node of the constructed AST, or nil if parsing fails.
func (k *KamailioASTBuilder) BuildAST(sourceCode []byte) *ASTNode {
	n := k.parser.Parse(sourceCode)
	if n == nil {
		return nil
	}
	return &ASTNode{
		Node: n,
	}
}
