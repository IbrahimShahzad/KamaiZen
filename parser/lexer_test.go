package parser_test

import (
	"KamaiZen/parser"
	"reflect"
	"testing"
)

func TestLexer_Tokenise(t *testing.T) {
	tests := []struct {
		input    string
		expected []parser.Token
	}{
		{
			input: "123",
			expected: []parser.Token{
				&parser.NumberLiteral{Value: 123},
				&parser.EOFToken{},
			},
		},
		{
			input: `"hello"`,
			expected: []parser.Token{
				&parser.BasicToken{TypeVal: parser.STRING, LiteralVal: "hello"},
				&parser.EOFToken{},
			},
		},
		{
			input: "identifier",
			expected: []parser.Token{
				&parser.BasicToken{TypeVal: parser.IDENT, LiteralVal: "identifier"},
				&parser.EOFToken{},
			},
		},
		{
			input: `123 identifier "hello"`,
			expected: []parser.Token{
				&parser.NumberLiteral{Value: 123},
				&parser.BasicToken{TypeVal: parser.IDENT, LiteralVal: "identifier"},
				&parser.BasicToken{TypeVal: parser.STRING, LiteralVal: "hello"},
				&parser.EOFToken{},
			},
		},
		{
			input: "   123   ",
			expected: []parser.Token{
				&parser.NumberLiteral{Value: 123},
				&parser.EOFToken{},
			},
		},
		{
			input: "// this is a comment\n123",
			expected: []parser.Token{
				&parser.NumberLiteral{Value: 123},
				&parser.EOFToken{},
			},
		},
		{
			input: "/* this is a \n multi-line comment */ 123",
			expected: []parser.Token{
				&parser.NumberLiteral{Value: 123},
				&parser.EOFToken{},
			},
		},
		{
			input: `"multi-line"
					" string"`,
			expected: []parser.Token{
				&parser.BasicToken{TypeVal: parser.STRING, LiteralVal: "multi-line string"},
				&parser.EOFToken{},
			},
		},

		{
			input: `x = 123 "hello world"`,
			expected: []parser.Token{
				&parser.BasicToken{TypeVal: parser.IDENT, LiteralVal: "x"},
				&parser.BasicToken{TypeVal: parser.ASSIGN, LiteralVal: "="},
				&parser.NumberLiteral{Value: 123},
				&parser.BasicToken{TypeVal: parser.STRING, LiteralVal: "hello world"},
				&parser.EOFToken{},
			},
		},
		{
			input: `x = 123 + 456`,
			expected: []parser.Token{
				&parser.BasicToken{TypeVal: parser.IDENT, LiteralVal: "x"},
				&parser.BasicToken{TypeVal: parser.ASSIGN, LiteralVal: "="},
				&parser.NumberLiteral{Value: 123},
				&parser.BasicToken{TypeVal: parser.ADD_OP, LiteralVal: "+"},
				&parser.NumberLiteral{Value: 456},
				&parser.EOFToken{},
			},
		},
		{
			input: `x = 123 + 456 / 789 * 1011`,
			expected: []parser.Token{
				&parser.BasicToken{TypeVal: parser.IDENT, LiteralVal: "x"},
				&parser.BasicToken{TypeVal: parser.ASSIGN, LiteralVal: "="},
				&parser.NumberLiteral{Value: 123},
				&parser.BasicToken{TypeVal: parser.ADD_OP, LiteralVal: "+"},
				&parser.NumberLiteral{Value: 456},
				&parser.BasicToken{TypeVal: parser.DIV_OP, LiteralVal: "/"},
				&parser.NumberLiteral{Value: 789},
				&parser.BasicToken{TypeVal: parser.MUL_OP, LiteralVal: "*"},
				&parser.NumberLiteral{Value: 1011},
				&parser.EOFToken{},
			},
		},
		{
			input: `if ( x <= 0 ) { return y; } else { break; }`,
			expected: []parser.Token{
				&parser.BasicToken{TypeVal: parser.KEYWORD, LiteralVal: "if"},
				&parser.BasicToken{TypeVal: parser.LPAREN, LiteralVal: "("},
				&parser.BasicToken{TypeVal: parser.IDENT, LiteralVal: "x"},
				&parser.BasicToken{TypeVal: parser.LE_OP, LiteralVal: "<="},
				&parser.NumberLiteral{Value: 0},
				&parser.BasicToken{TypeVal: parser.RPAREN, LiteralVal: ")"},
				&parser.BasicToken{TypeVal: parser.LBRACE, LiteralVal: "{"},
				&parser.BasicToken{TypeVal: parser.KEYWORD, LiteralVal: "return"},
				&parser.BasicToken{TypeVal: parser.IDENT, LiteralVal: "y"},
				&parser.BasicToken{TypeVal: parser.SEMICOLON, LiteralVal: ";"},
				&parser.BasicToken{TypeVal: parser.RBRACE, LiteralVal: "}"},
				&parser.BasicToken{TypeVal: parser.KEYWORD, LiteralVal: "else"},
				&parser.BasicToken{TypeVal: parser.LBRACE, LiteralVal: "{"},
				&parser.BasicToken{TypeVal: parser.KEYWORD, LiteralVal: "break"},
				&parser.BasicToken{TypeVal: parser.SEMICOLON, LiteralVal: ";"},
				&parser.BasicToken{TypeVal: parser.RBRACE, LiteralVal: "}"},
				&parser.EOFToken{},
			},
		},
		{
			input: `$var(x) = 123;`,
			expected: []parser.Token{
				&parser.CoreVariableToken{TypeVal: parser.CORE_VARIABLE, VariableType: "var", VariableName: "x"},
				&parser.BasicToken{TypeVal: parser.ASSIGN, LiteralVal: "="},
				&parser.NumberLiteral{Value: 123},
				&parser.BasicToken{TypeVal: parser.SEMICOLON, LiteralVal: ";"},
				&parser.EOFToken{},
			},
		},
		{
			input: `request_route { x = 123; }`,
			expected: []parser.Token{
				&parser.BasicToken{TypeVal: parser.ROUTE, LiteralVal: "request_route"},
				&parser.BasicToken{TypeVal: parser.LBRACE, LiteralVal: "{"},
				&parser.BasicToken{TypeVal: parser.IDENT, LiteralVal: "x"},
				&parser.BasicToken{TypeVal: parser.ASSIGN, LiteralVal: "="},
				&parser.NumberLiteral{Value: 123},
				&parser.BasicToken{TypeVal: parser.SEMICOLON, LiteralVal: ";"},
				&parser.BasicToken{TypeVal: parser.RBRACE, LiteralVal: "}"},
				&parser.EOFToken{},
			},
		},
		{
			input: `if ($var(x) == $avp(y)){ return $var(z); }`,
			expected: []parser.Token{
				&parser.BasicToken{TypeVal: parser.KEYWORD, LiteralVal: "if"},
				&parser.BasicToken{TypeVal: parser.LPAREN, LiteralVal: "("},
				&parser.CoreVariableToken{TypeVal: parser.CORE_VARIABLE, VariableType: "var", VariableName: "x"},
				&parser.BasicToken{TypeVal: parser.EQ_OP, LiteralVal: "=="},
				&parser.CoreVariableToken{TypeVal: parser.CORE_VARIABLE, VariableType: "avp", VariableName: "y"},
				&parser.BasicToken{TypeVal: parser.RPAREN, LiteralVal: ")"},
				&parser.BasicToken{TypeVal: parser.LBRACE, LiteralVal: "{"},
				&parser.BasicToken{TypeVal: parser.KEYWORD, LiteralVal: "return"},
				&parser.CoreVariableToken{TypeVal: parser.CORE_VARIABLE, VariableType: "var", VariableName: "z"},
				&parser.BasicToken{TypeVal: parser.SEMICOLON, LiteralVal: ";"},
				&parser.BasicToken{TypeVal: parser.RBRACE, LiteralVal: "}"},
				&parser.EOFToken{},
			},
		},
		{
			input: `# this is a comment
					#!define VAR 123
					## This is also a comment
					/* This is a multi-line
					 *  comment */
					#!trydef X 456`,
			expected: []parser.Token{
				&parser.BasicToken{TypeVal: parser.PREPROC, LiteralVal: "define"},
				&parser.BasicToken{TypeVal: parser.IDENT, LiteralVal: "VAR"},
				&parser.NumberLiteral{Value: 123},
				&parser.BasicToken{TypeVal: parser.PREPROC, LiteralVal: "trydef"},
				&parser.BasicToken{TypeVal: parser.IDENT, LiteralVal: "X"},
				&parser.NumberLiteral{Value: 456},
				&parser.EOFToken{},
			},
		},
		{
			input: `#!KAMAILIO`,
			expected: []parser.Token{
				&parser.BasicToken{TypeVal: parser.PREPROC, LiteralVal: "KAMAILIO"},
				&parser.EOFToken{},
			},
		},
		// TODO: Fix this test
		// {
		// 	input: `!!KAMAILIO`,
		// 	expected: []parser.Token{
		// 		&parser.BasicToken{TypeVal: parser.PREPROC, LiteralVal: "KAMAILIO"},
		// 		&parser.EOFToken{},
		// 	},
		// },
		{
			input: `#!include "file.so"`,
			expected: []parser.Token{
				&parser.BasicToken{TypeVal: parser.PREPROC, LiteralVal: "include"},
				&parser.BasicToken{TypeVal: parser.STRING, LiteralVal: "file.so"},
				&parser.EOFToken{},
			},
		},
		{
			input: `include_file("file.so")`,
			expected: []parser.Token{
				&parser.BasicToken{TypeVal: parser.KEYWORD, LiteralVal: "include_file"},
				&parser.BasicToken{TypeVal: parser.LPAREN, LiteralVal: "("},
				&parser.BasicToken{TypeVal: parser.STRING, LiteralVal: "file.so"},
				&parser.BasicToken{TypeVal: parser.RPAREN, LiteralVal: ")"},
				&parser.EOFToken{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("test %s panicked: %v", tt.input, r)
				}
			}()

			lexer := parser.NewLexer([]byte(tt.input))
			tokens := lexer.Tokenise()

			if !reflect.DeepEqual(tokens, tt.expected) {
				for i, token := range tokens {
					if i >= len(tt.expected) {
						t.Errorf("unexpected token: %v", token)
						continue
					}
					if !reflect.DeepEqual(token, tt.expected[i]) {
						t.Errorf("expected: %v, got: %v", tt.expected[i], token)
					}
				}
			}
		})
	}
}
