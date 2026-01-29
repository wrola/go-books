package errors

import (
	"errors"
)

// Common errors used across the library module
var (
	ErrNotFound  = errors.New("not found")
	ErrDatabase  = errors.New("database error")
	ErrInvalidID = errors.New("invalid ID")
)