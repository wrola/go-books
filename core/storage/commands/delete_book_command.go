package commands

import (
	"context"
	"errors"
	"strings"

	"books/core/storage/repositories/interfaces"
)

type DeleteBookCommand struct {
	ISBN string 
}

type DeleteBookCommandHandler struct {
	repo interfaces.BookRepository
}

func NewDeleteBookCommandHandler(repo interfaces.BookRepository) *DeleteBookCommandHandler {
	return &DeleteBookCommandHandler{
		repo: repo,
	}
}

func (h *DeleteBookCommandHandler) Handle(ctx context.Context, cmd interface{}) error {
	if cmd == nil {
		return ErrInvalidCommandType
	}

	command, ok := cmd.(*DeleteBookCommand)
	if !ok {
		return ErrInvalidCommandType
	}

	if strings.TrimSpace(command.ISBN) == "" {
		return errors.New("book ID cannot be empty")
	}

	_, err := h.repo.FindByISBN(ctx, command.ISBN)
	if err != nil {
		return err 
	}

	return h.repo.Delete(ctx, command.ISBN)
}