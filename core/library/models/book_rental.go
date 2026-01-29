package models

import (
	"time"
)

type BookRental struct {
	BookID         string     `json:"book_id"`
	UserID         string     `json:"user_id"`
	BorrowedAt     time.Time  `json:"borrowed_at"`
	ReturnDeadline time.Time  `json:"return_deadline"`
	ReturnedAt     *time.Time `json:"returned_at,omitempty"`
}

func NewBookRental(bookID, userID string) *BookRental {
	now := time.Now()
	deadline := now.Add(14 * 24 * time.Hour)

	return &BookRental{
		BookID:         bookID,
		UserID:         userID,
		BorrowedAt:     now,
		ReturnDeadline: deadline,
	}
}

func (b *BookRental) IsReturned() bool {
	return b.ReturnedAt != nil
}

func (b *BookRental) MarkAsReturned() {
	now := time.Now()
	b.ReturnedAt = &now
}

func (b *BookRental) IsOverdue() bool {
	if b.IsReturned() {
		return false
	}
	return time.Now().After(b.ReturnDeadline)
}

func (b *BookRental) DaysUntilDue() int {
	if b.IsReturned() {
		return 0
	}

	duration := b.ReturnDeadline.Sub(time.Now())
	return int(duration.Hours() / 24)
}

func BookIsAvailable(bookID string, rentals []*BookRental) bool {
	for _, rental := range rentals {
		if rental.BookID == bookID && !rental.IsReturned() {
			return false
		}
	}
	return true
}
