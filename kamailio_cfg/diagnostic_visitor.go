package kamailio_cfg

import (
	"KamaiZen/lsp"
	"KamaiZen/settings"

	"github.com/rs/zerolog/log"
	sitter "github.com/smacker/go-tree-sitter"
)

// DiagnosticVisitor is a struct that collects diagnostics during the visit of a Kamailio configuration.
// It holds a slice of lsp.Diagnostic which contains the diagnostics found.
type DiagnosticVisitor struct {
	diagnostics []lsp.Diagnostic
}

// NewDiagnosticVisitor creates and returns a new instance of DiagnosticVisitor.
//
// Returns:
//
//	*DiagnosticVisitor - A new instance of DiagnosticVisitor.
func NewDiagnosticVisitor() *DiagnosticVisitor {
	return &DiagnosticVisitor{}
}

// createDiagnostic creates a new diagnostic message with the given parameters.
// It constructs an lsp.Diagnostic with the specified message, range, and severity.
//
// Parameters:
//
//	message string - The diagnostic message.
//	start sitter.Point - The starting position of the diagnostic range.
//	end sitter.Point - The ending position of the diagnostic range.
//	severity lsp.DiagnosticSeverity - The severity level of the diagnostic.
//
// Returns:
//
//	lsp.Diagnostic - The constructed diagnostic message.
func createDiagnostic(
	message string,
	start sitter.Point,
	end sitter.Point,
	severity lsp.DiagnosticSeverity) lsp.Diagnostic {
	return lsp.Diagnostic{
		Range: lsp.Range{
			Start: lsp.Position{
				Line:      int(start.Row),
				Character: int(start.Column),
			},
			End: lsp.Position{
				Line:      int(end.Row),
				Character: int(end.Column),
			},
		},
		Message:  message,
		Severity: severity,
	}
}

// Visit traverses the AST starting from the given node and collects diagnostics.
// It recursively visits all child nodes and adds diagnostics that can't be found using queries.
//
// Parameters:
//
//	node *ASTNode - The starting node for the traversal.
//	a *Analyzer - The analyzer used for additional context during the visit.
//
// Returns:
//
//	error - An error if the visit fails, otherwise nil.
func (d *DiagnosticVisitor) Visit(node *ASTNode, a *Analyzer) error {
	// NOTE: Add diagnostics that can't be found using queries
	// Visit and add diagnostics that can't be found using queries

	// Traverse the children
	for i := 0; i < int(node.Node.ChildCount()); i++ {
		child := node.Node.Child(i)
		d.Visit(&ASTNode{Node: child}, a)
	}

	// Query ALL deprecated comments just once
	if node.Node.NamedChildCount() == 0 {

	}
	return nil
}

func getXMLPaths(node *ASTNode, a *Analyzer) []sitter.Node {
	var xml_nodes []sitter.Node
	// TODO: fix grammar. right now skipping the xml errors
	qe, err := NewQueryExecutor(_XML_QUERY, node.Node, a.GetParser().language)
	if err != nil {
		log.Error().Err(err).Msg("Error creating query")
		return nil
	}
	for {
		match, ok := qe.NextMatch()
		if !ok {
			break
		}
		for _, capture := range match.Captures {
			node := capture.Node
			xml_nodes = append(xml_nodes, *node)
		}
	}
	return xml_nodes
}

// addSyntaxErrors identifies and collects syntax errors in the given AST node.
// It uses a query executor to find syntax errors and creates diagnostics for each error found.
//
// Parameters:
//
//	node *ASTNode - The AST node to be checked for syntax errors.
//	a *Analyzer - The analyzer used to get the parser and language information.
func (d *DiagnosticVisitor) addSyntaxErrors(node *ASTNode, a *Analyzer) {
	var diagnostics []lsp.Diagnostic
	qe, err := NewQueryExecutor(_ERROR_QUERY, node.Node, a.GetParser().language)
	if err != nil {
		log.Error().Err(err).Msg("Error creating query")
		return
	}
	xml_nodes := getXMLPaths(node, a)
	for {
		match, ok := qe.NextMatch()
		if !ok {
			break
		}
		for _, capture := range match.Captures {
			node := capture.Node
			skip := false
			for _, xml := range xml_nodes {
				if node.StartByte() >= xml.StartByte() && node.EndByte() <= xml.EndByte() {
					skip = true
				}
			}
			if skip {
				continue
			}
			diagnostics = append(diagnostics,
				createDiagnostic("Syntax error", node.StartPoint(), node.EndPoint(), lsp.ERROR))
		}
	}
	d.diagnostics = append(d.diagnostics, diagnostics...)
}

