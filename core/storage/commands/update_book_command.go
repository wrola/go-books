package commands

import (
	"context"
	"errors"
	"strings"

	"books/core/storage/models"
	"books/core/storage/repositories/interfaces"
)

type UpdateBookCommand struct {
	ISBN  string 
	Title string
	Author string
}

type UpdateBookCommandHandler struct {
	repo interfaces.BookRepository
}

func NewUpdateBookCommandHandler(repo interfaces.BookRepository) *UpdateBookCommandHandler {
	return &UpdateBookCommandHandler{
		repo: repo,
	}
}

func (h *UpdateBookCommandHandler) Handle(ctx context.Context, cmd interface{}) error {
	if cmd == nil {
		return ErrInvalidCommandType
	}

	command, ok := cmd.(*UpdateBookCommand)
	if !ok {
		return ErrInvalidCommandType
	}

	if strings.TrimSpace(command.ISBN) == "" {
		return errors.New("book ISBN cannot be empty")
	}

	if command.Title == "" && command.Author == "" {
		return errors.New("at least one field must be provided for update")
	}

	bookToUpdate, err := h.repo.FindByISBN(ctx, command.ISBN)
	if err != nil {
		return err 
	}

	newBook := &models.Book{
		ISBN:        bookToUpdate.ISBN,
		Title:       bookToUpdate.Title,
		Author:      bookToUpdate.Author,
		PublishedAt: bookToUpdate.PublishedAt,
	}

	if command.Title != "" {
		newBook.Title = command.Title
	}
	if command.Author != "" {
		newBook.Author = command.Author
	}

	if err := h.repo.Delete(ctx, command.ISBN); err != nil {
		return err
	}

	return h.repo.Save(ctx, newBook)
}