package lsp

import "KamaiZen/settings"

// DefinitionProviderRequest represents a request for definition information.
// It contains the request metadata and the parameters for the definition request.
type DefinitionProviderRequest struct {
	Request
	Params DefinitionProviderParams `json:"params"`
}

// DefinitionProviderParams contains the parameters for the DefinitionProviderRequest.
// It includes the text document position parameters.
type DefinitionProviderParams struct {
	TextDocuemntPositionParams
}

// DefinitionProviderResponse represents the response to a DefinitionProviderRequest.
// It contains the response metadata and the location of the definition.
type DefinitionProviderResponse struct {
	Response
	Result Location `json:"result"`
}

// DefinitionProvider represents a provider for definition information.
// It includes the contents of the definition and the range within the document.
type DefinitionProvider struct {
	Contents MarkupContent `json:"contents"`
	Range    *Range        `json:"range,omitempty"`
}

// NewDefinitionProviderResponse creates and returns a new DefinitionProviderResponse.
// It initializes the response with the given ID and sets the location of the definition.
//
// Parameters:
//
//	id int - The ID of the response.
//	contents string - The contents of the definition.
//
// Returns:
//
//	DefinitionProviderResponse - The initialized response.
func NewDefintionProviderResponse(id int, contents string, uri DocumentURI, start Position, end Position) DefinitionProviderResponse {
	return DefinitionProviderResponse{
		Response: Response{
			RPC: settings.RPC_VERSION,
			ID:  id,
		},
		Result: Location{
			URI: uri,
			Range: Range{
				Start: start,
				End:   end,
			},
		},
	}
}
