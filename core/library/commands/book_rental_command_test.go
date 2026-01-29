package commands

import (
	"books/core/library/errors"
	"books/core/library/models"
	"books/core/library/repositories"
	storage_models "books/core/storage/models"
	"context"
	stderrors "errors"
	"testing"
)

// Mock repository for testing
type mockBookRepository struct {
	books       map[string]*storage_models.Book
	rentals     map[string]*models.BookRental
	userRentals map[string][]*models.BookRental
	saveErr     error
}

func newMockRepository() *mockBookRepository {
	return &mockBookRepository{
		books:       make(map[string]*storage_models.Book),
		rentals:     make(map[string]*models.BookRental),
		userRentals: make(map[string][]*models.BookRental),
	}
}

// Implement repositories.BookRepository interface
func (m *mockBookRepository) GetBookByISBN(ctx context.Context, isbn string) (*storage_models.Book, error) {
	book, exists := m.books[isbn]
	if !exists {
		return nil, stderrors.New("book not found")
	}
	return book, nil
}

func (m *mockBookRepository) BookExists(ctx context.Context, isbn string) (bool, error) {
	_, exists := m.books[isbn]
	return exists, nil
}

func (m *mockBookRepository) GetActiveBookRentalByBookID(ctx context.Context, bookID string) (*models.BookRental, error) {
	rental, exists := m.rentals[bookID]
	if !exists {
		return nil, errors.ErrNotFound
	}
	return rental, nil
}

func (m *mockBookRepository) GetAllUserRentals(ctx context.Context, userID string) ([]*models.BookRental, error) {
	rentals, exists := m.userRentals[userID]
	if !exists {
		return []*models.BookRental{}, nil
	}
	return rentals, nil
}

func (m *mockBookRepository) SaveBookRental(ctx context.Context, rental *models.BookRental) error {
	if m.saveErr != nil {
		return m.saveErr
	}
	m.rentals[rental.BookID] = rental

	// Also add to user rentals
	if _, exists := m.userRentals[rental.UserID]; !exists {
		m.userRentals[rental.UserID] = []*models.BookRental{}
	}
	userRentals := m.userRentals[rental.UserID]
	userRentals = append(userRentals, rental)
	m.userRentals[rental.UserID] = userRentals

	return nil
}

// Verify the mock implements the interface
var _ repositories.BookRepository = (*mockBookRepository)(nil)

type bookRentalTestCase struct {
	name        string
	setupRepo   func(repo *mockBookRepository)
	command     BookRentalCommand
	expectError bool
	errorMsg    string
}

func getBookRentalTestCases() []bookRentalTestCase {
	return []bookRentalTestCase{
		{
			name: "Successfully borrow a book",
			setupRepo: func(repo *mockBookRepository) {
				repo.books["book1"] = &storage_models.Book{ISBN: "book1", Title: "Test Book"}
			},
			command: BookRentalCommand{
				BookID: "book1",
				UserID: "user1",
			},
			expectError: false,
		},
		{
			name:        "Invalid command type",
			setupRepo:   func(repo *mockBookRepository) {},
			command:     BookRentalCommand{}, // We'll pass a different type in the test
			expectError: true,
			errorMsg:    "invalid command type",
		},
		{
			name:      "Book does not exist",
			setupRepo: func(repo *mockBookRepository) {},
			command: BookRentalCommand{
				BookID: "nonexistent",
				UserID: "user1",
			},
			expectError: true,
			errorMsg:    "book not found in storage",
		},
		{
			name: "Book already borrowed by someone else",
			setupRepo: func(repo *mockBookRepository) {
				repo.books["book1"] = &storage_models.Book{ISBN: "book1", Title: "Test Book"}
				rental := models.NewBookRental("book1", "user2")
				repo.rentals["book1"] = rental
			},
			command: BookRentalCommand{
				BookID: "book1",
				UserID: "user1",
			},
			expectError: true,
			errorMsg:    "book is already borrowed by someone else",
		},
		{
			name: "User already has this book",
			setupRepo: func(repo *mockBookRepository) {
				repo.books["book1"] = &storage_models.Book{ISBN: "book1", Title: "Test Book"}
				rental := models.NewBookRental("book1", "user1")
				repo.userRentals["user1"] = []*models.BookRental{rental}
			},
			command: BookRentalCommand{
				BookID: "book1",
				UserID: "user1",
			},
			expectError: true,
			errorMsg:    "you already have this book, please return it before renting again",
		},
		{
			name: "Book is returned and can be rented again",
			setupRepo: func(repo *mockBookRepository) {
				repo.books["book1"] = &storage_models.Book{ISBN: "book1", Title: "Test Book"}
				rental := models.NewBookRental("book1", "user1")
				rental.MarkAsReturned()
				repo.userRentals["user1"] = []*models.BookRental{rental}
			},
			command: BookRentalCommand{
				BookID: "book1",
				UserID: "user1",
			},
			expectError: false,
		},
		{
			name: "Save error",
			setupRepo: func(repo *mockBookRepository) {
				repo.books["book1"] = &storage_models.Book{ISBN: "book1", Title: "Test Book"}
				repo.saveErr = errors.ErrDatabase
			},
			command: BookRentalCommand{
				BookID: "book1",
				UserID: "user1",
			},
			expectError: true,
			errorMsg:    "database error",
		},
	}
}

func TestBookRentalCommandHandler_Handle(t *testing.T) {
	tests := getBookRentalTestCases()

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Setup
			repo := newMockRepository()
			tc.setupRepo(repo)
			handler := NewBookRentalCommandHandler(repo)

			// Execute
			var err error
			if tc.name == "Invalid command type" {
				// Test with wrong command type
				err = handler.Handle(context.Background(), "not a command")
			} else {
				err = handler.Handle(context.Background(), tc.command)
			}

			// Assert
			if tc.expectError {
				if err == nil {
					t.Errorf("expected error but got nil")
					return
				}
				if err.Error() != tc.errorMsg {
					t.Errorf("expected error message '%s', got '%s'", tc.errorMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("expected no error but got: %v", err)
					return
				}

				// Only check rental details if we expect the operation to succeed
				rental, exists := repo.rentals[tc.command.BookID]
				if !exists {
					t.Errorf("rental was not saved in the repository")
					return
				}
				if rental.BookID != tc.command.BookID {
					t.Errorf("expected BookID '%s', got '%s'", tc.command.BookID, rental.BookID)
				}
				if rental.UserID != tc.command.UserID {
					t.Errorf("expected UserID '%s', got '%s'", tc.command.UserID, rental.UserID)
				}
				if rental.IsReturned() {
					t.Errorf("expected rental to not be returned, but it was")
				}
			}
		})
	}
}