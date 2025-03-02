package kamailio_cfg

import (
	"github.com/rs/zerolog/log"
	sitter "github.com/smacker/go-tree-sitter"
)

const (
	_AVP_IDENTIFIER     = "$avp"
	_VAR_IDENTIFIER     = "$var"
	_DLG_VAR_IDENTIFIER = "$dlg_var"
)

const (
	_AVP_SCOPE = "Transaction"
	_DLG_SCOPE = "Dialog"
	_VAR_SCOPE = "Local"
)

// keeps track of variables declared int he program

// global variables
// gloabl variables are avps
// var avpVariables map[string]Variable
//
// // local variables
// var localVariables map[string]Variable
//
// var dlgVariables map[string]Variable

type Position struct {
	start sitter.Point
	end   sitter.Point
}

// Variable struct
type Variable struct {
	name       string
	value      string
	scope      string
	identifier string
	position   Position
}

func (v Variable) GetDocs() string {

	var header string

	// scope also determines the type of variable
	switch v.scope {
	case _AVP_SCOPE:
		header = "## User defined AVP\n\n\t" + v.name + "\n\n"
	case _DLG_SCOPE:
		header = "## User defined Dialog Variable\n\n\t" + v.name + "\n\n"
	case _VAR_SCOPE:
		header = "## User defined Local Variable\n\n\t" + v.name + "\n\n"
	default:
		// This shouldn't happen
		header = "## User defined Variable\n\n\t" + v.name + "\n\n"
	}
	return header + "### Value\n\n\t" + v.value + "\n\n" + "### Scope\n\n\t" + v.scope + "\n"
}

type Variables struct {
	avpVariables   map[string]Variable
	localVariables map[string]Variable
	dlgVariables   map[string]Variable
}

func (v *Variables) GetAVPVariables() map[string]Variable {
	return v.avpVariables
}

func (v *Variables) GetAVPVariable(name string) Variable {
	id := _AVP_IDENTIFIER + "(" + name + ")"
	return v.avpVariables[id]
}

func (v *Variables) GetLocalVariable(name string) Variable {
	id := _VAR_IDENTIFIER + "(" + name + ")"
	return v.localVariables[id]
}

func (v *Variables) GetDlgVariable(name string) Variable {
	id := _DLG_VAR_IDENTIFIER + "(" + name + ")"
	return v.dlgVariables[id]
}

func (v *Variables) GetLocalVariables() map[string]Variable {
	return v.localVariables
}

func (v *Variables) GetDlgVariables() map[string]Variable {
	return v.dlgVariables
}

func NewVariables() *Variables {
	return &Variables{
		avpVariables:   make(map[string]Variable),
		localVariables: make(map[string]Variable),
		dlgVariables:   make(map[string]Variable),
	}
}

func ExtractVariables(a *Analyzer, source_code []byte) Variables {
	avpVariables := make(map[string]Variable)
	localVariables := make(map[string]Variable)
	dlgVariables := make(map[string]Variable)
	addAVPVariable := func(name string, value string, identifier string, position Position) {
		avpVariables[name] = Variable{name, value, _AVP_SCOPE, identifier, position}
	}

	addLocalVariable := func(name string, value string, scope string, identifier string, position Position) {
		localVariables[name] = Variable{name, value, scope, identifier, position}
	}

	addDlgVariable := func(name string, value string, scope string, identifier string, position Position) {
		dlgVariables[name] = Variable{name, value, scope, identifier, position}
	}

	q, err := NewQueryExecutor(_ASSINGMENT_QUERY, a.ast.Node, a.builder.parser.language)
	if err != nil {
		log.Error().Err(err).Msg("Error creating query executor")
		return Variables{}
	}
	for {
		match, ok := q.NextMatch()
		if !ok {
			break
		}
		for _, capture := range match.Captures {
			node := capture.Node
			variable := node.ChildByFieldName("left")
			if variable.Type() == PseudoVariableNodeType { //|| variable.Type() == PseudoVariableExpressionNodeType {
				pc := variable.NamedChild(0)
				if pc.Type() == PseudoContentNodeType {
					v := pc.NamedChild(0)
					switch v.Type() {
					case AVPNodeType:
						_id := v.ChildByFieldName("name").Child(0).Content(source_code)
						_name := _AVP_IDENTIFIER + "(" + _id + ")"
						_val := node.ChildByFieldName("right").Content(source_code)
						addAVPVariable(_name, _val, _id, Position{
							start: node.StartPoint(),
							end:   node.EndPoint(),
						})
					case VARNodeType:
						_id := v.ChildByFieldName("name").Child(0).Content(source_code)
						_name := _VAR_IDENTIFIER + "(" + _id + ")"
						_val := node.ChildByFieldName("right").Content(source_code)
						_scope := _VAR_SCOPE
						addLocalVariable(_name, _val, _scope, _id, Position{
							start: node.StartPoint(),
							end:   node.EndPoint(),
						})
					case DlgVarNodeType:
						_id := v.ChildByFieldName("name").Child(0).Content(source_code)
						_name := _DLG_VAR_IDENTIFIER + "(" + _id + ")"
						_val := node.ChildByFieldName("right").Content(source_code)
						_scope := _DLG_SCOPE
						addDlgVariable(_name, _val, _scope, _id, Position{
							start: node.StartPoint(),
							end:   node.EndPoint(),
						})
					default:
						continue
					}
				}
			}
		}
	}
	return Variables{
		avpVariables:   avpVariables,
		localVariables: localVariables,
		dlgVariables:   dlgVariables,
	}
}
