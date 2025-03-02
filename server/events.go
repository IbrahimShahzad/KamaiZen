package server

import (
	"KamaiZen/file"
	"KamaiZen/file_state"
	"KamaiZen/lsp"
	"KamaiZen/settings"
	"encoding/json"

	"github.com/rs/zerolog/log"
)

const (
	MethodInitialize            = "initialize"
	MethodInitialized           = "initialized"
	MethodDidOpen               = "textDocument/didOpen"
	MethodDidChange             = "textDocument/didChange"
	MethodHover                 = "textDocument/hover"
	MethodDefinition            = "textDocument/definition"
	MethodFormatting            = "textDocument/formatting"
	MethodCompletion            = "textDocument/completion"
	MethodConfiguration         = "workspace/Configuration"
	MethodConfigurationResponse = ""
)

// EventManager manages event handlers for different methods.
type EventManager struct {
	handlers map[string]func(contents []byte)
}

// NewEventManager creates and returns a new EventManager instance.
func NewEventManager() *EventManager {
	return &EventManager{
		handlers: make(map[string]func(contents []byte)),
	}
}

// RegisterHandler registers a handler function for a specific method.
// method: The name of the method for which the handler is being registered.
// handler: The function to handle the event. It takes an state_manager.State, the contents as a byte slice, and a channel for state_manager.State.
func (em *EventManager) RegisterHandler(method string, handler func(contents []byte)) {
	log.Info().Str("method", method).Msg("Registering handler")
	em.handlers[method] = handler
}

// Dispatch calls the registered handler for the given method.
// method: The name of the method for which the handler is being dispatched.
// state: The current state of the state_manager.
// contents: The contents to be passed to the handler as a byte slice.
// analyser_channel: A channel for state_manager.State to communicate with the handler.
func (em *EventManager) Dispatch(method string, contents []byte) {
	if handler, found := em.handlers[method]; found {
		handler(contents)
		return
	}
	log.Error().Str("method", method).Msg("No handler found")
}

// handleInitialized handles the 'initialized' notification.
// state: The current state of the state_manager.
// contents: The contents of the notification as a byte slice.
// analyser_channel: A channel for state_manager.State to communicate with the handler.
func handleInitialized(contents []byte) {
	var notification lsp.InitializedNotification
	log.Info().Str("contents", string(contents)).Msg("Received initialized notification")
	if e := json.Unmarshal(contents, &notification); e != nil {
		log.Error().Err(e).Msg("Error unmarshalling initialized notfication")
		return
	}
	GetServerInstance().SetState(ServerInitialized)
	log.Info().Msgf("Received initialized notification with %v", notification)
}

// handleInitialize handles the 'initialize' request.
// state: The current state of the state_manager.
// contents: The contents of the request as a byte slice.
// analyser_channel: A channel for state_manager.State to communicate with the handler.
func handleInitialize(contents []byte) {
	var request lsp.InitializeRequest
	if e := json.Unmarshal(contents, &request); e != nil {
		log.Error().Err(e).Msg("Error unmarshalling initialize request")
		return
	}
	log.Info().
		Str("client", request.Params.ClientInfo.Name).
		Str("version", request.Params.ClientInfo.Version).
		Msg("Connected... Sending workspace configuration request")
	config_request := lsp.NewWorkspaceConfigurationRequest(
		request.ID,
		lsp.ConfigurationParams{
			Items: []lsp.ConfigurationItem{
				{
					Section: "kamaizen",
				},
			},
		})
	GetServerInstance().SetState(ServerInitializing)
	lsp.WriteResponse(config_request)
}

func handleWorkspaceConfiguration(contents []byte) {
	var response lsp.WorkspaceConfigurationResponse
	if e := json.Unmarshal(contents, &response); e != nil {
		log.Error().Err(e).Msg("Error unmarshalling workspace configuration response")
		return
	}
	var initialize_response lsp.InitializeResponse
	initialize_response = lsp.NewInitializeResponse(response.ID)
	lsp.WriteResponse(initialize_response)
	GetServerInstance().addKamailioMethods(
		settings.NewLSPSettings(
			response.Result[0].KamailioSourcePath,
			"", // TODO: Add root path
			response.Result[0].Loglevel,
			response.Result[0].EnableDeprecatedCommentHint,
			response.Result[0].EnableDiagnostics))
}

// handleDidOpen handles the 'didOpen' notification.
// state: The current state of the state_manager.
// contents: The contents of the notification as a byte slice.
// analyser_channel: A channel for state_manager.State to communicate with the handler.
func handleDidOpen(contents []byte) {
	log.Error().Msg("Did open Diagnostic")
	var notification lsp.DidOpenTextDocumentNotification
	if e := json.Unmarshal(contents, &notification); e != nil {
		log.Error().Err(e).Msg("Error unmarshalling didOpen notification")
		return
	}
	uri := notification.Params.TextDocument.URI
	text := []byte(notification.Params.TextDocument.Text)

	log.Info().Str("uri", string(uri)).Msg("Opened document")
	server := GetServerInstance()
	server.fileStates.AddState(
		notification.Params.TextDocument.URI,
		file_state.NewState(
			file.NewFile(uri, text),
		),
	)

	server.diagnosticsMu.Lock()
	defer server.diagnosticsMu.Unlock()
	server.diagnostics[uri] = server.fileStates.OpenDocument(uri, string(text), server.diagnostics[uri])
	log.Error().Msgf("Diagnostics -------- count: %d %v", server.diagnostics[uri].GetDiagnosticsCount(), server.diagnostics[uri])
	if server.diagnostics[uri].GetDiagnosticsCount() > 0 && settings.GlobalSettings.EnableDiagnostics {
		lsp.WriteResponse(
			lsp.NewPublishDiagnosticNotification(
				notification.Params.TextDocument.URI,
				server.diagnostics[uri].GetDiagnostics(),
			),
		)
	}
}

