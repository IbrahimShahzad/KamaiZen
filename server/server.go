package server

import (
	"KamaiZen/document_manager"
	"KamaiZen/file"
	sm "KamaiZen/file_state"
	"KamaiZen/kamailio_cfg"
	"KamaiZen/lsp"
	"KamaiZen/rpc"
	"KamaiZen/settings"
	"bufio"
	"os"
	"sync"

	"github.com/rs/zerolog/log"
)

type ServerState int

const (
	ServerCreated ServerState = iota
	ServerInitializing
	ServerInitialized
	ServerShutDown
)

func (s ServerState) String() string {
	return [...]string{"ServerCreated", "ServerInitializing", "ServerInitialized", "ServerShutDown"}[s]
}

func (s *ServerState) setState(newState ServerState) {
	s = &newState
}

type Server struct {
	eventManager *EventManager
	stateMu      sync.Mutex
	state        *ServerState

	worksapce  file.Workspace
	fileStates *sm.FileStates

	diagnosticsMu         sync.Mutex // guards map and its values
	diagnostics           map[lsp.DocumentURI]*kamailio_cfg.DiagnosticVisitor
	cancelPrevDiagnostics func()
}

func (s *Server) SetState(newState ServerState) {
	s.stateMu.Lock()
	defer s.stateMu.Unlock()
	s.state.setState(newState)
}

func IsServerInitialized() bool {
	server := GetServerInstance()
	server.stateMu.Lock()
	defer server.stateMu.Unlock()
	return *server.state == ServerInitialized
}

// create a single instance of the server
var _serverInstance *Server

func NewServerInstance() *Server {
	if _serverInstance == nil {
		_serverInstance = &Server{
			eventManager: NewEventManager(),
			state:        new(ServerState),
			diagnostics:  make(map[lsp.DocumentURI]*kamailio_cfg.DiagnosticVisitor),
			worksapce:    file.NewWorkspace(),
			fileStates:   sm.NewFileStates(),
		}
		_serverInstance.state.setState(ServerCreated)
	}
	return _serverInstance
}

// GetServerInstance returns the single instance of the server.
func GetServerInstance() *Server {
	return _serverInstance
}

// StartServer starts the language server and listens for incoming messages from the client.
// It initializes the event manager, registers handlers for various methods, and processes incoming messages.
//
// Parameters:
func (s *Server) StartServer(wg *sync.WaitGroup) {
	defer wg.Done()
	scanner := bufio.NewScanner(os.Stdin)
	log.Info().Msg("Starting server")
	scanner.Split(rpc.Split)

	// Initialize EventManager and register handlers
	s.RegisterDefaultHandlers()

	for scanner.Scan() {
		msg := scanner.Bytes()
		method, contents, e := rpc.DecodeMessage(msg)
		if e != nil {
			log.Error().Err(e).Msg("Error decoding message")
			continue
		}
		handleMessage(method, contents, s.eventManager)
	}
}

func (s *Server) RegisterDefaultHandlers() {
	s.RegisterHandler(MethodInitialize, handleInitialize)
	s.RegisterHandler(MethodInitialized, handleInitialized)
	s.RegisterHandler(MethodDidOpen, handleDidOpen)
	s.RegisterHandler(MethodDidChange, handleDidChange)
	s.RegisterHandler(MethodDefinition, handleDefinition)
	s.RegisterHandler(MethodFormatting, handleFormatting)
	s.RegisterHandler(MethodConfigurationResponse, handleWorkspaceConfiguration)
}

func (s *Server) StopServer() {
	log.Info().Msg("Stopping server")
}

func (s *Server) RegisterHandler(method string, handler func(contents []byte)) {
	s.eventManager.RegisterHandler(method, handler)
}

func (s *Server) addKamailioMethods(settings settings.LSPSettings) {
	log.Debug().Str("path", settings.KamailioSourcePath).Msg("Kamailio src added")
	log.Debug().Msg("Adding Hover and Completion methods")
	document_manager.Initialise(settings)
	s.RegisterHandler(MethodHover, handleHover)
	s.RegisterHandler(MethodCompletion, handleCompletion)
}
