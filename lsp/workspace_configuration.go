package lsp

type ConfigurationParams struct {
	Items []ConfigurationItem `json:"items"`
}

type ConfigurationItem struct {
	ScopeUri string `json:"scopeUri"`
	Section  string `json:"section"`
}

// type ConfigurationResponse struct {
// 	// it can be either	error or ConfigrationObject
// 	Result interface{} `json:"result"`
// }

type ConfigurationResponse struct {
	Result []ConfigurationObject `json:"result"`
}

type ConfigurationObject struct {
	KamailioSourcePath          string `json:"kamailioSourcePath"`
	Loglevel                    int    `json:"logLevel"`
	EnableDeprecatedCommentHint bool   `json:"enableDeprecatedCommentHint"`
	EnableDiagnostics           bool   `json:"enableDiagnostics"`
}

type ConfigurationItemValue struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type ConfigurationItemValueResponse struct {
	Value string `json:"value"`
}

type ConfigurationItemValueParams struct {
	Item ConfigurationItem `json:"item"`
}

type WorkspaceConfigurationRequest struct {
	Request
	Params ConfigurationParams `json:"params"`
}

type WorkspaceConfigurationResponse struct {
	Response
	Result []ConfigurationObject `json:"result"`
}

func NewWorkspaceConfigurationRequest(id int, params ConfigurationParams) WorkspaceConfigurationRequest {
	return WorkspaceConfigurationRequest{
		Request: Request{
			ID:     id,
			Method: "workspace/configuration",
		},
		Params: params,
	}
}
