package kamailio_cfg

import "KamaiZen/document_manager"

type Function struct {
	document_manager.FunctionDocumentation
	position Position
}

func NewFunction(name string, position Position) Function {
	fdocs := document_manager.GetAllAvailableFunctionDocs()
	docs := fdocs.GetFunction(name)
	if docs == nil {
		return Function{}
	}
	return Function{*docs, position}
}

type Functions struct {
	funcs map[string]Function
	refs  map[string][]Position
}

func NewFunctions() Functions {
	return Functions{
		funcs: make(map[string]Function),
		refs:  make(map[string][]Position),
	}
}

const _FUNCTION_CALL_QUERY = `(call_expression
    function: (_) @name
    arguments: (argument_list) @arguments
    ) @reference.call`

func (f *Functions) AddFunctions(a *Analyzer, sc []byte) {
	q, err := NewQueryExecutor(
		_FUNCTION_CALL_QUERY,
		a.ast.Node,
		a.builder.parser.language,
	)
	if err != nil {
		return
	}
	_callTag := "reference.call"
	_nameTag := "name"
	_argumentsTag := "arguments"

	for {
		match, ok := q.NextMatch()
		if !ok {
			break
		}
		for _, capture := range match.Captures {
			node := capture.Node
			captureName := q.query.CaptureNameForId(capture.Index)
			if captureName == _callTag {
				// add the function to the list of references
				// get the function name
				functionName := capture.Node.Content(sc)
				// get the position of the function call
				pos := Position{node.StartPoint(), node.EndPoint()}
				// add the function to the list of references
				f.refs[functionName] = append(f.refs[functionName], pos)
			}
			if captureName == _nameTag {
				functionName := capture.Node.Content(sc)

				// add only if the function is not already in the list
				if _, ok := f.funcs[functionName]; ok {
					continue
				}
				f.funcs[functionName] = NewFunction(
					functionName,
					Position{node.StartPoint(), node.EndPoint()},
				)
			}
			if captureName == _argumentsTag {
				// TODO: extend functions to keep track of arguments?
				// get the arguments
			}
		}
	}

}
