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
		Name     string
		input    []parser.Token
		expected *parser.ASTNode
	}{
		{
			Name: "Test 1",
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
						Type: parser.TOP_LEVEL_STATEMENT_NODE,
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
										Name:  nil,
										Type:  parser.OPERATOR_NODE,
										Value: "=",
									},
									{
										Name:  "right",
										Type:  parser.NUMBER_NODE,
										Value: 123,
									},
								},
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
			Name: "Test 2",
			input: []parser.Token{
				&parser.BasicToken{TypeVal: parser.IDENT, LiteralVal: "x"},
				&parser.BasicToken{TypeVal: parser.ASSIGN, LiteralVal: "="},
				&parser.BasicToken{TypeVal: parser.STRING, LiteralVal: "Hello World!"},
				&parser.EOFToken{},
			},
			expected: &parser.ASTNode{
				Name: nil,
				Type: parser.ROOT_NODE,
				Value: []*parser.ASTNode{
					{
						Name: nil,
						Type: parser.TOP_LEVEL_STATEMENT_NODE,
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
										Name:  nil,
										Type:  parser.OPERATOR_NODE,
										Value: "=",
									},
									{
										Name:  "right",
										Type:  parser.STRING_NODE,
										Value: "Hello World!",
									},
								},
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
			Name: "Test 3",
			input: []parser.Token{
				&parser.BasicToken{TypeVal: parser.PREPROC, LiteralVal: "KAMAILIO"},
				&parser.EOFToken{},
			},
			expected: &parser.ASTNode{
				Name: nil,
				Type: parser.ROOT_NODE,
				Value: []*parser.ASTNode{
					{
						Name: nil,
						Type: parser.TOP_LEVEL_STATEMENT_NODE,
						Value: []*parser.ASTNode{
							{
								Name:  nil,
								Type:  parser.FILE_STARTER_NODE,
								Value: "KAMAILIO",
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
			Name: "Test 4",
			input: []parser.Token{
				&parser.BasicToken{TypeVal: parser.ROUTE, LiteralVal: "request_route"},
				&parser.BasicToken{TypeVal: parser.LBRACE, LiteralVal: "{"},
				&parser.BasicToken{TypeVal: parser.IDENT, LiteralVal: "x"},
				&parser.BasicToken{TypeVal: parser.ASSIGN, LiteralVal: "="},
				&parser.NumberLiteral{Value: 123},
				&parser.BasicToken{TypeVal: parser.SEMICOLON, LiteralVal: ";"},
				&parser.BasicToken{TypeVal: parser.RBRACE, LiteralVal: "}"},
				&parser.EOFToken{},
			},
			expected: &parser.ASTNode{
				Name: nil,
				Type: parser.ROOT_NODE,
				Value: []*parser.ASTNode{
					{
						Name: nil,
						Type: parser.TOP_LEVEL_STATEMENT_NODE,
						Value: []*parser.ASTNode{
							{
								Name: nil,
								Type: parser.REQUEST_ROUTE_NODE,
								Value: []*parser.ASTNode{
									{
										Name: nil,
										Type: parser.BLOCK_NODE,
										Value: []*parser.ASTNode{
											{
												Name: nil,
												Type: parser.STATEMENT_NODE,
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
																Name:  nil,
																Type:  parser.OPERATOR_NODE,
																Value: "=",
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
														Type:  parser.EOS_NODE,
														Value: nil,
													},
												},
											},
										},
									},
								},
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
			Name: "Test 5",
			input: []parser.Token{
				&parser.BasicToken{TypeVal: parser.PREPROC, LiteralVal: "include"},
				&parser.BasicToken{TypeVal: parser.STRING, LiteralVal: "file.so"},
				&parser.EOFToken{},
			},
			expected: &parser.ASTNode{
				Name: nil,
				Type: parser.ROOT_NODE,
				Value: []*parser.ASTNode{
					{
						Name: nil,
						Type: parser.TOP_LEVEL_STATEMENT_NODE,
						Value: []*parser.ASTNode{
							{
								Name: "include",
								Type: parser.INCLUDE_FILE_NODE,
								Value: []*parser.ASTNode{
									{
										Name:  nil,
										Type:  parser.STRING_NODE,
										Value: "file.so",
									},
								},
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
			Name: "Test 6",
			input: []parser.Token{
				&parser.BasicToken{TypeVal: parser.ROUTE, LiteralVal: "route"},
				&parser.BasicToken{TypeVal: parser.LBRACKET, LiteralVal: "["},
				&parser.BasicToken{TypeVal: parser.IDENT, LiteralVal: "NEW_ROUTE"},
				&parser.BasicToken{TypeVal: parser.RBRACKET, LiteralVal: "]"},
				&parser.BasicToken{TypeVal: parser.LBRACE, LiteralVal: "{"},
				&parser.BasicToken{TypeVal: parser.IDENT, LiteralVal: "x"},
				&parser.BasicToken{TypeVal: parser.ASSIGN, LiteralVal: "="},
				&parser.NumberLiteral{Value: 123},
				&parser.BasicToken{TypeVal: parser.SEMICOLON, LiteralVal: ";"},
				&parser.BasicToken{TypeVal: parser.RBRACE, LiteralVal: "}"},
				&parser.EOFToken{},
			},
			expected: &parser.ASTNode{
				Name: nil,
				Type: parser.ROOT_NODE,
				Value: []*parser.ASTNode{
					{
						Name: nil,
						Type: parser.TOP_LEVEL_STATEMENT_NODE,
						Value: []*parser.ASTNode{
							{
								Name: "NEW_ROUTE",
								Type: parser.ROUTE_NODE,
								Value: []*parser.ASTNode{
									{
										Name: nil,
										Type: parser.BLOCK_NODE,
										Value: []*parser.ASTNode{
											{
												Name: nil,
												Type: parser.STATEMENT_NODE,
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
																Name:  nil,
																Type:  parser.OPERATOR_NODE,
																Value: "=",
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
														Type:  parser.EOS_NODE,
														Value: nil,
													},
												},
											},
										},
									},
								},
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
			Name: "Test 7",
			input: []parser.Token{
				&parser.BasicToken{TypeVal: parser.ROUTE, LiteralVal: "route"},
				&parser.BasicToken{TypeVal: parser.LBRACKET, LiteralVal: "["},
				&parser.NumberLiteral{Value: 1},
				&parser.BasicToken{TypeVal: parser.RBRACKET, LiteralVal: "]"},
				&parser.BasicToken{TypeVal: parser.LBRACE, LiteralVal: "{"},
				&parser.BasicToken{TypeVal: parser.IDENT, LiteralVal: "x"},
				&parser.BasicToken{TypeVal: parser.ASSIGN, LiteralVal: "="},
				&parser.NumberLiteral{Value: 123},
				&parser.BasicToken{TypeVal: parser.SEMICOLON, LiteralVal: ";"},
				&parser.BasicToken{TypeVal: parser.IDENT, LiteralVal: "y"},
				&parser.BasicToken{TypeVal: parser.ASSIGN, LiteralVal: "="},
				&parser.NumberLiteral{Value: 456},
				&parser.BasicToken{TypeVal: parser.SEMICOLON, LiteralVal: ";"},
				&parser.CoreVariableToken{TypeVal: parser.CORE_VARIABLE, VariableType: "var", VariableName: "x"},
				&parser.BasicToken{TypeVal: parser.ASSIGN, LiteralVal: "="},
				&parser.CoreVariableToken{TypeVal: parser.CORE_VARIABLE, VariableType: "avp", VariableName: "y"},
				&parser.BasicToken{TypeVal: parser.SEMICOLON, LiteralVal: ";"},
				&parser.BasicToken{TypeVal: parser.RBRACE, LiteralVal: "}"},
				&parser.EOFToken{},
			},
			expected: &parser.ASTNode{
				Name: nil,
				Type: parser.ROOT_NODE,
				Value: []*parser.ASTNode{
					{
						Name: nil,
						Type: parser.TOP_LEVEL_STATEMENT_NODE,
						Value: []*parser.ASTNode{
							{
								Name: "1",
								Type: parser.ROUTE_NODE,
								Value: []*parser.ASTNode{
									{
										Name: nil,
										Type: parser.BLOCK_NODE,
										Value: []*parser.ASTNode{
											{
												Name: nil,
												Type: parser.STATEMENT_NODE,
												Value: []*parser.ASTNode{
													{
														Name: nil,
														Type: parser.ASSIGNMENT_NODE,
														Value: []*parser.ASTNode{
															{Name: "left", Type: parser.IDENTIFIER_NODE, Value: "x"},
															{Name: nil, Type: parser.OPERATOR_NODE, Value: "="},
															{Name: "right", Type: parser.NUMBER_NODE, Value: 123},
														},
													},
													{Name: nil, Type: parser.EOS_NODE, Value: nil},
												},
											},
											{
												Name: nil,
												Type: parser.STATEMENT_NODE,
												Value: []*parser.ASTNode{
													{
														Name: nil,
														Type: parser.ASSIGNMENT_NODE,
														Value: []*parser.ASTNode{
															{Name: "left", Type: parser.IDENTIFIER_NODE, Value: "y"},
															{Name: nil, Type: parser.OPERATOR_NODE, Value: "="},
															{Name: "right", Type: parser.NUMBER_NODE, Value: 456},
														},
													},
													{Name: nil, Type: parser.EOS_NODE, Value: nil},
												},
											},
											{
												Name: nil,
												Type: parser.STATEMENT_NODE,
												Value: []*parser.ASTNode{
													{
														Name: nil,
														Type: parser.ASSIGNMENT_NODE,
														Value: []*parser.ASTNode{
															{
																Name: "left",
																Type: parser.CORE_VAR_VAR_NODE,
																Value: []*parser.ASTNode{{
																	Name:  nil,
																	Type:  parser.IDENTIFIER_NODE,
																	Value: "x",
																}},
															},
															{Name: nil, Type: parser.OPERATOR_NODE, Value: "="},
															{
																Name: "right",
																Type: parser.CORE_VAR_AVP_NODE,
																Value: []*parser.ASTNode{{
																	Name:  nil,
																	Type:  parser.IDENTIFIER_NODE,
																	Value: "y",
																}},
															},
														},
													},
													{Name: nil, Type: parser.EOS_NODE, Value: nil},
												},
											},
										},
									},
								},
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
			Name: "Test 8",
			input: []parser.Token{
				&parser.BasicToken{TypeVal: parser.PREPROC, LiteralVal: "KAMAILIO"},
				&parser.BasicToken{TypeVal: parser.ROUTE, LiteralVal: "request_route"},
				&parser.BasicToken{TypeVal: parser.LBRACE, LiteralVal: "{"},
				&parser.BasicToken{TypeVal: parser.KEYWORD, LiteralVal: "if"},
				&parser.BasicToken{TypeVal: parser.LPAREN, LiteralVal: "("},
				&parser.CoreVariableToken{TypeVal: parser.CORE_VARIABLE, VariableType: "var", VariableName: "x"},
				&parser.BasicToken{TypeVal: parser.REGEX_OP, LiteralVal: "=~"},
				&parser.BasicToken{TypeVal: parser.STRING, LiteralVal: "sip:.*@.*"},
				&parser.BasicToken{TypeVal: parser.RPAREN, LiteralVal: ")"},
				&parser.BasicToken{TypeVal: parser.LBRACE, LiteralVal: "{"},
				&parser.BasicToken{TypeVal: parser.KEYWORD, LiteralVal: "exit"},
				&parser.BasicToken{TypeVal: parser.SEMICOLON, LiteralVal: ";"},
				&parser.BasicToken{TypeVal: parser.RBRACE, LiteralVal: "}"},
				&parser.BasicToken{TypeVal: parser.RBRACE, LiteralVal: "}"},
				&parser.EOFToken{},
			},
			expected: &parser.ASTNode{
				Name: nil,
				Type: parser.ROOT_NODE,
				Value: []*parser.ASTNode{
					{
						Name: nil,
						Type: parser.TOP_LEVEL_STATEMENT_NODE,
						Value: []*parser.ASTNode{
							{
								Name:  nil,
								Type:  parser.FILE_STARTER_NODE,
								Value: "KAMAILIO",
							},
						},
					},
					{
						Name: nil,
						Type: parser.TOP_LEVEL_STATEMENT_NODE,
						Value: []*parser.ASTNode{
							{
								Name: nil,
								Type: parser.REQUEST_ROUTE_NODE,
								Value: []*parser.ASTNode{
									{
										Name: nil,
										Type: parser.BLOCK_NODE,
										Value: []*parser.ASTNode{
											{
												Name: "if",
												Type: parser.STATEMENT_NODE,
												Value: []*parser.ASTNode{
													{
														Name: "condition",
														Type: parser.BINARY_EXPR_NODE,
														Value: []*parser.ASTNode{
															{
																Name: "left",
																Type: parser.CORE_VAR_VAR_NODE,
																Value: []*parser.ASTNode{
																	{
																		Name:  nil,
																		Type:  parser.IDENTIFIER_NODE,
																		Value: "x",
																	},
																},
															},
															{
																Name:  "operator",
																Type:  parser.OPERATOR_NODE,
																Value: "=~",
															},
															{
																Name:  "right",
																Type:  parser.STRING_NODE,
																Value: "sip:.*@.*",
															},
														},
													},
													{
														Name: "consequence",
														Type: parser.BLOCK_NODE,
														Value: []*parser.ASTNode{
															{
																Name: nil,
																Type: parser.STATEMENT_NODE,
																Value: []*parser.ASTNode{
																	{
																		Name:  nil,
																		Type:  parser.KEYWORD_NODE,
																		Value: "exit",
																	},
																	{
																		Name:  nil,
																		Type:  parser.EOS_NODE,
																		Value: nil,
																	},
																},
															},
														},
													},
												},
											},
										},
									},
								},
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
		expected := tt.expected
		parser.UpdateASTLevel(expected)
		if !actual.Equals(tt.expected) {
			t.Logf("Test failed for %s\n", tt.Name)
			t.Errorf("\n-------------------\nexpected=\n%s\n---------------------\ngot=\n%s", expected, actual)
		} else {
			t.Logf("Test passed for %s\n", tt.Name)
		}
	}

}
