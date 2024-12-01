package parser

import (
	"fmt"
	"log"
	"reflect"
	"strconv"
)

// TODO: use logging instead of fmt.Println
type Parser struct {
	tokens           []Token
	pos              int
	Root             *ASTNode
	tranasctionStack []int // Stack to store transaction checkpoints
}

func (p *Parser) startTransaction() {
	p.tranasctionStack = append(p.tranasctionStack, p.pos)
}

func (p *Parser) commitTransaction() {
	if len(p.tranasctionStack) > 0 {
		p.tranasctionStack = p.tranasctionStack[:len(p.tranasctionStack)-1] // pop
	}
}

func (p *Parser) rollbackTransaction() {
	if len(p.tranasctionStack) > 0 {
		p.pos = p.tranasctionStack[len(p.tranasctionStack)-1]               // restore position
		p.tranasctionStack = p.tranasctionStack[:len(p.tranasctionStack)-1] // pop
	}
}

func (p *Parser) rollBackWithError(msg string) {
	log.Printf("Rolling back with error: %s", msg)
	p.rollbackTransaction()
}

func (p *Parser) withTransaction(parser func() *ASTNode) *ASTNode {
	p.startTransaction()

	node := parser()
	if node == nil || node.Type == EMPTY_NODE {
		p.rollbackTransaction()
		return errorNode("Transaction Failed")
	}

	if node.Type == ERROR_NODE {
		p.rollBackWithError(node.Value.(string))
		return node
	}

	p.commitTransaction()
	return node
}

func (p *Parser) updateASTLevel() {
	p.Root.level = 0
	p.updateASTLevelRecursive(p.Root)
}

func (p *Parser) updateASTLevelRecursive(node *ASTNode) {
	if node.Value == nil {
		return
	}
	if _, ok := node.Value.([]*ASTNode); !ok {
		// it is a leaf node
		return
	}
	for _, child := range node.Value.([]*ASTNode) {
		child.level = node.level + 1
		p.updateASTLevelRecursive(child)
	}
}

func NewParser(tokens []Token) *Parser {
	return &Parser{
		tokens: tokens,
		pos:    0,
		Root:   nil,
	}
}

func (p *Parser) peek() Token {
	if p.pos >= len(p.tokens) {
		return nil
	}
	return p.tokens[p.pos]
}

func (p *Parser) peekNext() Token {
	if p.pos+1 >= len(p.tokens) {
		return nil
	}
	return p.tokens[p.pos+1]
}

func (p *Parser) next() Token {
	if p.pos >= len(p.tokens) {
		return nil
	}
	p.pos++
	return p.tokens[p.pos-1]
}

func (p *Parser) consume(t TokenType) bool {
	if p.peek().Type() == t {
		p.pos++
		return true
	}
	return false
}

func (p *Parser) consumeTokenWithLiteral(t TokenType, literal string) bool {
	if p.peek().Type() == t && p.peek().Literal().(string) == literal {
		p.pos++
		return true
	}
	return false
}

func (p *Parser) unConsume() {
	if p.pos > 0 {
		p.pos--
	}
}

func (p *Parser) parseEOS() *ASTNode {
	// fmt.Println("Parsing EOS")
	if p.peek().Type() == EOF {
		return errorNode("Expected semicolon but got EOF")
	}
	if p.consume(SEMICOLON) {
		return &ASTNode{Name: nil, Type: EOS_NODE, Value: nil, level: 0}
	}
	return errorNode("Expected semicolon but got " + p.tokens[p.pos].Literal().(string))
}

func (p *Parser) parseBlock() *ASTNode {
	fmt.Println("Parsing block ", p.peek().Type(), p.peek().Literal())
	if p.peek().Type() == EOF {
		return errorNode("Expected block but got EOF")
	}
	if p.consume(LBRACE) {
		node := repeat(p.parseStatement)()
		fmt.Println("Block items: ", node)
		node.Type = BLOCK_NODE
		if node.Type == ERROR_NODE {
			p.unConsume()
			return errorNode("Error parsing block")
		}
		if p.consume(RBRACE) {
			return node
		}
		p.unConsume()
		return errorNode("Expected closing brace but got " + p.tokens[p.pos].Literal().(string))
	}
	return errorNode("Expected opening brace but got " + p.tokens[p.pos].Literal().(string))
}

