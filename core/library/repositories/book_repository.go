package repositories

import (
	"context"

	"books/core/library/models"
	storage_models "books/core/storage/models"
)

type BookRepository interface {
	GetBookByISBN(ctx context.Context, isbn string) (*storage_models.Book, error)
	BookExists(ctx context.Context, isbn string) (bool, error)

	GetActiveBookRentalByBookID(ctx context.Context, bookID string) (*models.BookRental, error)
	GetAllUserRentals(ctx context.Context, userID string) ([]*models.BookRental, error)
	SaveBookRental(ctx context.Context, rental *models.BookRental) error
}
