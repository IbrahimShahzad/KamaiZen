package parser

import (
	"errors"
	"fmt"
	"strings"
)

type ASTNodeType string

const (
	ASSIGNMENT_NODE   ASTNodeType = "Assignment"
	EMPTY_NODE        ASTNodeType = "Empty"
	ROOT_NODE         ASTNodeType = "Root"
	STRING_NODE       ASTNodeType = "String"
	EOF_NODE          ASTNodeType = "EOF"
	NUMBER_NODE       ASTNodeType = "Number"
	IDENTIFIER_NODE   ASTNodeType = "Identifier"
	CORE_VAR_VAR_NODE ASTNodeType = "var"
	CORE_VAR_AVP_NODE ASTNodeType = "avp"
	ERROR_NODE        ASTNodeType = "Error"
)

type ASTNode struct {
	Name  interface{}
	Type  ASTNodeType
	Value interface{}
}

func (n *ASTNode) String() string {
	// if value is of type []*ASTNode, do not print the value
	_, ok := n.Value.([]*ASTNode)
	if !ok {
		return fmt.Sprintf("ASTNode{Name: %v, Type: %v, Value: %v}", n.Name, n.Type, n.Value)
	}
	return fmt.Sprintf("ASTNode{Name: %v, Type: %v, Value: * }", n.Name, n.Type)
}

var EmptyNode = &ASTNode{Name: nil, Type: EMPTY_NODE, Value: nil}

func (n *ASTNode) addChild(child *ASTNode) error {
	if n.Value == nil {
		n.Value = []*ASTNode{child}
	} else {
		// assert n.Value is []*ASTNode
		if _, ok := n.Value.([]*ASTNode); !ok {
			return errors.New("Value is not a slice of ASTNode")
		}
		n.Value = append(n.Value.([]*ASTNode), child)
	}
	return nil
}

func (n *ASTNode) getChild(index int) (*ASTNode, error) {
	if n.Value == nil {
		return nil, errors.New("No children found")
	}
	// assert n.Value is []*ASTNode
	if _, ok := n.Value.([]*ASTNode); !ok {
		return nil, errors.New("Value is not a slice of ASTNode")
	}
	if index < 0 || index >= len(n.Value.([]*ASTNode)) {
		return nil, errors.New("Index out of range")
	}
	return n.Value.([]*ASTNode)[index], nil
}

func (n *ASTNode) getChildByName(name string) (*ASTNode, error) {
	if n.Value == nil {
		return nil, errors.New("No children found")
	}
	// assert n.Value is map[string]*ASTNode
	if _, ok := n.Value.(map[string]*ASTNode); !ok {
		return nil, errors.New("Value is not a map of string to ASTNode")
	}
	if _, ok := n.Value.(map[string]*ASTNode)[name]; !ok {
		return nil, errors.New("Name not found")
	}
	return n.Value.(map[string]*ASTNode)[name], nil
}

// Combinators
func seq(parsers ...func() *ASTNode) func() *[]ASTNode {
	return func() *[]ASTNode {
		var nodes *[]ASTNode
		for _, parser := range parsers {
			node := parser()
			if node == nil {
				return nil // Fail if any parser in the sequence fails
			}
			if nodes == nil {
				nodes = &[]ASTNode{*node}
			} else {
				// Adding sibling nodes
				*nodes = append(*nodes, *node)
			}
		}
		return nodes
	}
}

func choice(parsers ...func() *ASTNode) func() *ASTNode {
	return func() *ASTNode {
		for _, parser := range parsers {
			node := parser()
			if node != nil {
				return node // Return the first successful parser
			}
		}
		return nil
	}
}

func optional(parser func() *ASTNode) func() *ASTNode {
	return func() *ASTNode {
		node := parser()
		if node == nil {
			return EmptyNode
		}
		return node
	}
}

type Parser struct {
	tokens []Token
	pos    int
	Root   *ASTNode
}

func NewParser(tokens []Token) *Parser {
	return &Parser{
		tokens: tokens,
		pos:    0,
	}
}

func (p *Parser) Parse() *ASTNode {
	r := &ASTNode{
		Name: nil,
		Type: ROOT_NODE,
	}
	sibs := seq(
		p.parseAssignment,
		p.parseEOF,
	)()
	for _, sib := range *sibs {
		err := r.addChild(&sib)
		if err != nil {
			fmt.Println("Error adding sib")
		}
	}
	p.Root = r
	p.Root.Print(0)
	return p.Root
}