func (p *Parser) parseFileStarter() *ASTNode {
	// fmt.Println("Parsing file starter")
	if p.peek().Type() == EOF {
		return errorNode("Expected file starter but got EOF")
	}
	if p.consume(PREPROC) {
		starter, ok := p.tokens[p.pos-1].Literal().(string)
		if !ok {
			p.unConsume()
			return errorNode("Invalid file starter " + p.tokens[p.pos].Literal().(string))
		}
		switch starter {
		case "SER":
			fallthrough
		case "KAMAILIO":
			fallthrough
		case "OPENSER":
			fallthrough
		case "MAXCOMPAT":
			fallthrough
		case "ALL":
			return &ASTNode{
				Name:  nil,
				Type:  FILE_STARTER_NODE,
				Value: starter,
			}
		default:
			p.unConsume()
			return errorNode("Invalid file starter " + starter)
		}

	}
	return errorNode("Expected file starter but got " + p.tokens[p.pos].Literal().(string))
}

func (p *Parser) parseRequestRouteKeyword() *ASTNode {
	// fmt.Println("Parsing request route keyword")
	if p.peek().Type() == EOF {
		return errorNode("Expected request_route keyword but got EOF")
	}
	if p.consume(ROUTE) {
		val := p.tokens[p.pos-1].Literal().(string) // consumed token
		var route_type ASTNodeType
		switch val {
		case "request_route":
			route_type = REQUEST_ROUTE_NODE
		case "reply_route":
			route_type = REPLY_ROUTE_NODE
		case "failure_route":
			route_type = FAILURE_ROUTE_NODE
		case "onreply_route":
			route_type = ONREPLY_ROUTE_NODE
		case "branch_route":
			route_type = BRANCH_ROUTE_NODE
		case "local_route":
			route_type = LOCAL_ROUTE_NODE
		case "startup_route":
			route_type = STARTUP_ROUTE_NODE
		case "route":
			route_type = ROUTE_NODE
		default:
			p.unConsume()
			errorNode("Invalid route type " + val)
		}
		return &ASTNode{
			Name:  nil,
			Type:  route_type,
			Value: nil,
		}

	}
	return errorNode("Expected request_route keyword but got " + p.tokens[p.pos].Literal().(string))
}

func (p *Parser) parseRequestRouteName() *ASTNode {
	// fmt.Println("Parsing request route name")
	if p.peek().Type() == EOF {
		return errorNode("Expected request_route name but got EOF")
	}
	// name should be  '[', identifier or number or string, ']'
	if p.consume(LBRACKET) {
		name := choice(
			p.parseIdentifier,
			p.parseNumber,
			p.parseString,
		)()
		if name == nil {
			return errorNode("Error parsing request_route name")
		}
		if p.consume(RBRACKET) {
			return name
		}
		return errorNode("Expected closing bracket but got " + p.tokens[p.pos].Literal().(string))
	}
	return errorNode("Expected opening bracket but got " + p.tokens[p.pos].Literal().(string))
}

func (p *Parser) parseRequestRoute() *ASTNode {
	// the format is as follows
	// request_route [optional_name] Block
	fmt.Println("Parsing request route", p.peek().Type(), p.peek().Literal())
	if p.peek().Type() == EOF {
		return errorNode("Expected request_route but got EOF")
	}
	r := p.parseRequestRouteKeyword()
	if r == nil ||
		r.Type == ERROR_NODE ||
		r.Type == EMPTY_NODE {
		return errorNode("Error parsing request_route")
	}

	// By this point we have the route type
	name := optional(p.parseRequestRouteName)()

	if name != nil && name.Type != EMPTY_NODE {
		var n interface{}
		switch reflect.TypeOf(name.Value).Kind() {
		case reflect.String:
			n = name.Value.(string)
		case reflect.Int:
			n = strconv.Itoa(name.Value.(int))
		default:
			n = nil
		}
		r.Name = n
	}

	if r.Type == ROUTE_NODE && r.Name == nil {
		r.addChild(errorNode("Route name is required"))
	}

	fmt.Println("Request route node", r, "parsing block")
	block := p.parseBlock()

	if block == nil || block.Type == ERROR_NODE {
		r.addChild(errorNode("Error parsing block: " + block.Value.(string)))
		return r
	}

	r.addChild(block)
	return r
}

