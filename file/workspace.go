package file

import (
	"KamaiZen/lsp"
	"sync"
)

type Folder = lsp.DocumentURI

type Workspace struct {
	mu          sync.Mutex                // protects the fields below
	files       map[lsp.DocumentURI]*File // the files in the workspace
	initialized bool                      // true if the workspace has been initialized
	root        Folder                    // the root of the workspace
}

// NewWorkspace creates and returns a new Workspace instance.
func NewWorkspace() Workspace {
	return Workspace{
		files:       make(map[lsp.DocumentURI]*File),
		initialized: false,
	}
}
