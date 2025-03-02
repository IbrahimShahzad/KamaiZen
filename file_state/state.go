package file_state

import (
	"KamaiZen/file"
	"KamaiZen/kamailio_cfg"
	"KamaiZen/lsp"
	"fmt"

	"github.com/rs/zerolog/log"
)

// FileState represents the current state of the file and the analyzer associated.
type FileState struct {
	f        *file.File              // The file associated with the state.
	analyser *kamailio_cfg.Analyzer  // The analyzer used for parsing and analyzing the documents.
	vars     *kamailio_cfg.Variables // vairables associated with the file
	funcs    kamailio_cfg.Functions  // functions associated with the file
	routes   kamailio_cfg.Routes     // routes associated with the file
	isRoot   bool                    // we are going to consider the file containing the request_route as the root file
}

var state FileState

// GetState returns the current state.
//
// Returns:
//
//	*State - The current state.
func GetState() *FileState {
	return &state
}

// SetState sets the current state to the given state.
//
// Parameters:
//
//	s State - The new state to be set.
func SetState(s FileState) {
	state = s
}

// updateState updates the state with the given document URI and text.
// It also reinitializes the analyzer.
//
// Parameters:
//
//	DocumentURI lsp.DocumentURI - The URI of the document to be updated.
//	text string - The new text content of the document.
func (s *FileState) updateState(text []byte) {
	s.f.SetContent(text)
	s.analyser.Renew()
}

// NewState creates and returns a new instance of State.
// It initializes the Documents map.
//
// Returns:
//
//	State - The initialized state.
func NewState(f *file.File) *FileState {
	return &FileState{
		f:        f,
		analyser: kamailio_cfg.NewAnalyzer(),
		vars:     kamailio_cfg.NewVariables(),
		funcs:    kamailio_cfg.NewFunctions(),
	}
}

func (s *FileState) GetAnalyzer() *kamailio_cfg.Analyzer {
	return s.analyser
}

// GetDocument returns the text content of the document with the given URI.
//
// Parameters:
//
//	uri lsp.DocumentURI - The URI of the document.
//
// Returns:
//
//	string - The text content of the document.
func (s *FileState) GetDocument() []byte {
	return s.f.Content()
}

// SetDocument sets the text content of the document with the given URI.
//
// Parameters:
//
//	text string - The new text content of the document.
func (s *FileState) SetDocument(text []byte) {
	s.f.SetContent([]byte(text))
}

// OpenDocument opens the document with the given URI and text, and returns the diagnostics.
//
// Parameters:
//
//	content []byte - The text content of the document.
func (s *FileState) openDocument(content []byte, d *kamailio_cfg.DiagnosticVisitor) {
	s.SetDocument(content)
	s.analyser.Build(content)
	s.analyser.GetAST().Accept(d, s.analyser)
	vars := kamailio_cfg.ExtractVariables(s.analyser, []byte(content))
	s.vars = &vars
	s.routes = kamailio_cfg.FetchRoutes(s.analyser, []byte(content))
	if s.routes.ContainsRequestRoute() {
		s.isRoot = true
	}
	d.GetQueryDiagnostics(s.analyser)
}

// UpdateDocument updates the document with the given content, and returns the diagnostics.
//
// Parameters:
//
//	content []byte - The new text content of the document.
//
// Returns:
//
//	[]lsp.Diagnostic - The list of diagnostics.
func (s *FileState) updateDocument(content []byte, d *kamailio_cfg.DiagnosticVisitor) {
	s.updateState(content)
	// for now we will parse the whole document
	s.analyser.Build(content)
	s.analyser.GetAST().Accept(d, s.analyser)
	vars := kamailio_cfg.ExtractVariables(s.analyser, []byte(content))
	s.vars = &vars
	s.routes = kamailio_cfg.FetchRoutes(s.analyser, []byte(content))
	if s.routes.ContainsRequestRoute() {
		s.isRoot = true
	}
	d.GetQueryDiagnostics(s.analyser)
}

// Hover returns the hover information for the given document URI and position.
//
// Parameters:
//
//	id int - The ID of the hover request.
//	uri lsp.DocumentURI - The URI of the document.
//	position lsp.Position - The position within the document.
//
// Returns:
//
//	lsp.HoverResponse - The hover response.
func (s *FileState) hover(id int, position lsp.Position) lsp.HoverResponse {
	return lsp.NewHoverResponse(id,
		fmt.Sprintf("%s", s.GetNodeDocsAtPosition(position, s.f.Content())))
}

// Definition returns the definition information for the given document URI and position.
//
// Parameters:
//
//	id int - The ID of the definition request.
//	position lsp.Position - The position within the document.
//
// Returns:
//
//	lsp.DefinitionProviderResponse - The definition response.
func (s *FileState) definition(
	id int,
	position lsp.Position,
) lsp.DefinitionProviderResponse {
	log.Error().Msg("Definition not implemented")
	r := s.GetRouteDefinitionAtPosition(position)
	log.Info().Msgf("Route definition: %v", r)
	if r == nil {
		return lsp.NewDefintionProviderResponse(
			id,
			"No definition found",
			nil,
		)
	}
	return lsp.NewDefintionProviderResponse(
		id,
		r.Content,
		&lsp.Location{
			URI: s.f.URI,
			Range: lsp.Range{
				Start: lsp.Position{
					Line:      int(r.StartPoint.Row),
					Character: int(r.StartPoint.Column),
				},
				End: lsp.Position{
					Line:      int(r.EndPoint.Row),
					Character: int(r.EndPoint.Column),
				},
			},
		},
	)
}

// TextDocumentCompletion returns the completion items for the given document URI and position.
//
// Parameters:
//
//	id int - The ID of the completion request.
//	position lsp.Position - The position within the document.
//
// Returns:
//
//	lsp.CompletionResponse - The completion response.
func (s *FileState) textDocumentCompletion(id int, position lsp.Position) lsp.CompletionResponse {
	return lsp.NewCompletionResponse(id, s.GetCompletionItems())
}

func (s *FileState) formatting(id int, options lsp.FormattingOptions) lsp.DocumentFormattingResponse {
	// TODO: Implement formatting
	// visitor := kamailio_cfg.NewFormattingVisitor()
	// s.Analyzer.GetAST().Accept(visitor, s.Analyzer)
	// edits := visitor.GetEdits()
	// return lsp.NewDocumentFormattingResponse(id, edits)
	log.Info().Msg("Formatting document")
	new_text := kamailio_cfg.FixIndent(string(s.f.Content()))
	// TODO:
	// we can use this text to update the parse tree
	// and then visit the tree to get further formatting
	s.analyser.Build([]byte(new_text[0].NewText))
	return lsp.NewDocumentFormattingResponse(id, kamailio_cfg.FixIndent(string(s.f.Content())))
}
