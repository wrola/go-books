package handlers

import (
	"context"
	"time"

	"books/core/storage/commands"
	"books/core/storage/models"
	"books/core/storage/repositories/interfaces"
)

// CreateBookCommand represents a command to create a new book
type CreateBookCommand struct {
	ISBN  string
	Title string
	Author string
}

// CreateBookHandler handles the creation of new books
type CreateBookHandler struct {
	bookRepository interfaces.BookStoragePostgresRepository
}

// NewCreateBookHandler creates a new CreateBookHandler
func NewCreateBookHandler(bookRepository interfaces.BookStoragePostgresRepository) *CreateBookHandler {
	return &CreateBookHandler{
		bookRepository: bookRepository,
	}
}

// Handle processes the CreateBookCommand
func (h *CreateBookHandler) Handle(ctx context.Context, command interface{}) error {
	cmd, ok := command.(*CreateBookCommand)
	if !ok {
		return commands.ErrInvalidCommandType
	}

	book, err := models.NewBook(cmd.ISBN, cmd.Title, cmd.Author, time.Now())
	if err != nil {
		return err
	}

	return h.bookRepository.Save(ctx, book)
}