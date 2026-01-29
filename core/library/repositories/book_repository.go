package repositories

import (
	"books/core/library/models"
	storage_models "books/core/storage/models"
	"context"
)

// BookRepository defines the repository interface for book operations
type BookRepository interface {
	// Book storage operations
	GetBookByISBN(ctx context.Context, isbn string) (*storage_models.Book, error)
	BookExists(ctx context.Context, isbn string) (bool, error)

	// Book rental operations
	GetActiveBookRentalByBookID(ctx context.Context, bookID string) (*models.BookRental, error)
	GetAllUserRentals(ctx context.Context, userID string) ([]*models.BookRental, error)
	SaveBookRental(ctx context.Context, rental *models.BookRental) error
}