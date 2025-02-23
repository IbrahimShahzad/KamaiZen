package kamailio_cfg

import (
	sitter "github.com/smacker/go-tree-sitter"
)

type NamedRoute struct {
	Name       string
	Content    string
	StartPoint sitter.Point
	EndPoint   sitter.Point
}

func (nr NamedRoute) String() string {
	return nr.Name
}

func (nr NamedRoute) addContent(content string) {
	nr.Content = content
}

func QueryRoute(a *Analyzer, source_code []byte) *NamedRoute {
	q, err := NewQueryExecutor(
		_ROUTE_DECLARATION_QUERY,
		a.ast.Node,
		a.builder.parser.language,
	)
	if err != nil {
		return nil
	}

	_routeTag := "definition.function"
	_nameTag := "name"

	var nodeContent string
	var route *NamedRoute
	for {
		match, ok := q.NextMatch()
		if !ok {
			break
		}
		for _, capture := range match.Captures {
			node := capture.Node
			captureName := q.query.CaptureNameForId(capture.Index)
			if captureName == _routeTag {
				nodeText := string(capture.Node.Content(source_code))
				nodeContent = nodeText
			}
			if captureName == _nameTag {
				nodeText := string(capture.Node.Content(source_code))
				route = &NamedRoute{
					Name:       nodeText,
					StartPoint: node.StartPoint(),
					EndPoint:   node.EndPoint(),
				}
			}
		}
	}
	if route != nil {
		route.addContent(nodeContent)
		return route
	}
	return nil
}
