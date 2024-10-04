package parser

// Combinators
func seq(parsers ...func() *ASTNode) func() *ASTNode {
	return func() *ASTNode {
		// create a single parent node
		// add all the children to the parent node
		var nodes ASTNode
		for _, parser := range parsers {
			node := parser()
			if node == nil || node.Type == ERROR_NODE {
				return nil
			} else if node == EmptyNode {
				continue
			}
			nodes.addChild(node)
		}
		return &nodes
	}
}

func choice(parsers ...func() *ASTNode) func() *ASTNode {
	return func() *ASTNode {
		for _, parser := range parsers {
			node := parser()
			if node != nil && node != EmptyNode && node.Type != ERROR_NODE {
				return node
			}
			// else continue to the next parser
		}
		return nil
	}
}

func optional(parser func() *ASTNode) func() *ASTNode {
	return func() *ASTNode {
		node := parser()
		if node == nil || node.Type == ERROR_NODE {
			return EmptyNode
		}
		return node
	}
}

func repeat(parser func() *ASTNode) func() *ASTNode {
	return func() *ASTNode {
		var nodes ASTNode
		for {
			node := parser()
			if node == nil || node == EmptyNode || node.Type == ERROR_NODE {
				break
			}
			nodes.addChild(node)
		}
		return &nodes
	}
}
