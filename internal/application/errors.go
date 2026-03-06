package application

import "errors"

var (
	ErrFileTooLarge       = errors.New("file too large")
	ErrEmptyText          = errors.New("empty text")
	ErrUnsupportedEncoding = errors.New("unsupported encoding")
	ErrNoTerms            = errors.New("no terms found after filtering")
	ErrAnalysisCancelled  = errors.New("analysis cancelled")
)
