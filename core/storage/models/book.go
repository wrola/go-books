package models

import (
	"errors"
	"strings"
	"time"
)

// Book represents a book in the library storage system
type Book struct {
	ISBN        string    `json:"isbn"`
	Title       string    `json:"title"`
	Author      string    `json:"author"`
	PublishedAt time.Time `json:"published_at"`
}

// NewBook creates a new Book instance
func NewBook(isbn, title, author string, publishedAt time.Time) (*Book, error) {
	if err := validateBook(title, author, isbn); err != nil {
		return nil, err
	}

	return &Book{
		ISBN:        isbn,
		Title:       title,
		Author:      author,
		PublishedAt: publishedAt,
	}, nil
}

// ValidateBook checks if the book data is valid
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
