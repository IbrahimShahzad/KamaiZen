package settings

type LSPSettings struct {
	KamailioSourcePath     string `json:"kamailioSourcePath"`
	LogLevel               int    `json:"logLevel"`
	DeprecatedCommentHints bool   `json:"deprecatedCommentHints"`
	EnableDiagnostics      bool   `json:"enableDiagnostics"`
}

var GlobalSettings LSPSettings

// NewLSPSettings creates and returns a new instance of LSPSettings.
// It initializes the settings with the given Kamailio source path, root directory, and log level.
//
// Parameters:
//
//	kamailioSourcePath string - The path to the Kamailio source code.
//	rootDir string - The root directory for the language server.
//	log_level logger.LOGLEVEL - The logging level for the language server.
//	dch - Deprecated Comments Hints enabled/disabled
//	diag - Diagnostics enabled/disabled
//
// Returns:
//
//	LSPSettings - The initialized settings.
func NewLSPSettings(kamailioSourcePath string, rootDir string, log_level int, dch bool, diag bool) LSPSettings {
	GlobalSettings = LSPSettings{
		KamailioSourcePath:     kamailioSourcePath,
		LogLevel:               log_level,
		DeprecatedCommentHints: dch,
		EnableDiagnostics:      diag,
	}
	return GlobalSettings
}

const RPC_VERSION = "2.0"
const KAMAIZEN_VERSION = "0.1.0"
const MY_NAME = "KamaiZen"
