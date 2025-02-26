package kamailio_cfg

import (
	sitter "github.com/smacker/go-tree-sitter"
)

const (
	ReturnNodeType                   = "return_statement"
	CoreFunctionNodeType             = "core_function_statement"
	StatementNodeType                = "statement"
	ParenthesizedExpressionNodeType  = "parenthesized_expression"
	AssignmentExpressionNodeType     = "assignment_expression"
	TopLevelAssignmentNodeType       = "top_level_assignment_expression"
	CompoundStatementNodeType        = "compound_statement"
	CallExpressionNodeType           = "call_expression"
	IdentifierNodeType               = "identifier"
	PseudoVariableNodeType           = "pseudo_variable"
	VARNodeType                      = "var_"
	AVPNodeType                      = "avp_var"
	DlgVarNodeType                   = "dlg_var"
	PseudoContentNodeType            = "pseudo_content"
	PseudoVariableExpressionNodeType = "pvar_expression"
	ExpressionNodeType               = "expression"
	BlockEndNodeType                 = "block_end"
	UnaryExpressionNodeType          = "unary_expression"
	BinaryExpressionNodeType         = "binary_expression"
	CaseStatementNodeType            = "case_statement"
	IFStatementNodeType              = "if_statement"
)

// UpdateTree updates the given parse tree by applying an edit operation.
// It modifies the tree to reflect changes between the leftNode and rightNode,
// using the provided formatted content.
//
// Parameters:
//
//	tree *sitter.Tree - The parse tree to be updated.
//	leftNode *sitter.Node - The starting node of the edit operation.
//	rightNode *sitter.Node - The ending node of the edit operation.
//	formattedContent string - The new content to replace the old content between leftNode and rightNode.
func UpdateTree(tree *sitter.Tree, leftNode, rightNode *sitter.Node, formattedContent string) {
	tree.Edit(sitter.EditInput{
		StartIndex:  leftNode.StartByte(),
		OldEndIndex: rightNode.EndByte(),
		NewEndIndex: uint32(len(formattedContent)),
		StartPoint:  leftNode.StartPoint(),
		OldEndPoint: rightNode.EndPoint(),
		NewEndPoint: sitter.Point{
			Row:    leftNode.StartPoint().Row,
			Column: uint32(len(formattedContent)),
		},
	})
}
