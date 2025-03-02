package file

import (
	"KamaiZen/lsp"
	"bytes"
	"crypto/sha256"
	"fmt"
	"io"
	"sync"
)

type File struct {
	sync.Mutex                 // protects the fields below
	URI        lsp.DocumentURI // the file's URI
	content    []byte          // the file's content
	hash       Hash            // the file's hash
	saved      bool            // true if the file has been saved
	hashOnSave Hash            // the file's hash when it was last saved
}

// SetContent sets the content of the file.
func (f *File) SetContent(content []byte) {
	f.Lock()
	defer f.Unlock()
	f.content = content
}

// Content returns the content of the file.
func (f *File) Content() []byte {
	f.Lock()
	defer f.Unlock()
	return f.content
}

// Hash returns the hash of the file.
func (f *File) Hash() Hash {
	f.Lock()
	defer f.Unlock()
	return f.hash
}

// HasChangedSinceSave returns true if the file has changed since it was last saved.
func (f *File) HasChangedSinceSave() bool {
	f.Lock()
	defer f.Unlock()
	return f.hash != f.hashOnSave
}

// Saved returns true if the file has been saved.
func (f *File) Saved() bool {
	f.Lock()
	defer f.Unlock()
	return f.saved
}

// HasChanged returns true if the file has changed since the given hash.
func (f *File) HasChanged(hash Hash) bool {
	f.Lock()
	defer f.Unlock()
	if f.hash != hash {
		f.saved = false
		return true
	}
	return false
}

// Save saves the file.
func (f *File) Save() {
	f.Lock()
	defer f.Unlock()
	f.hashOnSave = f.hash
	f.saved = true
}

// NewFile creates a new file with the given URI and content.
func NewFile(uri lsp.DocumentURI, content []byte) *File {
	// create hash
	hash, err := calculateFileHash(content)
	if err != nil {
		panic("I will die right here right now because I can't calculate the hash of a file")
	}
	return &File{
		URI:     uri,
		content: content,
		hash:    hash,
		saved:   true,
	}
}

type Hash [sha256.Size]byte

// calculateFileHash calculates the hash of the given data.
func calculateFileHash(data []byte) (Hash, error) {
	hasher := sha256.New()
	_, err := io.Copy(hasher, bytes.NewReader(data))
	if err != nil {
		return Hash{}, err
	}
	return Hash(hasher.Sum(nil)), nil
}

// String returns the digest as a string of hex digits.
func (h Hash) String() string {
	return fmt.Sprintf("%64x", [sha256.Size]byte(h))
}