func (p *Parser) parseTopLevelStatement() *ASTNode {
	fmt.Println("Parsing top level statement", p.peek().Type(), p.peek().Literal())
	if p.peek().Type() == EOF {
		return errorNode("Expected top level statement but got EOF")
	}
	child := choice(
		p.parseFileStarter,
		p.parseTopLevelAssignment,
		p.parseRequestRoute,
		// p.parseBlock,
		p.parseEOF,
	)()
	if child == nil {
		return errorNode("Error parsing top level statement")
	}
	node := &ASTNode{
		Name:   nil,
		Type:   TOP_LEVEL_STATEMENT_NODE,
		Value:  nil,
		Parent: nil,
		level:  0,
	}
	err := node.addChild(child)
	if err != nil {
		fmt.Println("Error adding children to statement node")
	}
	return node
}

func getPrecedence(op TokenType) int {
	switch op {
	case ADD_OP, SUB_OP:
		return 1
	case MUL_OP, DIV_OP:
		return 2
	case OR_OP:
		return 3
	case AND_OP:
		return 4
	case BIN_OR_OP:
		return 5
	case BIN_XOR_OP:
		return 6
	case BIN_AND_OP:
		return 7
	case EQ_OP, REGEX_OP, NE_OP:
		return 8
	case GT_OP, GE_OP, LT_OP, LE_OP:
		return 9
	// TODO: check if these operators are required
	// case OPERATORS.BIN_LSHIFT, OPERATORS.BIN_RSHIFT:
	// 	return 10
	default:
		return 0
	}
}

func (p *Parser) parseBinaryOperator(precedence int) *ASTNode {
	if p.peek().Type() == EOF {
		return errorNode("Expected binary operator but got EOF")
	}
	if !IsBinaryOperator(p.peek()) {
		return errorNode("Expected binary operator but got " + p.tokens[p.pos].Literal().(string))
	}
	op := p.tokens[p.pos]
	p.pos++
	node := &ASTNode{
		Name:  nil,
		Type:  OPERATOR_NODE,
		Value: op.Literal().(string),
	}
	return node
}

func (p *Parser) parseBinaryExpression() *ASTNode {
	// binary expression is of the form
	// expression operator expression
	if p.peek().Type() == EOF || p.peekNext().Type() == EOF {
		return errorNode("Expected binary expression but got EOF")
	}
	// operators with precedence
	node := &ASTNode{
		Name:  nil,
		Type:  BINARY_EXPR_NODE,
		Value: nil,
	}

	// operator
	if IsBinaryOperator(p.peekNext()) {
		left := choice(
			p.parseNonBinaryExpression,
		)()
		if left == nil {
			return errorNode("Error parsing left expression")
		}
		left.Name = "left"
		node.addChild(left)
		fmt.Println("Left node", left)

		op := p.parseBinaryOperator(getPrecedence(p.peekNext().Type()))
		if op == nil {
			return errorNode("Error parsing operator")
		}
		op.Name = "operator"
		node.addChild(op)
		fmt.Println("Operator node", op)

		right := p.parseNonBinaryExpression()
		if right == nil {
			return errorNode("Error parsing right expression")
		}
		right.Name = "right"
		node.addChild(right)
		fmt.Println("Right node", right)

		return node

	}

	fmt.Println("[B1] Error parsing binary expression")
	return errorNode("Expected binary expression but got " + p.tokens[p.pos].Literal().(string))
}

func (p *Parser) parseExpression() *ASTNode {
	// fmt.Println("Parsing expression")
	if p.peek().Type() == EOF {
		return errorNode("Expected expression but got EOF")
	}
	return choice(
		// p.parseParenthesizedExpression,
		p.parseBinaryExpression, // this could be (expression operator expression)
		p.parseNonBinaryExpression,
	)()
}

