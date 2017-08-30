package database

import "errors"

var (
	// ErrDuplicated .
	ErrDuplicated = errors.New("Duplicated entry")
	// ErrNotFound .
	ErrNotFound = errors.New("Not found")
)
