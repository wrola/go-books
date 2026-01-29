package repositories

import (
	"context"
	"database/sql"
	"fmt"

	"books/core/storage/models"
	"books/core/storage/repositories/interfaces"

	_ "github.com/lib/pq"
)

// BookStoragePostgresRepository implements BookStoragePostgresRepository interface using PostgreSQL
type BookStoragePostgresRepository struct {
	db *sql.DB
}

// NewBookStoragePostgresRepository creates a new PostgreSQL book repository
func NewBookStoragePostgresRepository(db *sql.DB) *BookStoragePostgresRepository {
	return &BookStoragePostgresRepository{
		db: db,
	}
}

// Save adds or updates a book in the repository
func (r *BookStoragePostgresRepository) Save(ctx context.Context, book *models.Book) error {
	query := `
		INSERT INTO books (isbn, title, author, published_at)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (isbn) DO UPDATE
		SET title = $2, author = $3, published_at = $4
	`

	_, err := r.db.ExecContext(ctx, query,
		book.ISBN,
		book.Title,
		book.Author,
		book.PublishedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to save book: %w", err)
	}
	return nil
}

// FindAll returns all books in the repository
func (r *BookStoragePostgresRepository) FindAll(ctx context.Context) ([]*models.Book, error) {
	query := `SELECT isbn, title, author, published_at FROM books`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query books: %w", err)
	}
	defer rows.Close()

	var books []*models.Book
	for rows.Next() {
		book := &models.Book{}
		err := rows.Scan(&book.ISBN, &book.Title, &book.Author, &book.PublishedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan book: %w", err)
		}
		books = append(books, book)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate books: %w", err)
	}

	return books, nil
}

// FindByISBN returns a book by its ISBN
func (r *BookStoragePostgresRepository) FindByISBN(ctx context.Context, isbn string) (*models.Book, error) {
	query := `SELECT isbn, title, author, published_at FROM books WHERE isbn = $1`

	book := &models.Book{}
	err := r.db.QueryRowContext(ctx, query, isbn).Scan(
		&book.ISBN,
		&book.Title,
		&book.Author,
		&book.PublishedAt,
	)

	if err == sql.ErrNoRows {
		return nil, interfaces.ErrBookNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find book by ISBN: %w", err)
	}

	return book, nil
}

// Delete removes a book from the repository
func (r *BookStoragePostgresRepository) Delete(ctx context.Context, isbn string) error {
	query := `DELETE FROM books WHERE isbn = $1`

	result, err := r.db.ExecContext(ctx, query, isbn)
	if err != nil {
		return fmt.Errorf("failed to delete book: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return interfaces.ErrBookNotFound
	}

	return nil
}

// Ensure BookStoragePostgresRepository implements BookRepository interface
var _ interfaces.BookRepository = (*BookStoragePostgresRepository)(nil)