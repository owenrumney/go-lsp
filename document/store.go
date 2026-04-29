package document

import (
	"fmt"
	"sort"
	"sync"

	"github.com/owenrumney/go-lsp/lsp"
)

// Store tracks open text documents and applies LSP document sync messages.
type Store struct {
	mu   sync.RWMutex
	docs map[lsp.DocumentURI]*Document
}

// NewStore creates an empty document store.
func NewStore() *Store {
	return &Store{docs: make(map[lsp.DocumentURI]*Document)}
}

// Open records a newly opened document.
func (s *Store) Open(params *lsp.DidOpenTextDocumentParams) (*Document, error) {
	doc := newDocument(params.TextDocument)

	s.mu.Lock()
	defer s.mu.Unlock()
	s.docs[doc.URI()] = doc
	return doc.snapshot(), nil
}

// Change applies all content changes from a didChange notification.
func (s *Store) Change(params *lsp.DidChangeTextDocumentParams) (*Document, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	doc, ok := s.docs[params.TextDocument.URI]
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrDocumentNotFound, params.TextDocument.URI)
	}

	for _, change := range params.ContentChanges {
		if err := doc.ApplyChange(change, params.TextDocument.Version); err != nil {
			return nil, err
		}
	}

	return doc.snapshot(), nil
}

// Close removes a document from the store.
func (s *Store) Close(params *lsp.DidCloseTextDocumentParams) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.docs, params.TextDocument.URI)
}

// Get returns a snapshot of an open document.
func (s *Store) Get(uri lsp.DocumentURI) (*Document, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	doc, ok := s.docs[uri]
	if !ok {
		return nil, false
	}
	return doc.snapshot(), true
}

// Text returns the full text for an open document.
func (s *Store) Text(uri lsp.DocumentURI) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	doc, ok := s.docs[uri]
	if !ok {
		return "", false
	}
	return doc.text, true
}

// Version returns the current version for an open document.
func (s *Store) Version(uri lsp.DocumentURI) (int, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	doc, ok := s.docs[uri]
	if !ok {
		return 0, false
	}
	return doc.version, true
}

// URIs returns all open document URIs in lexical order.
func (s *Store) URIs() []lsp.DocumentURI {
	s.mu.RLock()
	defer s.mu.RUnlock()

	uris := make([]lsp.DocumentURI, 0, len(s.docs))
	for uri := range s.docs {
		uris = append(uris, uri)
	}
	sort.Slice(uris, func(i, j int) bool { return uris[i] < uris[j] })
	return uris
}

func (d *Document) snapshot() *Document {
	cp := *d
	cp.lineStarts = append([]int(nil), d.lineStarts...)
	return &cp
}
