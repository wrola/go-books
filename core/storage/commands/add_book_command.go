package commands

import (
	"context"
	"errors"
	"strings"
	"time"

	"books/core/storage/models"
	"books/core/storage/repositories/interfaces"
)

// AddBookCommand represents the command to add a new book
type AddBookCommand struct {
	ISBN  string
	Title string
	Author string
}

// AddBookCommandHandler handles AddBookCommand
type AddBookCommandHandler struct {
	repo interfaces.BookStoragePostgresRepository
}

// NewAddBookCommandHandler creates a new AddBookCommandHandler
func NewAddBookCommandHandler(repo interfaces.BookStoragePostgresRepository) *AddBookCommandHandler {
	return &AddBookCommandHandler{
		repo: repo,
	}
}

// Handle processes the AddBookCommand
func (h *AddBookCommandHandler) Handle(ctx context.Context, cmd interface{}) error {
	command, ok := cmd.(AddBookCommand)
	if !ok {
		return ErrInvalidCommandType
	}

	// Validate inputs
	if strings.TrimSpace(command.Title) == "" {
		return errors.New("title cannot be empty")
	}

	if strings.TrimSpace(command.Author) == "" {
		return errors.New("author cannot be empty")
	}

	if strings.TrimSpace(command.ISBN) == "" {
		return errors.New("ISBN cannot be empty")
	}

	book, err := models.NewBook(command.ISBN, command.Title, command.Author, time.Now())
	if err != nil {
		return err
	}

	return h.repo.Save(ctx, book)
}