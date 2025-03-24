package models

import (
	"errors"
	"strings"

	"github.com/google/uuid"
)

// Book represents a book entity
type Book struct {
	ID     string
	Title  string
	Author string
	ISBN   string
}

// NewBook creates a new book with validation
func NewBook(title, author, isbn string) (*Book, error) {
	// Validate book data
	if err := validateBook(title, author, isbn); err != nil {
		return nil, err
	}

	return &Book{
		ID:     uuid.New().String(),
		Title:  title,
		Author: author,
		ISBN:   isbn,
	}, nil
}

// validateBook checks if the book data is valid
func validateBook(title string, author string, isbn string) error {
	if strings.TrimSpace(title) == "" {
		return errors.New("title cannot be empty")
	}

	if strings.TrimSpace(author) == "" {
		return errors.New("author cannot be empty")
	}

	if strings.TrimSpace(isbn) == "" {
		return errors.New("ISBN cannot be empty")
	}

	return nil
}
