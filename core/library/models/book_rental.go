package models

import (
	"time"
)

// BookRental represents a book that has been borrowed from the library
type BookRental struct {
	BookID          string     `json:"book_id"`
	UserID          string     `json:"user_id"`
	BorrowedAt      time.Time  `json:"borrowed_at"`
	ReturnDeadline  time.Time  `json:"return_deadline"`
	ReturnedAt      *time.Time `json:"returned_at,omitempty"`
}

// NewBookRental creates a new BookRental record for a book
func NewBookRental(bookID, userID string) *BookRental {
	now := time.Now()
	deadline := now.Add(14 * 24 * time.Hour) // 2 weeks borrowing period

	return &BookRental{
		BookID:         bookID,
		UserID:         userID,
		BorrowedAt:     now,
		ReturnDeadline: deadline,
	}
}

// IsReturned checks if the book has been returned
func (b *BookRental) IsReturned() bool {
	return b.ReturnedAt != nil
}

// MarkAsReturned marks a book as returned
func (b *BookRental) MarkAsReturned() {
	now := time.Now()
	b.ReturnedAt = &now
}

// IsOverdue checks if the rental is overdue
func (b *BookRental) IsOverdue() bool {
	if b.IsReturned() {
		return false
	}
	return time.Now().After(b.ReturnDeadline)
}

// DaysUntilDue returns the number of days until the book is due
// returns negative number if overdue
func (b *BookRental) DaysUntilDue() int {
	if b.IsReturned() {
		return 0
	}

	duration := b.ReturnDeadline.Sub(time.Now())
	return int(duration.Hours() / 24)
}

// BookIsAvailable checks if a book is available for borrowing
// This is a utility function to check if any rental records exist for a book
func BookIsAvailable(bookID string, rentals []*BookRental) bool {
	for _, rental := range rentals {
		if rental.BookID == bookID && !rental.IsReturned() {
			return false
		}
	}
	return true
}