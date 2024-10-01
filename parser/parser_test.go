package parser_test

import (
	"KamaiZen/parser"
	"testing"
)

func TestASTNode_isEqual(t *testing.T) {
	tests := []struct {
		input1   *parser.ASTNode
		input2   *parser.ASTNode
		expected bool
	}{
		{
			input1: &parser.ASTNode{
				Name:  nil,
				Type:  parser.IDENTIFIER_NODE,
				Value: "x",
			},
			input2: &parser.ASTNode{
				Name:  nil,
				Type:  parser.IDENTIFIER_NODE,
				Value: "x",
			},
			expected: true,
		},
		{
			input1: &parser.ASTNode{
				Name:  nil,
				Type:  parser.IDENTIFIER_NODE,
				Value: "x",
			},
			input2: &parser.ASTNode{
				Name:  nil,
				Type:  parser.IDENTIFIER_NODE,
				Value: "y",
			},
			expected: false,
		},
		{
			input1: &parser.ASTNode{
				Name:  nil,
				Type:  parser.IDENTIFIER_NODE,
				Value: "x",
			},
			input2: &parser.ASTNode{
				Name:  nil,
				Type:  parser.STRING_NODE,
				Value: "x",
			},
			expected: false,
		},
		{
			input1: &parser.ASTNode{
				Name:  nil,
				Type:  parser.IDENTIFIER_NODE,
				Value: "x",
			},
			input2: &parser.ASTNode{
				Name:  nil,
				Type:  parser.NUMBER_NODE,
				Value: 3,
			},
			expected: false,
		},
		{
			input1: &parser.ASTNode{
				Name: nil,
				Type: parser.ASSIGNMENT_NODE,
				Value: []*parser.ASTNode{
					{
						Name:  "left",
						Type:  parser.IDENTIFIER_NODE,
						Value: "x",
					},
					{
						Name:  "right",
						Type:  parser.NUMBER_NODE,
						Value: 123,
					},
				},
			},
			input2: &parser.ASTNode{
				Name: nil,
				Type: parser.ASSIGNMENT_NODE,
				Value: []*parser.ASTNode{
					{
						Name:  "left",
						Type:  parser.IDENTIFIER_NODE,
						Value: "x",
					},
					{
						Name:  "right",
						Type:  parser.NUMBER_NODE,
						Value: 123,
					},
				},
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		actual := tt.input1.Equals(tt.input2)
		if actual != tt.expected {
			t.Errorf("expected=%v, got=%v", tt.expected, actual)
		}
	}
}

func TestParser_Parse(t *testing.T) {
	tests := []struct {
		input    []parser.Token
		expected *parser.ASTNode
	}{
		{
			input: []parser.Token{
				&parser.BasicToken{TypeVal: parser.IDENT, LiteralVal: "x"},
				&parser.BasicToken{TypeVal: parser.ASSIGN, LiteralVal: "="},
				&parser.NumberLiteral{Value: 123},
				&parser.EOFToken{},
			},
			expected: &parser.ASTNode{
				Name: nil,
				Type: parser.ROOT_NODE,
				Value: []*parser.ASTNode{
					{
						Name: nil,
						Type: parser.ASSIGNMENT_NODE,
						Value: []*parser.ASTNode{
							{
								Name:  "left",
								Type:  parser.IDENTIFIER_NODE,
								Value: "x",
							},
							{
								Name:  "right",
								Type:  parser.NUMBER_NODE,
								Value: 123,
							},
						},
					},
					{
						Name:  nil,
						Type:  parser.EOF_NODE,
						Value: nil,
					},
				},
			},
		},
		{
			input: []parser.Token{
				&parser.BasicToken{TypeVal: parser.IDENT, LiteralVal: "x"},
				&parser.BasicToken{TypeVal: parser.ASSIGN, LiteralVal: "="},
				&parser.BasicToken{TypeVal: parser.STRING, LiteralVal: "string"},
				&parser.EOFToken{},
			},
			expected: &parser.ASTNode{
				Name: nil,
				Type: parser.ROOT_NODE,
				Value: []*parser.ASTNode{
					{
						Name: nil,
						Type: parser.ASSIGNMENT_NODE,
						Value: []*parser.ASTNode{
							{
								Name:  "left",
								Type:  parser.IDENTIFIER_NODE,
								Value: "x",
							},
							{
								Name:  "right",
								Type:  parser.STRING_NODE,
								Value: "string",
							},
						},
					},
					{
						Name:  nil,
						Type:  parser.EOF_NODE,
						Value: nil,
					},
				},
			},
		},
	}

	for _, tt := range tests {
		p := parser.NewParser(tt.input)
		actual := p.Parse()
		if !actual.Equals(tt.expected) {
			t.Errorf("expected=%v, got=%v", tt.expected, actual)
			t.Logf("Actual: ")
			if actual.Value != nil {
				for _, n := range actual.Value.([]*parser.ASTNode) {
					if n.Value != nil {
						for _, n := range n.Value.([]*parser.ASTNode) {
							t.Logf("%s\n", n)
						}
					}
				}
			}
			t.Logf("Expected: ")
			if tt.expected.Value != nil {
				for _, n := range tt.expected.Value.([]*parser.ASTNode) {
					if n.Value != nil {
						for _, n := range n.Value.([]*parser.ASTNode) {
							t.Logf("%s\n", n)
						}
					}
				}
			}
		}
	}

}
