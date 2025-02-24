package lsp

// DidChangeTextDocumentNotification represents a notification sent to the server
// when a text document is changed. It contains the notification metadata and the
// parameters for the change.
type DidChangeTextDocumentNotification struct {
	Notification
	Params DidChangeTextDocumentParams `json:"params"`
}

// DidChangeTextDocumentParams contains the parameters for the DidChangeTextDocumentNotification.
// It includes the versioned text document identifier and the list of content changes.
type DidChangeTextDocumentParams struct {
	TextDocument   VersionedTextDocumentIdentifier  `json:"textDocument"`
	ContentChanges []TextDocumentContentChangeEvent `json:"contentChanges"`
}

// TextDocumentContentChangeEvent represents a change event in a text document.
// It includes the range of the change and the new text content.
type TextDocumentContentChangeEvent struct {
	Range *Range `json:"range,omitempty"`
	Text  string `json:"text"`
}

// Apply applies the content change event to the given text.
// It modifies the text based on the range and new content specified in the change event.
//
// Parameters:
//
//	text string - The original text to be modified.
//
// Returns:
//
//	string - The modified text after applying the change event.
func (change TextDocumentContentChangeEvent) Apply(text string) string {
	if change.Range != nil {
		start := change.Range.Start
		end := change.Range.End
		text = text[:start.Line] + change.Text + text[end.Line:]
	} else {
		text = change.Text
	}
	return text
}
