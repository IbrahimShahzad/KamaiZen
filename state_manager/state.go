package state_manager

import (
	"KamaiZen/kamailio_cfg"
	"KamaiZen/logger"
	"KamaiZen/lsp"
	"KamaiZen/settings"
	"fmt"
)

type State struct {
	Documents map[lsp.DocumentURI]string // A map of document URIs to their corresponding text content.
	Analyzer  *kamailio_cfg.Analyzer     // The analyzer used for parsing and analyzing the documents.
}

var state State

// GetState returns the current state.
//
// Returns:
//
//	*State - The current state.
func GetState() *State {
	return &state
}

// SetState sets the current state to the given state.
//
// Parameters:
//
//	s State - The new state to be set.
func SetState(s State) {
	state = s
}

// updateState updates the state with the given document URI and text.
// It also reinitializes the analyzer.
//
// Parameters:
//
//	DocumentURI lsp.DocumentURI - The URI of the document to be updated.
//	text string - The new text content of the document.
func (s *State) updateState(DocumentURI lsp.DocumentURI, text string) {
	state.Documents[DocumentURI] = text
	s.Analyzer = kamailio_cfg.NewAnalyzer()
}

// InitializeState initializes and returns a new state.
// It creates a new state, initializes the analyzer, and logs the initialization.
//
// Returns:
//
//	State - The initialized state.
func InitializeState() State {
	state = NewState()
	state.Analyzer = kamailio_cfg.NewAnalyzer()
	logger.Debugf("Parser: %v", state.Analyzer.GetParser())
	logger.Debug("State initialized")
	kamailio_cfg.InitialiseVariables()
	return state
}

// NewState creates and returns a new instance of State.
// It initializes the Documents map.
//
// Returns:
//
//	State - The initialized state.
func NewState() State {
	return State{
		Documents: make(map[lsp.DocumentURI]string),
	}
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
func (s *State) GetDocument(uri lsp.DocumentURI) string {
	return s.Documents[uri]
}

// SetDocument sets the text content of the document with the given URI.
//
// Parameters:
//
//	uri lsp.DocumentURI - The URI of the document.
//	text string - The new text content of the document.
func (s *State) SetDocument(uri lsp.DocumentURI, text string) {
	s.Documents[uri] = text
}

func (s *State) RegisterSubscribers() {
	// register the subscribers for the events
	// subscriber may include the Parser to update the AST
	// TODO: this goes along with the NotifySubscribers method
}

func (s *State) NotifySubsrcibers(uri lsp.DocumentURI, changes []lsp.TextDocumentContentChangeEvent) {
	// notify the subscribers of the changes
	// the parser and the analyser will be listening to the changes
	// analyzer_channel <- s.ChangeDocument(uri, changes)
	// TODO: think about the implementation
}

// OpenDocument opens the document with the given URI and text, and returns the diagnostics.
//
// Parameters:
//
//	uri lsp.DocumentURI - The URI of the document.
//	text string - The text content of the document.
//
// Returns:
//
//	[]lsp.Diagnostic - The list of diagnostics.
func (s *State) OpenDocument(uri lsp.DocumentURI, text string) []lsp.Diagnostic {
	s.Documents[uri] = text
	s.Analyzer.Build([]byte(text))

	visitor := kamailio_cfg.NewDiagnosticVisitor()
	s.Analyzer.GetAST().Accept(visitor, s.Analyzer)
	kamailio_cfg.ExtractVariables(s.Analyzer, []byte(text))
	visitor.GetQueryDiagnostics(s.Analyzer.GetAST(), s.Analyzer)
	if settings.GlobalSettings.EnableDiagnostics {
		return visitor.GetDiagnostics()
	}
	return []lsp.Diagnostic{}
}

// UpdateDocument updates the document with the given URI and text, and returns the diagnostics.
//
// Parameters:
//
//	uri lsp.DocumentURI - The URI of the document.
//	text string - The new text content of the document.
//
// Returns:
//
//	[]lsp.Diagnostic - The list of diagnostics.
func (s *State) UpdateDocument(uri lsp.DocumentURI, text string) []lsp.Diagnostic {
	s.updateState(uri, text)
	// for now we will parse the whole document
	s.Analyzer.Build([]byte(text))
	visitor := kamailio_cfg.NewDiagnosticVisitor()
	s.Analyzer.GetAST().Accept(visitor, s.Analyzer)
	kamailio_cfg.ExtractVariables(s.Analyzer, []byte(text))
	visitor.GetQueryDiagnostics(s.Analyzer.GetAST(), s.Analyzer)
	return visitor.GetDiagnostics()
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
func (s *State) Hover(id int, uri lsp.DocumentURI, position lsp.Position) lsp.HoverResponse {
	return lsp.NewHoverResponse(id,
		fmt.Sprintf("%s", GetNodeDocsAtPosition(uri, position, []byte(s.Documents[uri]))))
}

// Definition returns the definition information for the given document URI and position.
//
// Parameters:
//
//	id int - The ID of the definition request.
//	uri lsp.DocumentURI - The URI of the document.
//	position lsp.Position - The position within the document.
//
// Returns:
//
//	lsp.DefinitionProviderResponse - The definition response.
func (s *State) Definition(id int, uri lsp.DocumentURI, position lsp.Position) lsp.DefinitionProviderResponse {
	// TODO: Implement definition
	// lookup the text content type from the analysis and return the definition
	return lsp.NewDefintionProviderResponse(id, fmt.Sprintf("File: %s :: Definition at line %d, character %d", uri, position.Line, position.Character))
}

// TextDocumentCompletion returns the completion items for the given document URI and position.
//
// Parameters:
//
//	id int - The ID of the completion request.
//	uri lsp.DocumentURI - The URI of the document.
//	position lsp.Position - The position within the document.
//
// Returns:
//
//	lsp.CompletionResponse - The completion response.
func (s *State) TextDocumentCompletion(id int, uri lsp.DocumentURI, position lsp.Position) lsp.CompletionResponse {
	logger.Debug("Completion request for document with URI: ", uri)
	items := GetCompletionItems(uri)
	return lsp.NewCompletionResponse(id, items)
}

func (s *State) Formatting(id int, uri lsp.DocumentURI, options lsp.FormattingOptions) lsp.DocumentFormattingResponse {
	// TODO: Implement formatting
	// visitor := kamailio_cfg.NewFormattingVisitor()
	// s.Analyzer.GetAST().Accept(visitor, s.Analyzer)
	// edits := visitor.GetEdits()
	// return lsp.NewDocumentFormattingResponse(id, edits)
	return lsp.NewDocumentFormattingResponse(id, []lsp.TextEdit{})
}
