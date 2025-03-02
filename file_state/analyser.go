package file_state

import (
	"KamaiZen/document_manager"
	"KamaiZen/kamailio_cfg"
	"KamaiZen/lsp"
	"regexp"

	"github.com/rs/zerolog/log"
	sitter "github.com/smacker/go-tree-sitter"
)

// TODO: major refactoring required

// type Tree struct {
// 	nodes map[lsp.DocumentURI]*sitter.Node
// }
//
// // module level state tree
// var treeCache Tree
//
// // NewTree creates and returns a new instance of Tree.
// // It initializes the nodes map.
// //
// // Returns:
// //
// //	Tree - The initialized state tree.
// func NewTree() Tree {
// 	return Tree{
// 		nodes: make(map[lsp.DocumentURI]*sitter.Node),
// 	}
// }
//
// // AddNode adds a node to the state tree for the given document URI.
// //
// // Parameters:
// //
// //	uri lsp.DocumentURI - The URI of the document.
// //	node *sitter.Node - The AST node to be added.
// func (s *Tree) AddNode(uri lsp.DocumentURI, node *sitter.Node) {
// 	s.nodes[uri] = node
// }
//
// // TraverseNode traverses the AST starting from the given node and logs the traversal.
// // It recursively visits all child nodes and logs their information.
// //
// // Parameters:
// //
// //	uri lsp.DocumentURI - The URI of the document.
// //	node *sitter.Node - The starting node for the traversal.
// //	padding int - The padding used for indentation.
// func (s *Tree) TraverseNode(uri lsp.DocumentURI, node *sitter.Node, padding int) {
// 	// traverse the node and print the node
// 	var i uint32
// 	childCount := node.ChildCount()
// 	for i = 0; i < childCount; i++ {
// 		// Print spaces for padding
// 		child := node.Child(int(i))
// 		s.TraverseNode(uri, child, padding+2)
// 	}
// }

// GetNodeDocsAtPosition retrieves the documentation for the node at the given position in the source code.
// Parameters:
// - uri: The URI of the document.
// - position: The position within the document.
// - source_code: The source code as a byte slice.
// Returns:
// - The documentation string for the node at the specified position.
func (s *FileState) GetNodeDocsAtPosition(position lsp.Position, source_code []byte) string {
	node := s.analyser.GetAST().Node
	nodeAtPosition := getNodeAtPosition(node, position)
	if nodeAtPosition == nil {
		log.Error().Msg("Node at position is nil")
		return ""
	}
	switch nodeAtPosition.Type() {
	case kamailio_cfg.IdentifierNodeType:
		switch nodeAtPosition.Parent().Parent().Type() {
		case kamailio_cfg.CallExpressionNodeType:
			functionName := getFunctionName(nodeAtPosition, source_code)
			return document_manager.FindFunctionInAllModules(functionName)
		case kamailio_cfg.AVPNodeType:
			variableName := nodeAtPosition.Content(source_code)
			return s.vars.GetAVPVariable(variableName).GetDocs()
		case kamailio_cfg.VARNodeType:
			variableName := nodeAtPosition.Content(source_code)
			return s.vars.GetLocalVariable(variableName).GetDocs()
		case kamailio_cfg.DlgVarNodeType:
			variableName := nodeAtPosition.Content(source_code)
			return s.vars.GetDlgVariable(variableName).GetDocs()
		}
	case kamailio_cfg.AVPNodeType:
		variableName := nodeAtPosition.ChildByFieldName("name").NamedChild(0).Content(source_code)
		return s.vars.GetAVPVariable(variableName).GetDocs()
	case kamailio_cfg.VARNodeType:
		variableName := nodeAtPosition.ChildByFieldName("name").NamedChild(0).Content(source_code)
		return s.vars.GetLocalVariable(variableName).GetDocs()
	case kamailio_cfg.DlgVarNodeType:
		variableName := nodeAtPosition.ChildByFieldName("name").NamedChild(0).Content(source_code)
		return s.vars.GetDlgVariable(variableName).GetDocs()
	}
	word := nodeAtPosition.Content(source_code)
	// drop special characters
	var nonAlphanumericRegex = regexp.MustCompile(`[^a-zA-Z0-9 _]+`)
	key := nonAlphanumericRegex.ReplaceAllString(word, "")
	docs := document_manager.GetCookBookDocs(key)
	if docs != "" {
		return docs
	}
	log.Error().Str("word", word).Str("key", key).Msg("Documentation not found")
	return "Documentation not found"
}

// getNodeAtPosition finds the node at the specified position within the given AST node.
// Parameters:
// - node: The root AST node.
// - position: The position within the document.
// Returns:
// - The node at the specified position or nil if no such node exists.
func getNodeAtPosition(node *sitter.Node, position lsp.Position) *sitter.Node {
	if node == nil {
		return nil
	}
	nodeAtPosition := node.NamedDescendantForPointRange(
		sitter.Point{
			Row:    uint32(position.Line),
			Column: uint32(position.Character),
		},
		sitter.Point{
			Row:    uint32(position.Line),
			Column: uint32(position.Character),
		})
	return nodeAtPosition
}