// handleMessage handles incoming messages and dispatches them to the appropriate handler.
// state: The current state of the state_manager.
// method: The name of the method for which the handler is being dispatched.
// contents: The contents of the message as a byte slice.
// analyser_channel: A channel for state_manager.State to communicate with the handler.
// eventManager: The EventManager instance to use for dispatching the message.
func handleMessage(method string, contents []byte, eventManager *EventManager) {
	log.Info().Str("method", method).Msg("Received message with method")
	eventManager.Dispatch(method, contents)
}

// handleDidChange handles the 'didChange' notification.
// state: The current state of the state_manager.
// contents: The contents of the notification as a byte slice.
// analyser_channel: A channel for state_manager.State to communicate with the handler.
func handleDidChange(contents []byte) {
	var notification lsp.DidChangeTextDocumentNotification
	if e := json.Unmarshal(contents, &notification); e != nil {
		log.Error().Err(e).Msg("Error unmarshalling didChange notification")
		return
	}
	uri := notification.Params.TextDocument.URI
	server := GetServerInstance()
	for _, change := range notification.Params.ContentChanges {
		server.diagnosticsMu.Lock()
		server.diagnostics[uri].ClearDiagnostics()
		server.diagnostics[uri] = server.fileStates.UpdateDocument(uri, change.Text, server.diagnostics[uri])
		log.Error().Msgf("Diagnostics -------- count: %d %v", server.diagnostics[uri].GetDiagnosticsCount(), settings.GlobalSettings.EnableDiagnostics)
		defer server.diagnosticsMu.Unlock()
		if server.diagnostics[uri].GetDiagnosticsCount() > 0 &&
			settings.GlobalSettings.EnableDiagnostics {
			lsp.WriteResponse(
				lsp.NewPublishDiagnosticNotification(
					notification.Params.TextDocument.URI,
					server.diagnostics[uri].GetDiagnostics(),
				),
			)
			return
		}
		// clear diagnostics
		log.Debug().Str("uri", string(uri)).Msg("Clearing diagnostics for document")
		server.diagnosticsMu.Lock()
		server.diagnostics[uri].ClearDiagnostics()
		server.diagnosticsMu.Unlock()
		lsp.WriteResponse(
			lsp.NewPublishDiagnosticNotification(
				notification.Params.TextDocument.URI,
				[]lsp.Diagnostic{},
			),
		)

	}
}

// handleHover handles the 'hover' request.
// state: The current state of the state_manager.
// contents: The contents of the request as a byte slice.
// analyser_channel: A channel for state_manager.State to communicate with the handler.
func handleHover(contents []byte) {
	var request lsp.HoverRequest
	if e := json.Unmarshal(contents, &request); e != nil {
		log.Error().Err(e).Msg("Error unmarshalling hover request")
		return
	}
	lsp.WriteResponse(
		GetServerInstance().fileStates.Hover(
			request.ID,
			request.Params.TextDocument.URI,
			request.Params.Position,
		),
	)
}

// handleDefinition handles the 'definition' request.
// state: The current state of the state_manager.
// contents: The contents of the request as a byte slice.
// analyser_channel: A channel for state_manager.State to communicate with the handler.
func handleDefinition(contents []byte) {
	var request lsp.DefinitionProviderRequest
	if e := json.Unmarshal(contents, &request); e != nil {
		log.Error().Err(e).Msg("Error unmarshalling definition request")
		return
	}
	lsp.WriteResponse(
		GetServerInstance().fileStates.Definition(
			request.ID,
			request.Params.TextDocument.URI,
			request.Params.Position,
		),
	)
}

// handleFormatting handles the 'formatting' request.
// state: The current state of the state_manager.
// contents: The contents of the request as a byte slice.
// analyser_channel: A channel for state_manager.State to communicate with the handler.
func handleFormatting(contents []byte) {
	var request lsp.DocumentFormattingRequest
	if e := json.Unmarshal(contents, &request); e != nil {
		log.Error().Err(e).Msg("Error unmarshalling formatting request")
		return
	}
	lsp.WriteResponse(
		GetServerInstance().fileStates.Formatting(
			request.ID,
			request.Params.TextDocument.URI,
			request.Params.Options,
		),
	)
}

// handleCompletion handles the 'completion' request.
// state: The current state of the state_manager.
// contents: The contents of the request as a byte slice.
// analyser_channel: A channel for state_manager.State to communicate with the handler.
func handleCompletion(contents []byte) {
	var request lsp.CompletionRequest
	if e := json.Unmarshal(contents, &request); e != nil {
		log.Error().Err(e).Msg("Error unmarshalling completion request")
		return
	}
	lsp.WriteResponse(
		GetServerInstance().fileStates.TextDocumentCompletion(
			request.ID,
			request.Params.TextDocument.URI,
			request.Params.Position,
		),
	)
}
