package commands

import (
	"context"
	"errors"
	"strings"

	"books/core/storage/models"
	"books/core/storage/repositories/interfaces"
)

// UpdateBookCommand represents the command to update an existing book
type UpdateBookCommand struct {
	ISBN  string // ISBN of the book to update
	Title string
	Author string
}

// UpdateBookCommandHandler handles UpdateBookCommand
type UpdateBookCommandHandler struct {
	repo interfaces.BookStoragePostgresRepository
}

// NewUpdateBookCommandHandler creates a new UpdateBookCommandHandler
func NewUpdateBookCommandHandler(repo interfaces.BookStoragePostgresRepository) *UpdateBookCommandHandler {
	return &UpdateBookCommandHandler{
		repo: repo,
	}
}

// Handle processes the UpdateBookCommand
func (h *UpdateBookCommandHandler) Handle(ctx context.Context, cmd interface{}) error {
	command, ok := cmd.(UpdateBookCommand)
	if !ok {
		return ErrInvalidCommandType
	}

	// Validate OldISBN
	if strings.TrimSpace(command.ISBN) == "" {
		return errors.New("book ISBN cannot be empty")
	}

	// Check if at least one field is provided for update
	if command.Title == "" && command.Author == "" {
		return errors.New("at least one field must be provided for update")
	}

	// Get the book to update
	bookToUpdate, err := h.repo.FindByISBN(ctx, command.ISBN)
	if err != nil {
		return err // This will be ErrBookNotFound if the book doesn't exist
	}

	// Create a new book with updated fields
	newBook := &models.Book{
		ISBN:        bookToUpdate.ISBN,
		Title:       bookToUpdate.Title,
		Author:      bookToUpdate.Author,
		PublishedAt: bookToUpdate.PublishedAt,
	}

	// Update only the provided fields
	if command.Title != "" {
		newBook.Title = command.Title
	}
	if command.Author != "" {
		newBook.Author = command.Author
	}

	// Delete the old book if ISBN is being updated
	if err := h.repo.Delete(ctx, command.ISBN); err != nil {
		return err
	}

	// Save the updated book
	return h.repo.Save(ctx, newBook)
}