// addDeprecatedCommentHints identifies and collects deprecated comment hints in the given AST node.
// It uses a query executor to find deprecated comments and creates diagnostics for each hint found.
//
// Parameters:
//
//	node *ASTNode - The AST node to be checked for deprecated comments.
//	a *Analyzer - The analyzer used to get the parser and language information.
func (d *DiagnosticVisitor) addDeprecatedCommentHints(node *ASTNode, a *Analyzer) {
	var diagnostics []lsp.Diagnostic
	qe, err := NewQueryExecutor(_DEPRECATED_COMMENT_QUERY, node.Node, a.GetParser().language)
	if err != nil {
		log.Error().Err(err).Msg("Error creating query")
		return
	}
	for {
		match, ok := qe.NextMatch()
		if !ok {
			break
		}
		for _, capture := range match.Captures {
			node := capture.Node
			diagnostics = append(diagnostics,
				createDiagnostic("use /* comment */", node.StartPoint(), node.EndPoint(), lsp.HINT))
		}
	}
	d.diagnostics = append(d.diagnostics, diagnostics...)
}

// addUnreachableCodeWarnings identifies and collects warnings for unreachable code in the given AST node.
// It uses a query executor to find core functions and checks if any subsequent statements are unreachable.
// Diagnostics are created for each block of unreachable code found.
//
// Parameters:
//
//	node *ASTNode - The AST node to be checked for unreachable code.
//	a *Analyzer - The analyzer used to get the parser and language information.
func (d *DiagnosticVisitor) addUnreachableCodeWarnings(node *ASTNode, a *Analyzer) {
	var diagnostics []lsp.Diagnostic
	qe, err := NewQueryExecutor(_RETURN_STATEMENTS_QUERY, node.Node, a.GetParser().language)
	if err != nil {
		log.Error().Err(err).Msg("Error creating query")
		return
	}

	for {
		match, ok := qe.NextMatch()
		if !ok {
			break
		}
		for _, capture := range match.Captures {
			node := capture.Node
			s := node.Parent()
			if s.Type() == StatementNodeType && s.NextNamedSibling() != nil && s.NextNamedSibling().Type() == StatementNodeType {
				// the next named siblings (statements) are unreachable
				sibling_count := s.Parent().NamedChildCount()
				if sibling_count == 0 {
					log.Debug().Msg("No siblings found for core function")
					continue
				}
				start_node := s.NextNamedSibling()
				if start_node.Type() == BlockEndNodeType || start_node.NamedChild(0).Type() == CaseStatementNodeType {
					// not unreachable
					continue
				}
				end_node := s.Parent().NamedChild(int(sibling_count - 1))
				if end_node.Type() == BlockEndNodeType {
					end_node = end_node.PrevNamedSibling()
				}
				diagnostics = append(diagnostics,
					createDiagnostic("Unreachable code", start_node.StartPoint(), end_node.EndPoint(), lsp.WARNING))
			}
		}
	}
	d.diagnostics = append(d.diagnostics, diagnostics...)
}

// addInvalidExpressionErrors identifies and collects errors for invalid single expression statements in the given AST node.
// It uses a query executor to find expressions and checks if they are valid single expression statements.
// Diagnostics are created for each invalid expression statement found.
//
// Parameters:
//
//	node *ASTNode - The AST node to be checked for invalid expressions.
//	a *Analyzer - The analyzer used to get the parser and language information.
func (d *DiagnosticVisitor) addInvalidExpressionErrors(node *ASTNode, a *Analyzer) {
	var diagnostics []lsp.Diagnostic
	qe, err := NewQueryExecutor(_EXPRESSION_QUERY, node.Node, a.GetParser().language)
	if err != nil {
		log.Error().Err(err).Msg("Error creating query")
		return
	}
	for {
		match, ok := qe.NextMatch()
		if !ok {
			break
		}
		for _, capture := range match.Captures {
			node := capture.Node
			// check if the statement has an expression
			if node.Parent().Type() != StatementNodeType {
				continue
			}

			// check if the expression has no siblings
			if node.NextNamedSibling() != nil {
				continue
			}
			// check if the expression has only one child
			if node.ChildCount() > 1 {
				continue
			}

			// this is a single expression statement,
			// valid single expression statements are:
			// 1. core function
			// 2. assignment expression
			// 3. return statement
			// 4. call expression
			// 5. unary expression
			// 6. binary expression
			// 7. parenthesized expression

			if node.Child(0).Type() == CoreFunctionNodeType ||
				node.Child(0).Type() == AssignmentExpressionNodeType ||
				node.Child(0).Type() == ReturnNodeType ||
				node.Child(0).Type() == CallExpressionNodeType ||
				node.Child(0).Type() == UnaryExpressionNodeType ||
				node.Child(0).Type() == BinaryExpressionNodeType ||
				node.Child(0).Type() == ParenthesizedExpressionNodeType {
				continue
			}
			// Invalid single expression statement
			log.Debug().Str("node-type", node.Type()).Msg("invalid single expression statement found")
			diagnostics = append(diagnostics,
				createDiagnostic("Invalid statement", node.StartPoint(), node.EndPoint(), lsp.ERROR))

		}
	}
	d.diagnostics = append(d.diagnostics, diagnostics...)
}

