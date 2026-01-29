package commands

import (
	"context"
	stderrors "errors"

	"books/core/library/errors"
	"books/core/library/models"
	"books/core/library/repositories"
)

type BookRentalCommand struct {
	ID     string
	BookID string
	UserID string
}

type BookRentalCommandHandler struct {
	repo repositories.BookRepository
}

func NewBookRentalCommandHandler(repo repositories.BookRepository) *BookRentalCommandHandler {
	return &BookRentalCommandHandler{repo: repo}
}

func (h *BookRentalCommandHandler) Handle(ctx context.Context, cmd interface{}) error {
	command, ok := cmd.(BookRentalCommand)
	if !ok {
		return stderrors.New("invalid command type")
	}

	exists, err := h.repo.BookExists(ctx, command.BookID)
	if err != nil {
		return err
	}

	if !exists {
		return stderrors.New("book not found in storage")
	}

	activeRental, err := h.repo.GetActiveBookRentalByBookID(ctx, command.BookID)
	if err != nil && !stderrors.Is(err, errors.ErrNotFound) {
		return err
	}

	if activeRental != nil && !activeRental.IsReturned() {
		return stderrors.New("book is already borrowed by someone else")
	}

	userRentals, err := h.repo.GetAllUserRentals(ctx, command.UserID)
	if err != nil {
		return err
	}

	for _, rental := range userRentals {
		if rental.BookID == command.BookID && !rental.IsReturned() {
			return stderrors.New("you already have this book, please return it before renting again")
		}
	}

	rental := models.NewBookRental(command.BookID, command.UserID)

	err = h.repo.SaveBookRental(ctx, rental)
	if err != nil {
		return err
	}

	return nil
}
