package commands

import (
	"context"
	"errors"
	"strings"
)

// UpdateBookCommand represents the command to update an existing book
type UpdateBookCommand struct {
	ID     string
	Title  string
	Author string
	ISBN   string
}

// UpdateBookCommandHandler handles UpdateBookCommand
type UpdateBookCommandHandler struct {
	repo BookRepository
}

// NewUpdateBookCommandHandler creates a new UpdateBookCommandHandler
func NewUpdateBookCommandHandler(repo BookRepository) *UpdateBookCommandHandler {
	return &UpdateBookCommandHandler{repo: repo}
}

// Handle processes the UpdateBookCommand
func (h *UpdateBookCommandHandler) Handle(ctx context.Context, cmd interface{}) error {
	command, ok := cmd.(UpdateBookCommand)
	if !ok {
		return ErrInvalidCommandType
	}

	// Validate book ID
	if strings.TrimSpace(command.ID) == "" {
		return errors.New("book ID cannot be empty")
	}

	// Check if at least one field is provided for update
	if command.Title == "" && command.Author == "" && command.ISBN == "" {
		return errors.New("at least one field must be provided for update")
	}

	// Get the book to update
	bookToUpdate, err := h.repo.FindByID(ctx, command.ID)
	if err != nil {
		return err // This will be ErrBookNotFound if the book doesn't exist
	}

	// Update the book fields if provided
	if command.Title != "" {
		bookToUpdate.Title = command.Title
	}

	if command.Author != "" {
		bookToUpdate.Author = command.Author
	}

	if command.ISBN != "" {
		bookToUpdate.ISBN = command.ISBN
	}

	// Save the updated book
	return h.repo.Save(ctx, bookToUpdate)
}