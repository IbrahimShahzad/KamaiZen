package kamailio_cfg

import (
	"KamaiZen/lsp"
	"github.com/rs/zerolog/log"
	"strings"
)

type ASTVisitor interface {
	Visit(node *ASTNode, a *Analyzer) error
}

type FormattingVisitor struct {
	edits []lsp.TextEdit
}

func NewFormattingVisitor() *FormattingVisitor {
	return &FormattingVisitor{}
}

func (v *FormattingVisitor) GetEdits() []lsp.TextEdit {
	return v.edits
}

// TODO: This is a big todo
func (v *FormattingVisitor) Visit(node *ASTNode, a *Analyzer) error {

	edits := []lsp.TextEdit{}
	block_level := 0

	n := node.Node
	switch node.Node.Type() {
	case AssignmentExpressionNodeType:
		leftNode := n.ChildByFieldName("left")
		rightNode := n.ChildByFieldName("right")
		end_of_statment := n.NextSibling()
		if end_of_statment != nil {
			log.Info().Str("node-type", end_of_statment.Type()).Msg("End of statement")
		}
		// Add the left-hand side, the equals sign with correct spacing, and the right-hand side
		leftNodeContent := string(leftNode.Content(nil))
		rightNodeContent := string(rightNode.Content(nil))
		leftNodeContent = strings.TrimSpace(leftNodeContent)
		rightNodeContent = strings.TrimSpace(rightNodeContent)
		formattedContent := strings.Builder{}
		formattedContent.WriteString(leftNodeContent)
		formattedContent.WriteString(" = ") // Ensure exactly one space on both sides
		formattedContent.WriteString(rightNodeContent)
		log.Info().Str("Formatted content", formattedContent.String())
		edit := lsp.NewTextEdit(leftNode, rightNode, formattedContent.String())
		edits = append(edits, edit)
		// update the tree as well to reflect the changes
		// kamailio_cfg.UpdateTree(parser.GetTree(), leftNode, rightNode, formattedContent.String())
		// parser.UpdateTree([]byte(formattedContent.String()))
	case ParenthesizedExpressionNodeType:
		if n.ChildCount() != 3 {
			log.Info().Msg("Parenthesized expression does not have 3 children")
			break
		}

		leftBrace := n.Child(0)
		expression := n.Child(1)
		rightBrace := n.Child(2)
		formattedContent := strings.Builder{}
		formattedContent.WriteString("(")
		formattedContent.WriteString(string(expression.Content(nil)))
		formattedContent.WriteString(") ")
		edit := lsp.NewTextEdit(leftBrace, rightBrace, formattedContent.String())
		edits = append(edits, edit)
		// kamailio_cfg.UpdateTree(parser.GetTree(), leftBrace, rightBrace, formattedContent.String())
		// parser.UpdateTree([]byte(formattedContent.String()))
	case TopLevelAssignmentNodeType:
		key := n.ChildByFieldName("key")
		value := n.ChildByFieldName("value")
		formattedContent := strings.Builder{}
		formattedContent.WriteString(string(key.Content(nil)))
		formattedContent.WriteString("=")
		formattedContent.WriteString(string(value.Content(nil)))
		edit := lsp.NewTextEdit(key, value, formattedContent.String())
		edits = append(edits, edit)
		// kamailio_cfg.UpdateTree(parser.GetTree(), key, value, formattedContent.String())
		// parser.UpdateTree([]byte(formattedContent.String()))
	case CompoundStatementNodeType:
		leftBrace := n.Child(0)
		rightBrace := n.Child(int(n.ChildCount() - 1))
		edit := lsp.NewTextEdit(leftBrace, leftBrace, "{")
		edits = append(edits, edit)
		// kamailio_cfg.UpdateTree(parser.GetTree(), leftBrace, leftBrace, "{")
		// parser.UpdateTree([]byte("{"))
		edit = lsp.NewTextEdit(rightBrace, rightBrace, "}")
		edits = append(edits, edit)
		// kamailio_cfg.UpdateTree(parser.GetTree(), rightBrace, rightBrace, "}")
		// parser.UpdateTree([]byte("}"))
	case "if_statement":
		ifKeyword := n.Child(0) // if keyword
		formattedContent := strings.Builder{}
		content := "if "
		formattedContent.WriteString(content)
		edit := lsp.NewTextEdit(ifKeyword, ifKeyword, formattedContent.String())
		edits = append(edits, edit)
	case "block_start":
		block_level++
		log.Info().Msgf("INCREASING Block level: %d", block_level)
	case "block_end":
		block_level--
		log.Info().Msgf("DECREASING Block level: %d", block_level)
	case "core_function":
		content := string(n.Content(nil))
		edit := lsp.NewTextEdit(n, n, content)
		edits = append(edits, edit)
	case "call_expression":
		// content := string(sourceCode[node.StartByte():node.EndByte()])
		// content = stringPadding + content
		// edit := lsp.NewTextEdit(node, node, content)
		// edits = append(edits, edit)

		// default:
		// 	log.Printf("Node type: %v -- continue\n", node.Type())

	}
	for i := 0; i < int(node.Node.ChildCount()); i++ {
		child := node.Node.Child(i)
		v.Visit(&ASTNode{Node: child}, a)
	}
	v.edits = append(v.edits, edits...)
	return nil
}

