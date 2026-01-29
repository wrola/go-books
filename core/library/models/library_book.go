package models

import (
	"time"

	"books/core/storage/models"
)

type LibraryBook struct {
	ISBN        string    `json:"isbn"`
	Title       string    `json:"title"`
	Author      string    `json:"author"`
	PublishedAt time.Time `json:"published_at"`

	IsAvailable     bool       `json:"is_available"`
	CurrentBorrower string     `json:"current_borrower,omitempty"`
	DueDate         *time.Time `json:"due_date,omitempty"`
	IsOverdue       bool       `json:"is_overdue,omitempty"`
}

func NewLibraryBookFromStorageBook(book *models.Book, rentals []*BookRental) *LibraryBook {
	libraryBook := &LibraryBook{
		ISBN:        book.ISBN,
		Title:       book.Title,
		Author:      book.Author,
		PublishedAt: book.PublishedAt,
		IsAvailable: true,
	}

	for _, rental := range rentals {
		if rental.BookID == book.ISBN && !rental.IsReturned() {
			libraryBook.IsAvailable = false
			libraryBook.CurrentBorrower = rental.UserID
			libraryBook.DueDate = &rental.ReturnDeadline
			libraryBook.IsOverdue = rental.IsOverdue()
			break
		}
	}

	return libraryBook
}

func (lb *LibraryBook) DaysUntilDue() int {
	if lb.IsAvailable || lb.DueDate == nil {
		return 0
	}

	duration := lb.DueDate.Sub(time.Now())
	return int(duration.Hours() / 24)
}
