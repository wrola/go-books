package commands

import (
	"context"
	"errors"
	"strings"
)

// DeleteBookCommand represents the command to delete a book
type DeleteBookCommand struct {
	ID string
}

// DeleteBookCommandHandler handles DeleteBookCommand
type DeleteBookCommandHandler struct {
	repo BookRepository
}

// NewDeleteBookCommandHandler creates a new DeleteBookCommandHandler
func NewDeleteBookCommandHandler(repo BookRepository) *DeleteBookCommandHandler {
	return &DeleteBookCommandHandler{repo: repo}
}

// Handle processes the DeleteBookCommand
func (h *DeleteBookCommandHandler) Handle(ctx context.Context, cmd interface{}) error {
	command, ok := cmd.(DeleteBookCommand)
	if !ok {
		return ErrInvalidCommandType
	}

	// Validate book ID
	if strings.TrimSpace(command.ID) == "" {
		return errors.New("book ID cannot be empty")
	}

	// Check if the book exists first
	_, err := h.repo.FindByID(ctx, command.ID)
	if err != nil {
		return err // This will be ErrBookNotFound if the book doesn't exist
	}

	// Delete the book
	return h.repo.Delete(ctx, command.ID)
}