package commands

import (
	"context"
	"errors"
	"strings"

	"books/core/storage/repositories"
)

// DeleteBookCommand represents the command to delete a book
type DeleteBookCommand struct {
	ISBN string
}

// DeleteBookCommandHandler handles DeleteBookCommand
type DeleteBookCommandHandler struct {
	repo repositories.BookRepository
}

// NewDeleteBookCommandHandler creates a new DeleteBookCommandHandler
func NewDeleteBookCommandHandler(repo repositories.BookRepository) *DeleteBookCommandHandler {
	return &DeleteBookCommandHandler{repo: repo}
}

// Handle processes the DeleteBookCommand
func (h *DeleteBookCommandHandler) Handle(ctx context.Context, cmd interface{}) error {
	command, ok := cmd.(DeleteBookCommand)
	if !ok {
		return ErrInvalidCommandType
	}

	// Validate book ID
	if strings.TrimSpace(command.ISBN) == "" {
		return errors.New("book ID cannot be empty")
	}

	// Check if the book exists first
	_, err := h.repo.FindByISBN(ctx, command.ISBN)
	if err != nil {
		return err // This will be ErrBookNotFound if the book doesn't exist
	}

	// Delete the book
	return h.repo.Delete(ctx, command.ISBN)
}