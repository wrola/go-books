package models

import (
	"books/core/storage/models"
	"time"
)

// LibraryBook represents a read-only view of a book in the library
// with additional information about its availability
type LibraryBook struct {
	ISBN        string    `json:"isbn"`
	Title       string    `json:"title"`
	Author      string    `json:"author"`
	PublishedAt time.Time `json:"published_at"`

	// Availability information
	IsAvailable     bool       `json:"is_available"`
	CurrentBorrower string     `json:"current_borrower,omitempty"`
	DueDate         *time.Time `json:"due_date,omitempty"`
	IsOverdue       bool       `json:"is_overdue,omitempty"`
}

// NewLibraryBookFromStorageBook creates a new LibraryBook from a storage Book model
func NewLibraryBookFromStorageBook(book *models.Book, rentals []*BookRental) *LibraryBook {
	libraryBook := &LibraryBook{
		ISBN:        book.ISBN,
		Title:       book.Title,
		Author:      book.Author,
		PublishedAt: book.PublishedAt,
		IsAvailable: true,
	}

	// Find active rental for this book
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

// DaysUntilDue returns the number of days until the book is due
// returns negative number if overdue, 0 if available
func (lb *LibraryBook) DaysUntilDue() int {
	if lb.IsAvailable || lb.DueDate == nil {
		return 0
	}

	duration := lb.DueDate.Sub(time.Now())
	return int(duration.Hours() / 24)
}