func (p *Parser) parseNonBinaryExpression() *ASTNode {
	fmt.Println("Parsing non binary expression", p.peek().Type(), p.peek().Literal())
	if p.peek().Type() == EOF {
		return errorNode("Expected non binary expression but got EOF")
	}
	// $.assignment_expression,
	// $.pseudo_variable,
	// $.pvar_expression,
	// $.unary_expression,
	// $.cast_expression,
	// $.subscript_expression,
	// $.call_expression,
	// $.field_expression,
	// $.select_param,
	// $.identifier,
	// $.number_literal,
	// $.string,
	// $.true,
	// $.false,
	// $.null,
	// $.parenthesized_expression,
	return choice(
		// p.parseKeyword,
		p.parseNumber,
		p.parseString,
		p.parseIdentifier,
		p.ParseCoreVariable,
		p.parseParenthesizedExpression,
	)()
}

func (p *Parser) commaSeparatedExpressions() *ASTNode {
	// fmt.Println("Parsing comma separated expressions")
	if p.peek().Type() == EOF {
		return errorNode("Expected comma separated expressions but got EOF")
	}
	node := &ASTNode{
		Name:  nil,
		Type:  COMMA_SEPARATED_EXPR_NODE,
		Value: nil,
	}
	left := p.parseExpression()
	left.Name = "left"
	if left == nil {
		return errorNode("Error parsing expression")
	}
	if p.consume(COMMA) {
		// right is a choice between expression and comma separated expressions
		right := choice(
			p.parseExpression,
			p.commaSeparatedExpressions,
		)()
		if right == nil {
			return errorNode("Error parsing comma separated expressions")
		}
		right.Name = "right"

		node.addChild(left)
		node.addChild(right)
		return node
	}
	// return error
	return errorNode("Expected comma but got " + p.tokens[p.pos].Literal().(string))
}

func (p *Parser) parseParenthesizedExpression() *ASTNode {
	if p.peek().Type() == EOF {
		return errorNode("Expected parenthesized expression but got EOF")
	}
	if p.consume(LPAREN) {
		// find the closing parenthesis and parse the expression inside
		node := choice(
			p.parseExpression,
			p.commaSeparatedExpressions,
		)()
		if p.consume(RPAREN) {
			return node
		}
		return errorNode("Expected closing parenthesis but got " + p.tokens[p.pos].Literal().(string))
	}
	return errorNode("Expected opening parenthesis but got " + p.tokens[p.pos].Literal().(string))
}

func (p *Parser) parseIfStatement() *ASTNode {
	if p.peek().Type() == EOF {
		return errorNode("Expected if statement but got EOF")
	}

	if p.consumeTokenWithLiteral(KEYWORD, "if") {
		node := ASTNode{
			Name:  "if",
			Type:  IF_NODE,
			Value: nil,
		}
		condition := p.parseParenthesizedExpression()
		if condition == nil {
			p.unConsume()
			return errorNode("Error parsing if statement")
		}
		condition.Name = "condition"
		node.addChild(condition)
		consequence := p.parseBlock()
		if consequence == nil {
			p.unConsume()
			return errorNode("Error parsing if statement")
		}
		consequence.Name = "consequence"
		node.addChild(consequence)

		// optionals return empty node if not found
		elseBlock := optional(p.parseElseBlock)()
		if elseBlock.Type != EMPTY_NODE {
			// only add else block if it is not empty
			elseBlock.Name = "else"
			node.addChild(elseBlock)
		}
		return &node
	}
	return errorNode("Expected if statement but got " + p.tokens[p.pos].Literal().(string))
}

func (p *Parser) parseElseBlock() *ASTNode {
	if p.peek().Type() == EOF {
		return errorNode("Expected else block but got EOF")
	}
	if p.consumeTokenWithLiteral(KEYWORD, "else") {
		block := p.parseBlock()
		if block == nil {
			p.unConsume()
			return errorNode("Error parsing else block")
		}
		block.Name = "block"
		return block
	}
	return errorNode("Expected else block but got " + p.tokens[p.pos].Literal().(string))
}