// addInvalidAssignmentExpressionErrors analyzes the given AST node for invalid assignment expressions
// and adds corresponding diagnostic errors.
//
// Parameters:
// - node: The AST node to analyze.
// - a: The analyzer instance containing the parser and other analysis tools.
//
// This function ignores top-level assignments and only checks assignments within blocks.
func (d *DiagnosticVisitor) addInvalidAssignmentExpressionErrors(node *ASTNode, a *Analyzer) {
	var diagnostics []lsp.Diagnostic
	// These should only be for within the block, top level assignments are to be ignored here
	qe, err := NewQueryExecutor(_ASSINGMENT_EXPRESSION_QUERY, node.Node, a.GetParser().language)
	if err != nil {
		log.Error().Err(err).Msg("Error creating query")
		return
	}
	for {
		match, ok := qe.NextMatch()
		if !ok {
			break
		}
		for _, capture := range match.Captures {
			node := capture.Node
			// capture is statement -> expression -> assignment_expression
			n := node.NamedChild(0)
			n = n.NamedChild(0)
			if n == nil {
				diagnostics = append(diagnostics,
					createDiagnostic("Invalid assignment expression", node.StartPoint(), node.EndPoint(), lsp.ERROR))
				continue
			}

			if n.NamedChildCount() != 2 {
				diagnostics = append(diagnostics,
					createDiagnostic("Invalid assignment expression, left and right side are not valid", node.StartPoint(), node.EndPoint(), lsp.ERROR))
				continue
			}

			left := n.ChildByFieldName("left")
			if left.Type() != PseudoVariableNodeType && left.Type() != PseudoVariableExpressionNodeType {
				log.Debug().Str("left-hand-side", left.Type()).Msg("Invalid assignment")
				diagnostics = append(diagnostics,
					createDiagnostic("Invalid assignment: left-hand-side ", node.StartPoint(), node.EndPoint(), lsp.ERROR))
				continue
			}

			right := n.ChildByFieldName("right")
			if right.Type() != ExpressionNodeType || right.NamedChild(0).Type() == IdentifierNodeType {
				diagnostics = append(diagnostics,
					createDiagnostic("Invalid value on the right side of expression", node.StartPoint(), node.EndPoint(), lsp.ERROR))
				continue
			}
		}
	}
	d.diagnostics = append(d.diagnostics, diagnostics...)
}

// GetQueryDiagnostics collects various diagnostics for the given AST node.
// It checks for invalid expressions, deprecated comments, unreachable code, and syntax errors,
// and adds the corresponding diagnostics to the DiagnosticVisitor.
//
// Parameters:
//
//	node *ASTNode - The AST node to be checked for diagnostics.
//	a *Analyzer - The analyzer used to get the parser and language information.
func (d *DiagnosticVisitor) GetQueryDiagnostics(node *ASTNode, a *Analyzer) {
	// Since its not incremental, we can clear the diagnostics
	d.diagnostics = nil
	d.addInvalidExpressionErrors(node, a)
	d.addInvalidAssignmentExpressionErrors(node, a)
	d.addUnreachableCodeWarnings(node, a)
	// FIXME: fix false positives
	d.addSyntaxErrors(node, a)
	if settings.GlobalSettings.DeprecatedCommentHints {
		d.addDeprecatedCommentHints(node, a)
	}
}

// GetDiagnostics returns the collected diagnostics.
//
// Returns:
//
//	[]lsp.Diagnostic - A slice of diagnostics collected during the visit.
func (d *DiagnosticVisitor) GetDiagnostics() []lsp.Diagnostic {
	return d.diagnostics
}
