package parser

import (
	"errors"
	"fmt"
	"strings"
)

type ASTNodeType string

const (
	ASSIGNMENT_NODE          ASTNodeType = "Assignment"
	EMPTY_NODE               ASTNodeType = "Empty"
	ROOT_NODE                ASTNodeType = "Root"
	STRING_NODE              ASTNodeType = "String"
	EOF_NODE                 ASTNodeType = "EOF"
	NUMBER_NODE              ASTNodeType = "Number"
	IDENTIFIER_NODE          ASTNodeType = "Identifier"
	CORE_VAR_VAR_NODE        ASTNodeType = "var"
	CORE_VAR_AVP_NODE        ASTNodeType = "avp"
	ERROR_NODE               ASTNodeType = "Error"
	EOS_NODE                 ASTNodeType = "EOS"
	STATEMENT_NODE           ASTNodeType = "Statement"
	BLOCK_NODE               ASTNodeType = "Block"
	TOP_LEVEL_STATEMENT_NODE ASTNodeType = "TopLevelStatement"
	OPERATOR_NODE            ASTNodeType = "Operator"
	FILE_STARTER_NODE        ASTNodeType = "FileStarter"
)

type ASTNode struct {
	Name  interface{}
	Type  ASTNodeType
	Value interface{}
	level int
	// pointer to parent node
	Parent *ASTNode
}

var EmptyNode = &ASTNode{Name: nil, Type: EMPTY_NODE, Value: nil}

func errorNode(msg string) *ASTNode {
	return &ASTNode{Name: "ERROR", Type: ERROR_NODE, Value: msg}
}

func (n *ASTNode) IsLeaf() bool {
	if n.Value == nil {
		return true
	}
	_, ok := n.Value.([]*ASTNode)
	return !ok
}

func (n *ASTNode) getName() string {
	if n.Name == nil {
		return ""
	}
	return n.Name.(string) + ": "
}

func (n *ASTNode) GetParent() *ASTNode {
	return n.Parent
}

func (n *ASTNode) GetRoot() *ASTNode {
	if n.Parent == nil {
		return n
	}
	return n.Parent.GetRoot()
}

func (n *ASTNode) debugPrintAll() {
	fmt.Printf("Name: %v\nType: %v\nValue: %v\n Level: %v\n Parent: %v\n", n.Name, n.Type, n.Value, n.level, n.Parent)
}

func (n *ASTNode) String() string {
	// example output:
	// (ASSIGNMENT_NODE
	//   left: (IDENTIFIER_NODE x)
	//   (OPERATOR_NODE =)
	//   right:(STRING_NODE string))
	if n.Value == nil {
		return fmt.Sprintf("%s%s(%s)", strings.Repeat("  ", n.level), n.getName(), n.Type)
	}
	if _, ok := n.Value.([]*ASTNode); !ok {
		return fmt.Sprintf("%s%s(%s %v)", strings.Repeat("  ", n.level), n.getName(), n.Type, n.Value)
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%s%s(%s", strings.Repeat("  ", n.level), n.getName(), n.Type))

	for _, child := range n.Value.([]*ASTNode) {
		sb.WriteString(fmt.Sprintf("\n%s%s", strings.Repeat("  ", n.level), child.String()))
	}
	sb.WriteString(")")
	return sb.String()

}

func (n *ASTNode) addChild(child *ASTNode) error {
	// Set the parent node
	child.Parent = n
	child.level = n.level + 1
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
