package commands

import (
	"books/core/libary/errors"
	"books/core/libary/models"
	"books/core/libary/repositories"
	"context"
	stderrors "errors"
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

	// Check if book exists before attempting to get it
	exists, err := h.repo.BookExists(ctx, command.BookID)
	if err != nil {
		return err
	}

	if !exists {
		return stderrors.New("book not found in storage")
	}

	// Check if book is already borrowed by someone else
	activeRental, err := h.repo.GetActiveBookRentalByBookID(ctx, command.BookID)
	if err != nil && !stderrors.Is(err, errors.ErrNotFound) {
		return err
	}

	if activeRental != nil && !activeRental.IsReturned() {
		return stderrors.New("book is already borrowed by someone else")
	}

	// Get all user rentals to check if they already have this book
	userRentals, err := h.repo.GetAllUserRentals(ctx, command.UserID)
	if err != nil {
		return err
	}

	// Check if user already has an active rental for this book
	for _, rental := range userRentals {
		// If user has this book and it's not returned yet
		if rental.BookID == command.BookID && !rental.IsReturned() {
			return stderrors.New("you already have this book, please return it before renting again")
		}
	}

	// Create new rental
	rental := models.NewBookRental(command.BookID, command.UserID)

	// Save the rental record
	err = h.repo.SaveBookRental(ctx, rental)
	if err != nil {
		return err
	}

	return nil
}
