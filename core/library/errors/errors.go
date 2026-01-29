package errors

import (
	"errors"
)

var (
	ErrNotFound  = errors.New("not found")
	ErrDatabase  = errors.New("database error")
	ErrInvalidID = errors.New("invalid ID")
)
