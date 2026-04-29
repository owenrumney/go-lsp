package document

import "errors"

var (
	// ErrDocumentNotFound means a change or lookup referenced an unopened document.
	ErrDocumentNotFound = errors.New("document not found")

	// ErrInvalidPosition means an LSP position does not point into the document.
	ErrInvalidPosition = errors.New("invalid document position")

	// ErrInvalidRange means an LSP range is outside the document or has start after end.
	ErrInvalidRange = errors.New("invalid document range")

	// ErrVersionRegression means an update tried to move a document version backwards.
	ErrVersionRegression = errors.New("document version regression")
)