func (p *Parser) parseAssignmentStatement() *ASTNode {
	if p.peek().Type() == EOF {
		return errorNode("Expected assignment statement but got EOF")
	}
	return seq(
		p.parseAssignment,
		p.parseEOS,
	)()
}

func (p *Parser) parseKeywordStatement() *ASTNode {
	if p.peek().Type() == EOF {
		return errorNode("Expected keyword statement but got EOF")
	}
	current_pos := p.pos
	node := seq(
		p.parseKeyword,
		p.parseEOS,
	)()
	if node == nil {
		p.pos = current_pos
		return errorNode("Error parsing keyword statement")
	}
	node.Type = KEYWORD_NODE
	return node
}

func (p *Parser) parseExpressionStatement() *ASTNode {
	if p.peek().Type() == EOF {
		return errorNode("Expected expression statement but got EOF")
	}
	return seq(
		p.parseExpression,
		p.parseEOS,
	)()
}

func (p *Parser) parseStatement() *ASTNode {
	if p.peek().Type() == EOF {
		return errorNode("Expected identifier or core variable but got EOF")
	}
	node := choice(
		p.parseKeywordStatement,
		p.parseIfStatement,
		p.parseAssignmentStatement,
		p.parseExpressionStatement,
	)()
	if node == nil {
		return errorNode("Error parsing statement")
	}
	node.Type = STATEMENT_NODE
	return node
}

func (p *Parser) setAsRoot(node *ASTNode) {
	p.Root = node
	p.Root.Parent = nil
	p.Root.Type = ROOT_NODE
	p.updateASTLevel()
}

func (p *Parser) parseEOF() *ASTNode {
	if p.consume(EOF) {
		return &ASTNode{Name: nil, Type: EOF_NODE, Value: nil, level: 0}
	}
	return errorNode("Expected EOF but got " + p.tokens[p.pos].Literal().(string))
}

func (p *Parser) parseTopLevelAssignment() *ASTNode {
	if p.peek().Type() == EOF {
		return errorNode("Expected identifier or core variable but got EOF")
	}
	return parseWithSeq(
		ASSIGNMENT_NODE,
		p.assingnment_left_side,
		p.parseAssignOperator,
		p.assignment_right_side)
}

func (p *Parser) parseAssignment() *ASTNode {
	if p.peek().Type() == EOF {
		return errorNode("Expected identifier or core variable but got EOF")
	}
	return parseWithSeq(
		ASSIGNMENT_NODE,
		p.assingnment_left_side,
		p.parseAssignOperator,
		p.assignment_right_side)
}

func (p *Parser) assingnment_left_side() *ASTNode {
	if p.peek().Type() == EOF {
		return errorNode("Expected identifier or core variable but got EOF")
	}
	left := choice(
		p.parseIdentifier,
		p.ParseCoreVariable,
	)()
	if left == nil {
		return errorNode("Expected identifier or core variable but got " + p.tokens[p.pos].Literal().(string))
	}
	left.Name = "left"
	return left
}

func (p *Parser) ParseCoreVariable() *ASTNode {
	if p.peek().Type() == EOF {
		return errorNode("Expected core variable but got EOF")
	}
	if p.consume(CORE_VARIABLE) {
		coreVariableType := p.tokens[p.pos-1].(*CoreVariableToken).VariableType
		cn := CORE_VAR_VAR_NODE
		switch coreVariableType {
		case "var":
			cn = CORE_VAR_VAR_NODE
		case "avp":
			cn = CORE_VAR_AVP_NODE
		}
		node := &ASTNode{
			Name:  nil,
			Type:  cn,
			Value: nil,
		}
		child := &ASTNode{
			Name:  nil,
			Type:  IDENTIFIER_NODE,
			Value: p.tokens[p.pos-1].(*CoreVariableToken).VariableName,
		}
		node.addChild(child)
		return node
	}
	return errorNode("Expected core variable but got " + p.tokens[p.pos].Literal().(string))
}

