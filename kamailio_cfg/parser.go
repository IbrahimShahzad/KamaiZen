package kamailio_cfg

import (
	"context"

	"github.com/rs/zerolog/log"
	sitter "github.com/smacker/go-tree-sitter"
)

// TODO: Implement NotificationHandle for the parser and visitor functions.
// The AST should be parsed and updated in the same instance of the parser.
// This is because the parser instance keeps track of the tree and the tree
// is updated in place. If a new instance of the parser is created, the tree
// will be parsed from scratch and the previous tree will be lost.
// This is why the parser instance is stored in the analyser and the diagnostics
// functions.

// Watch document states and update the parser tree accordingly.
// Document state would watch stateTree and update the document state
// Whenever the docuemnt State is updated, the parser tree should be updated and diagnostics should be generated.

// Implement Memento for parser to undo changes in case of error.

// Parser is a struct that encapsulates the tree-sitter parser and its associated state.
// It holds references to the current and previous parse trees, as well as the language
// used for parsing.
type Parser struct {
	parser   *sitter.Parser
	language *sitter.Language
	tree     *sitter.Tree
	oldTree  *sitter.Tree
}

// NewParser creates and returns a new instance of Parser.
// It initializes the tree-sitter parser and sets the language for parsing.
//
// Returns:
//
//	*Parser - A new instance of Parser.
func NewParser() *Parser {
	log.Debug().Msg("Creating new parser")
	p := Parser{
		parser:   sitter.NewParser(),
		tree:     nil,
		oldTree:  nil,
		language: sitter.NewLanguage(Language()),
	}
	p.parser.SetLanguage(p.language)
	return &p
}

// GetLanguage returns the language used by the parser.
//
// Returns:
//
//	*sitter.Language - The language used by the parser.
func (p *Parser) GetLanguage() *sitter.Language {
	return p.language
}

// Parse parses the given source code and constructs the Abstract Syntax Tree (AST).
// It updates the current and previous parse trees accordingly.
//
// Parameters:
//
//	sourceCode []byte - The source code to be parsed into an AST.
//
// Returns:
//
//	*sitter.Node - The root node of the constructed AST, or nil if parsing fails.
func (p *Parser) Parse(sourceCode []byte) *sitter.Node {
	if p.language == nil {
		log.Fatal().Msg("Parser not initialized")
		return nil
	}
	tree, err := p.parser.ParseCtx(context.Background(), p.oldTree, sourceCode)
	if err != nil {
		log.Error().Err(err).Msg("Error parsing the source code")
		p.oldTree = nil
		return nil
	}
	if p.tree == nil {
		p.tree = tree
	} else {
		p.oldTree = p.tree
		p.tree = tree
	}
	n := p.tree.RootNode()
	return n
}

// GetTree returns the current parse tree.
//
// Returns:
//
//	*sitter.Tree - The current parse tree.
func (p *Parser) GetTree() *sitter.Tree {
	return p.tree
}
