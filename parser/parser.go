package parser

import (
	"fmt"
	"reflect"
	"strconv"
)

// TODO: use logging instead of fmt.Println
type Parser struct {
	tokens []Token
	pos    int
	Root   *ASTNode
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
	// fmt.Println("Parsing block", p.peek().Type())
	if p.peek().Type() == EOF {
		return errorNode("Expected block but got EOF")
	}
	if p.consume(LBRACE) {
		node := repeat(p.parseStatement)()
		node.Type = BLOCK_NODE
		if node.Type == ERROR_NODE {
			return errorNode("Error parsing block")
		}
		if p.consume(RBRACE) {
			return node
		}
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
	// fmt.Println("Parsing request route")
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

	block := p.parseBlock()

	if block == nil || block.Type == ERROR_NODE {
		r.addChild(errorNode("Error parsing request_route block"))
		return r
	}

	r.addChild(block)
	return r
}

func (p *Parser) parseTopLevelStatement() *ASTNode {
	// fmt.Println("Parsing top level statement")
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

func (p *Parser) parseStatement() *ASTNode {
	// fmt.Println("Parsing statement")
	if p.peek().Type() == EOF {
		return errorNode("Expected identifier or core variable but got EOF")
	}
	// TODO: make this a choice
	node := seq(p.parseAssignment, p.parseEOS)()
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
	// fmt.Println("Parsing top level assignment")
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
	// fmt.Println("ParsingAssignment")
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
	// fmt.Println("Parsing assignment left side")
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
	// fmt.Println("Parsing core variable")
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
	// fmt.Println("Parsing assignment right side")
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
	// fmt.Println("parseIdentifier")
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

func (p *Parser) parseNumber() *ASTNode {
	// fmt.Println("parseNumber")
	if p.consume(INT) {
		value, ok := p.tokens[p.pos-1].Literal().(int)
		if !ok {
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
	// fmt.Println("parseString")
	if p.consume(STRING) {
		value, ok := p.tokens[p.pos-1].Literal().(string)
		if !ok {
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
	// fmt.Println("parseAssignmentOperator")
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
