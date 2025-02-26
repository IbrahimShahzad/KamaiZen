package lsp

import "KamaiZen/settings"

// InitializeRequest represents a request to initialize the language server.
// It contains the request metadata and the parameters for initialization.
type InitializeRequest struct {
	Request
	Params InitializeRequestParams `json:"params"`
}

// InitializeRequestParams contains the parameters for the InitializeRequest.
// It includes information about the client.
type InitializeRequestParams struct {
	ClientInfo ClientInfo `json:"clientInfo"`
}

// ClientInfo represents information about the client making the request.
// It includes the client's name and version.
type ClientInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

// InitializeResponse represents the response to an InitializeRequest.
// It contains the response metadata and the result of the initialization.
type InitializeResponse struct {
	Response
	Result InitializeResult `json:"result"`
}

// InitializeResult contains the result of the initialization.
// It includes the server's capabilities and information about the server.
type InitializeResult struct {
	Capabilities ServerCapabilities `json:"capabilities"`
	ServerInfo   ServerInfo         `json:"serverInfo"`
}

// ServerCapabilities represents the capabilities of the language server.
// It includes various features supported by the server.
type ServerCapabilities struct {
	// TODO: incremental instead of full
	TextDocumentSync           TextDocumentSyncOptions `json:"textDocumentSync"`
	HoverProvider              bool                    `json:"hoverProvider"`
	DefinitionProvider         bool                    `json:"definitionProvider"`
	DocumentFormattingProvider bool                    `json:"documentFormattingProvider"`
	CompletionProvider         map[string]any          `json:"completionProvider"`
	DocumentHighlightProvider  bool                    `json:"documentHighlightProvider"`
	// TODO: Add more capabilities
	// CodeActionProvider bool `json:"codeActionProvider"`
}

// ServerInfo represents information about the language server.
// It includes the server's name and version.
type ServerInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

// NewInitializeResponse creates and returns a new InitializeResponse.
// It initializes the response with the given ID and sets the server capabilities and information.
//
// Parameters:
//
//	id int - The ID of the response.
//
// Returns:
//
//	InitializeResponse - The initialized response.
func NewInitializeResponse(id int) InitializeResponse {
	return InitializeResponse{
		Response: Response{
			RPC: "2.0",
			ID:  id,
		},
		Result: InitializeResult{
			Capabilities: ServerCapabilities{
				TextDocumentSync: TextDocumentSyncOptions{
					OpenClose: true,
					Change:    TEXT_DOCUMENT_SYNC_KIND_FULL,
					// FIXME: update this to incremental
					// Change:    TEXT_DOCUMENT_SYNC_KIND_INCREMENTAL,
				},
				HoverProvider:      true,
				DefinitionProvider: true,
				// FIXME: Update to a proper formatter
				DocumentFormattingProvider: true,
				CompletionProvider:         map[string]any{"resolveProvider": false},
				DocumentHighlightProvider:  false,
			},
			ServerInfo: ServerInfo{
				Name:    settings.MY_NAME,
				Version: settings.KAMAIZEN_VERSION,
			},
		},
	}
}
