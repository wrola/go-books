package handlers

import (
	"context"
	"time"

	"books/core/storage/commands"
	"books/core/storage/models"
	"books/core/storage/repositories/interfaces"
)

type CreateBookCommand struct {
	ISBN  string
	Title string
	Author string
}

type CreateBookHandler struct {
	bookRepository interfaces.BookRepository
}

func NewCreateBookHandler(bookRepository interfaces.BookRepository) *CreateBookHandler {
	return &CreateBookHandler{
		bookRepository: bookRepository,
	}
}

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