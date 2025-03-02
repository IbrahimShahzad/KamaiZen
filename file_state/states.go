package file_state

import (
	"KamaiZen/kamailio_cfg"
	"KamaiZen/lsp"
)

type FileStates struct {
	States map[lsp.DocumentURI]*FileState
}

// NewFileStates creates and returns a new instance of FileStates.
func NewFileStates() *FileStates {
	return &FileStates{
		States: make(map[lsp.DocumentURI]*FileState),
	}
}

// AddState adds the state of the document with the given URI.
func (fs *FileStates) AddState(uri lsp.DocumentURI, state *FileState) {
	fs.States[uri] = state
}

// GetState returns the state of the document with the given URI.
func (fs *FileStates) GetState(uri lsp.DocumentURI) *FileState {
	return fs.States[uri]
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
func (fs *FileStates) Definition(
	id int,
	uri lsp.DocumentURI,
	position lsp.Position,
) lsp.DefinitionProviderResponse {
	return fs.GetState(uri).definition(id, position)
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
func (s *FileStates) Hover(id int, uri lsp.DocumentURI, position lsp.Position) lsp.HoverResponse {
	return s.GetState(uri).hover(id, position)
}

func (s *FileStates) Formatting(id int, uri lsp.DocumentURI, options lsp.FormattingOptions) lsp.DocumentFormattingResponse {
	return s.GetState(uri).formatting(id, options)
}

func (s *FileStates) TextDocumentCompletion(id int, uri lsp.DocumentURI, position lsp.Position) lsp.CompletionResponse {
	return s.GetState(uri).textDocumentCompletion(id, position)
}

// OpenDocument opens the document with the given URI and text, and returns the diagnostics.
//
// Parameters:
//
//	uri lsp.DocumentURI - The URI of the document.
//	text string - The text content of the document.
func (s *FileStates) OpenDocument(uri lsp.DocumentURI, text string, d *kamailio_cfg.DiagnosticVisitor) *kamailio_cfg.DiagnosticVisitor {
	d = kamailio_cfg.NewDiagnosticVisitor()
	s.GetState(uri).openDocument([]byte(text), d)
	return d
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
func (s *FileStates) UpdateDocument(uri lsp.DocumentURI, text string, d *kamailio_cfg.DiagnosticVisitor) *kamailio_cfg.DiagnosticVisitor {
	d = kamailio_cfg.NewDiagnosticVisitor()
	s.GetState(uri).updateDocument([]byte(text), d)
	return d
}