// getFunctionName returns the name of the function at the given position in the AST node.
//
// Parameters:
//
//	node *sitter.Node - The AST node to be searched.
//	position lsp.Position - The position within the node.
//	source_code []byte - The source code of the document.
//
// Returns:
//
//	string - The name of the function at the given position.
func getFunctionName(node *sitter.Node, source_code []byte) string {
	if node == nil {
		return ""
	}
	if node.Type() == kamailio_cfg.IdentifierNodeType &&
		node.Parent().Parent().Type() == kamailio_cfg.CallExpressionNodeType &&
		node.Parent().Parent().FieldNameForChild(0) == "function" {
		return node.Content(source_code)
	}
	return ""
}

// TraverseNodeAndApply traverses the AST starting from the given node and applies the given function to each node.
// It recursively visits all child nodes and applies the function to each node.
//
// Parameters:
//
//	node *sitter.Node - The starting node for the traversal.
//	f func(*sitter.Node) - The function to be applied to each node.
func TraverseNodeAndApply(node *sitter.Node, f func(*sitter.Node)) {
	f(node)
	if node.ChildCount() == 0 {
		// get siblings
		if node.NextSibling() != nil {
			TraverseNodeAndApply(node.NextSibling(), f)
		}
		return
	}
	for i := 0; i < int(node.NamedChildCount()); i++ {
		TraverseNodeAndApply(node.NamedChild(i), f)
	}

}

// GetCompletionItems returns a list of completion items for the given document URI.
//
// Parameters:
//
//	uri lsp.DocumentURI - The URI of the document.
//
// Returns:
//
//	[]lsp.CompletionItem - A list of completion items.
func (s FileState) GetCompletionItems() []lsp.CompletionItem {
	var completionItems []lsp.CompletionItem
	functions := document_manager.GetAllAvailableFunctionDocs()
	for _, function := range functions {
		completionItems = append(completionItems, lsp.CompletionItem{
			Detail:        function.Name + "(" + function.Parameters + ")",
			Label:         function.Name + "(" + ")",
			Documentation: function.Description + "\n" + function.Example,
			Kind:          lsp.FUNCTION_COMPLETION,
		})
	}

	for header, description := range kamailio_cfg.SIPHeaders {
		completionItems = append(completionItems, lsp.CompletionItem{
			Detail:        "SIP Header",
			Label:         header,
			Documentation: description,
			Kind:          lsp.VARIABLE_COMPLETION,
		})
	}

	variables := s.vars.GetAVPVariables()
	for variable, value := range variables {
		completionItems = append(completionItems, lsp.CompletionItem{
			Detail:        "AVP",
			Label:         variable,
			Documentation: value.GetDocs(),
			Kind:          lsp.VARIABLE_COMPLETION,
		})
	}

	localVariables := s.vars.GetLocalVariables()
	for variable, value := range localVariables {
		completionItems = append(completionItems, lsp.CompletionItem{
			Detail:        "Local Variable",
			Label:         variable,
			Documentation: value.GetDocs(),
			Kind:          lsp.VARIABLE_COMPLETION,
		})
	}

	dlgVariables := s.vars.GetDlgVariables()
	for variable, value := range dlgVariables {
		completionItems = append(completionItems, lsp.CompletionItem{
			Detail:        "Dialog Variable",
			Label:         variable,
			Documentation: value.GetDocs(),
			Kind:          lsp.VARIABLE_COMPLETION,
		})
	}

	for module := range document_manager.GetAllAvailableModules() {
		completionItems = append(completionItems, lsp.CompletionItem{
			Detail:        "Module",
			Label:         module,
			Documentation: "Module " + module,
			Kind:          lsp.MODULE_COMPLETION,
		})
	}

	for c := range document_manager.GetAllCookBookKeys() {
		completionItems = append(completionItems, lsp.CompletionItem{
			Detail:        "Cookbook",
			Label:         c,
			Documentation: document_manager.GetCookBookDocs(c),
			Kind:          lsp.VARIABLE_COMPLETION,
		})
	}

	return completionItems

}

func (s *FileState) GetRouteDefinitionAtPosition(
	position lsp.Position,
) *kamailio_cfg.Route {
	node := s.analyser.GetAST().Node
	nodeAtPosition := getNodeAtPosition(node, position)
	if nodeAtPosition == nil {
		log.Error().Msg("Node at position is nil")
		return nil
	}
	// extract name of the route from node
	routeName := ""
	routeType := kamailio_cfg.RouteTypeUnknown
	if nodeAtPosition.Type() == kamailio_cfg.RouteCallNodeType {
		routeName = nodeAtPosition.NamedChild(0).Content(s.f.Content())
		routeType = kamailio_cfg.RouteTypeNamed
	} else if nodeAtPosition.Parent().Type() == kamailio_cfg.RouteCallNodeType {
		routeName = nodeAtPosition.Parent().NamedChild(0).Content(s.f.Content())
		routeType = kamailio_cfg.RouteTypeNamed
	}
	return s.routes.FindRouteByNameAndType(routeName, routeType)
}