func (p *Parser) assignment_right_side() *ASTNode {
	if p.peek().Type() == EOF {
		return errorNode("Expected number, string or core variable but got EOF")
	}
	right := choice(
		p.parseNumber,
		p.parseString,
		p.ParseCoreVariable,
	)()
	if right == nil {
		return errorNode("Expected number, string or core variable but got " + p.tokens[p.pos].Literal().(string))
	}
	right.Name = "right"
	return right
}

func (p *Parser) parseIdentifier() *ASTNode {
	if p.tokens[p.pos].Type() == EOF {
		return errorNode("Expected identifier but got EOF")
	}

	if p.consume(IDENT) {
		value, ok := p.tokens[p.pos-1].Literal().(string)
		if !ok {
			return errorNode("Invalid identifier " + p.tokens[p.pos].Literal().(string))
		}
		node := &ASTNode{
			Name:  nil,
			Type:  IDENTIFIER_NODE,
			Value: value,
		}
		return node
	}
	return errorNode("Expected identifier")
}

func (p *Parser) parseKeyword() *ASTNode {
	if p.tokens[p.pos].Type() == EOF {
		return errorNode("Expected keyword but got EOF")
	}
	if p.consume(KEYWORD) {
		value, ok := p.tokens[p.pos-1].Literal().(string)
		if !ok {
			p.unConsume()
			return errorNode("Invalid keyword " + p.tokens[p.pos].Literal().(string))
		}
		node := &ASTNode{
			Name:  nil,
			Type:  KEYWORD_NODE,
			Value: value,
		}
		return node
	}
	return errorNode("Expected keyword but got " + p.tokens[p.pos].Literal().(string))
}

func (p *Parser) parseNumber() *ASTNode {
	if p.consume(INT) {
		value, ok := p.tokens[p.pos-1].Literal().(int)
		if !ok {
			p.unConsume()
			return errorNode("Invalid number " + p.tokens[p.pos].Literal().(string))
		}
		node := &ASTNode{
			Name:  nil,
			Type:  NUMBER_NODE,
			Value: value,
		}
		return node
	}
	if p.tokens[p.pos].Type() == EOF {
		return errorNode("Expected number but got EOF")
	}
	return errorNode("Expected number but got " + p.tokens[p.pos].Literal().(string))
}

func (p *Parser) parseString() *ASTNode {
	if p.consume(STRING) {
		value, ok := p.tokens[p.pos-1].Literal().(string)
		if !ok {
			p.unConsume()
			return errorNode("Invalid string " + p.tokens[p.pos].Literal().(string))
		}
		node := &ASTNode{
			Name:  nil,
			Type:  STRING_NODE,
			Value: value,
		}
		return node
	}
	if p.tokens[p.pos].Type() == EOF {
		return errorNode("Expected string but got EOF")
	}
	return errorNode("Expected string but got " + p.tokens[p.pos].Literal().(string))
}

func (p *Parser) parseAssignOperator() *ASTNode {
	if p.consume(ASSIGN) {
		node := &ASTNode{
			Name:  nil,
			Type:  OPERATOR_NODE,
			Value: "=",
		}
		return node
	}
	if p.tokens[p.pos].Type() == EOF {
		return errorNode("Expected = but got EOF")
	}
	return errorNode("Expected assignment operator but got " + p.tokens[p.pos].Literal().(string))
}

func (p *Parser) parseOperator() *ASTNode {
	return choice(
		p.parseAssignOperator,
	)()
}

func parseWithSeq(nodeType ASTNodeType, parsers ...func() *ASTNode) *ASTNode {
	node := seq(parsers...)()
	if node == nil {
		return errorNode(fmt.Sprintf("Error parsing %s", nodeType))
	}
	node.Type = nodeType
	return node
}

func (p *Parser) Parse() *ASTNode {
	fmt.Printf("-----TOKENS------------\n")
	for _, token := range p.tokens {
		fmt.Printf("%v\n", token)
	}
	node := repeat(p.parseTopLevelStatement)()
	eof := p.parseEOF()
	if err := node.addChild(eof); err != nil {
		fmt.Println("Error adding EOF to root node")
	}
	p.setAsRoot(node)
	fmt.Printf("-----Parse Output------\n%v\n", p.Root)
	return p.Root
}