func (p *Parser) parseEOF() *ASTNode {
	if p.tokens[p.pos].Type() == EOF {
		p.pos++
		return &ASTNode{Name: nil, Type: EOF_NODE, Value: nil}
	}
	return nil
}

func (p *Parser) parseAssignment() *ASTNode {

	left := p.assingnment_left_side()
	if left == nil {
		fmt.Println("left is nil")
		return nil
	}
	x := p.parseOperator()
	if x == nil {
		fmt.Println("x is nil")
	}
	right := p.assignment_right_side()
	if right == nil {
		fmt.Println("right is nil")
		return nil
	}
	// create a new node
	node := &ASTNode{
		Name:  nil,
		Type:  ASSIGNMENT_NODE,
		Value: nil,
	}
	err := node.addChild(left)
	if err != nil {
		fmt.Println("Error adding left")
	}

	err = node.addChild(right)
	if err != nil {
		fmt.Println("Error adding right")
	}
	return node
}

func (p *Parser) assingnment_left_side() *ASTNode {
	l := choice(
		p.parseIdentifier,
		p.ParseCoreVariable,
	)()
	l.Name = "left"
	return l
}

func (p *Parser) ParseCoreVariable() *ASTNode {
	if p.tokens[p.pos].Type() == CORE_VARIABLE {
		coreVariableType := p.tokens[p.pos].(*CoreVariableToken).VariableType
		cn := CORE_VAR_VAR_NODE
		switch coreVariableType {
		case "var":
			cn = CORE_VAR_VAR_NODE
		case "avp":
			cn = CORE_VAR_AVP_NODE
		}
		node := &ASTNode{
			Name: nil,
			Type: cn,
			Value: &ASTNode{
				Name:  nil,
				Type:  IDENTIFIER_NODE,
				Value: p.tokens[p.pos].(*CoreVariableToken).Literal,
			},
		}
		p.pos++
		fmt.Printf("ParseCoreVariable: %v\n", node)
		return node
	}
	return nil
}

func (p *Parser) assignment_right_side() *ASTNode {
	r := choice(
		p.parseNumber,
		p.parseString,
	)()
	r.Name = "right"
	return r
}

func (p *Parser) parseIdentifier() *ASTNode {
	if p.tokens[p.pos].Type() == IDENT {
		node := &ASTNode{
			Name:  nil,
			Type:  IDENTIFIER_NODE,
			Value: p.tokens[p.pos].Literal().(string),
		}
		p.pos++
		return node
	}
	return nil
}

func (p *Parser) parseNumber() *ASTNode {
	if p.tokens[p.pos].Type() == INT {
		node := &ASTNode{
			Name:  nil,
			Type:  NUMBER_NODE,
			Value: p.tokens[p.pos].Literal().(int),
		}
		p.pos++
		return node
	}
	return nil
}

func (p *Parser) parseString() *ASTNode {
	if p.tokens[p.pos].Type() == STRING {
		node := &ASTNode{
			Name:  nil,
			Type:  STRING_NODE,
			Value: p.tokens[p.pos].Literal().(string),
		}
		p.pos++
		return node
	}
	return nil
}

func (p *Parser) parseOperator() *ASTNode {
	if p.tokens[p.pos].Type() == ASSIGN {
		node := &ASTNode{
			Name:  nil,
			Type:  ASSIGNMENT_NODE,
			Value: "=",
		}
		p.pos++
		return node
	}
	return nil
}

func (n *ASTNode) Equals(other *ASTNode) bool {

	if (n.Name != nil && other.Name != nil) && n.Name != other.Name {
		return false
	}

	if n.Type != other.Type {
		return false
	}

	if n.Value == nil && other.Value == nil {
		return true
	}

	nSlice, nOk := n.Value.([]*ASTNode)
	otherSlice, otherOk := other.Value.([]*ASTNode)
	if nOk && otherOk {
		// Compare lengths first to avoid out-of-bounds errors
		if len(nSlice) != len(otherSlice) {
			return false
		}
		// Compare each child node
		for i, child := range nSlice {
			if !child.Equals(otherSlice[i]) {
				return false
			}
		}
	} else if n.Value != other.Value {
		// Directly compare values if they are not slices
		return false
	}

	return true
}
func (n *ASTNode) Print(ident int) {
	if n == nil {
		fmt.Println("nil")
		return
	}

	fmt.Printf("%s%v\n", strings.Repeat("  ", ident), n)
	if n.Value != nil {
		if _, ok := n.Value.([]*ASTNode); ok {
			for _, child := range n.Value.([]*ASTNode) {
				child.Print(ident + 1)
			}
		}
	}

}
