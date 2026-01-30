package commands

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"books/core/storage/models"
	"books/core/storage/repositories/interfaces"
)

type AddBookCommand struct {
	ISBN  string
	Title string
	Author string
}

type AddBookCommandHandler struct {
	repo interfaces.BookRepository
}
func NewAddBookCommandHandler(repo interfaces.BookRepository) *AddBookCommandHandler {
	return &AddBookCommandHandler{
		repo: repo,
	}
}

func (h *AddBookCommandHandler) Handle(ctx context.Context, cmd interface{}) error {
	if cmd == nil {
		return ErrInvalidCommandType
	}

	command, ok := cmd.(*AddBookCommand)
	if !ok {
		return ErrInvalidCommandType
	}

	if strings.TrimSpace(command.Title) == "" {
		return errors.New("title cannot be empty")
	}

	if strings.TrimSpace(command.Author) == "" {
		return errors.New("author cannot be empty")
	}

	if strings.TrimSpace(command.ISBN) == "" {
		return errors.New("ISBN cannot be empty")
	}

	existingBook, err := h.repo.FindByISBN(ctx, command.ISBN)
	if err == nil && existingBook != nil {
		return fmt.Errorf("failed to save book: book with ISBN %s already exists", command.ISBN)
	}
	if err != nil && !errors.Is(err, interfaces.ErrBookNotFound) {
		return err
	}

	book, err := models.NewBook(command.ISBN, command.Title, command.Author, time.Now())
	if err != nil {
		return err
	}

	return h.repo.Save(ctx, book)
}