package kamailio_cfg

import (
	sitter "github.com/smacker/go-tree-sitter"
)

const (
	_ERROR_QUERY                 = "(ERROR) @error"
	_XML_QUERY                   = "(xml) @xml"
	_DEPRECATED_COMMENT_QUERY    = "(deprecated_comment) @deprecated"
	_CORE_FUNCTION_QUERY         = "(core_function) @core_statement"
	_RETURN_STATEMENTS_QUERY     = "(return_statement) @return"
	_FUNCTION_QUERY              = "(function: (expression)) @function"
	_STATEMENT_QUERY             = "(statement) @parent_statement"
	_EXPRESSION_QUERY            = "(expression) @expression_statement"
	_ASSINGMENT_QUERY            = "(assignment_expression) @assignment_expression"
	_ASSINGMENT_EXPRESSION_QUERY = "(statement (expression (assignment_expression))) @assignment_expression"
)

// QueryExecutor is a struct that encapsulates the execution of tree-sitter queries.
// It holds references to the query cursor, the query itself, and the AST node being queried.
type QueryExecutor struct {
	cursor *sitter.QueryCursor
	query  *sitter.Query
	node   *sitter.Node
}

// NewQueryExecutor creates and returns a new instance of QueryExecutor.
// It initializes the query and query cursor, and executes the query on the given AST node.
//
// Parameters:
//
//	queryStr string - The query string to be executed.
//	node *sitter.Node - The AST node to be queried.
//	language *sitter.Language - The language used for parsing.
//
// Returns:
//
//	*QueryExecutor - A new instance of QueryExecutor.
//	error - An error if the query creation fails, otherwise nil.
func NewQueryExecutor(queryStr string, node *sitter.Node, language *sitter.Language) (*QueryExecutor, error) {
	q, err := sitter.NewQuery([]byte(queryStr), language)
	if err != nil {
		return nil, err
	}
	cursor := sitter.NewQueryCursor()
	cursor.Exec(q, node)
	return &QueryExecutor{
		cursor: cursor,
		query:  q,
		node:   node,
	}, nil
}

// NextMatch returns the next match for the query executed by the QueryExecutor.
// It advances the query cursor to the next match and returns the match and a boolean indicating success.
//
// Returns:
//
//	*sitter.QueryMatch - The next query match.
//	bool - True if a match is found, false otherwise.
func (qe *QueryExecutor) NextMatch() (*sitter.QueryMatch, bool) {
	return qe.cursor.NextMatch()
}