// func applyEditToContent(content string, edit lsp.TextEdit) string {
// 	log.Info("Applying edits to content")
// 	sourceCode := []byte(content)
// 	// split source code 2d array of lines and characters
// 	lines := bytes.Split(sourceCode, []byte("\n"))
// 	var splitSourceCode [][]byte
// 	for _, line := range lines {
// 		var chars []byte
// 		for _, char := range line {
// 			chars = append(chars, char)
// 		}
// 		splitSourceCode = append(splitSourceCode, chars)
// 		// splitSourceCode = append(splitSourceCode, bytes.Split(line, nil))
// 	}
//
// 	log.Info("Edit: ", edit)
// 	startLine := edit.Range.Start.Line
// 	startChar := edit.Range.Start.Character
// 	endLine := edit.Range.End.Line
// 	endChar := edit.Range.End.Character
// 	newText := edit.NewText
// 	// update source code with new text
// 	// Replace the specified range with the new text
// 	newTextBytes := []byte(newText)
// 	splitSourceCode[startLine] = append(splitSourceCode[startLine][:startChar], append(newTextBytes, splitSourceCode[endLine][endChar:]...)...)
//
// 	// Reconstruct the source code from the modified 2D array
// 	var result []byte
// 	for _, line := range splitSourceCode {
// 		result = append(result, line...)
// 		result = append(result, '\n')
// 	}
// 	result = bytes.TrimSuffix(result, []byte("\n")) // Remove the last newline
//
// 	return string(result)
// }
//
// // All siblings should be on the same level
// func IndentKamailioCfg(content string, tabsize int) ([]lsp.TextEdit, error) {
// 	sourceCode := []byte(content)
// 	parser := kamailio_cfg.NewParser()
// 	parser.SetLanguage()
//
// 	log.Info("Indenting Kamailio configuration file")
// 	// this is the root node of the tree
// 	root := parser.Parse(sourceCode)
// 	edits := []lsp.TextEdit{}
// 	block_level := 0
//
// 	// go over the top level statements // each sibling is on the same level
// 	// top down call recursively IndentSiblings
// 	edits, err := IndentSiblings(root, sourceCode, tabsize, block_level)
// 	if err != nil {
// 		log.Info("Error indenting siblings: ", err)
// 		return nil, err
// 	}
// 	return edits, nil
// }
//
// func IndentSiblings(root *sitter.Node, sourceCode []byte, tabsize int, block_level int) ([]lsp.TextEdit, error) {
// 	edits := []lsp.TextEdit{}
// 	for i := 0; i < int(root.NamedChildCount()); i++ {
// 		node := root.NamedChild(i)
// 		stringPadding := strings.Repeat(" ", block_level*tabsize)
// 		switch node.Type() {
// 		case "block_start":
// 			block_level++
// 		case "block_end":
// 			block_level--
// 		default:
// 			// Add the left-hand side, the equals sign with correct spacing, and the right-hand side
// 			nodeContent := string(sourceCode[node.StartByte():node.EndByte()])
// 			nodeContent = strings.TrimSpace(nodeContent)
// 			nodeContent = stringPadding + nodeContent
// 			formattedContent := strings.Builder{}
// 			formattedContent.WriteString(nodeContent)
// 			edit := lsp.NewTextEdit(node, node, formattedContent.String())
// 			edits = append(edits, edit)
// 		}
// 	}
// 	return edits, nil
// }
//
// func formatAssignmentExpression(tree *sitter.Tree, node *sitter.Node, parser *sitter.Parser, sourceCode []byte) ([]lsp.TextEdit, *sitter.Tree, error) {
// 	edits := []lsp.TextEdit{}
// 	leftNode := node.ChildByFieldName("left")
// 	rightNode := node.ChildByFieldName("right")
// 	// Add the left-hand side, the equals sign with correct spacing, and the right-hand side
// 	leftNodeContent := string(sourceCode[leftNode.StartByte():leftNode.EndByte()])
// 	rightNodeContent := string(sourceCode[rightNode.StartByte():rightNode.EndByte()])
// 	leftNodeContent = strings.TrimSpace(leftNodeContent)
// 	rightNodeContent = strings.TrimSpace(rightNodeContent)
// 	formattedContent := strings.Builder{}
// 	formattedContent.WriteString(leftNodeContent)
// 	formattedContent.WriteString(" = ") // Ensure exactly one space on both sides
// 	formattedContent.WriteString(rightNodeContent)
// 	log.Info("Formatted content: ", formattedContent.String())
// 	edit := lsp.NewTextEdit(leftNode, rightNode, formattedContent.String())
// 	edits = append(edits, edit)
// 	// update the tree as well to reflect the changes
// 	kamailio_cfg.UpdateTree(tree, leftNode, rightNode, formattedContent.String())
// 	tree, err := parser.ParseCtx(context.Background(), tree, []byte(formattedContent.String()))
// 	if err != nil {
// 		log.Info("Error parsing the formatted content: ", err)
// 		return nil, nil, err
// 	}
// 	return edits, tree, nil
// }
