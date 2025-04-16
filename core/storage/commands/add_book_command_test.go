package commands

import (
	"context"
	"errors"
	"testing"

	"books/core/storage/models"
)

// MockBookRepository is a mock implementation of BookRepository for testing
type MockBookRepository struct {
	books          []*models.Book
	saveCalled     bool
	saveError      error
	findByIDError  error
	deleteError    error
}

func NewMockBookRepository() *MockBookRepository {
	return &MockBookRepository{
		books: make([]*models.Book, 0),
	}
}

func (r *MockBookRepository) Save(ctx context.Context, book *models.Book) error {
	r.saveCalled = true
	if r.saveError != nil {
		return r.saveError
	}

	// Try to find and update existing book
	for i, existingBook := range r.books {
		if existingBook.ISBN == book.ISBN {
			// Create a new book instance to avoid reference issues
			r.books[i] = &models.Book{
				ISBN:        book.ISBN,
				Title:       book.Title,
				Author:      book.Author,
				PublishedAt: book.PublishedAt,
			}
			return nil
		}
	}

	// If not found, append as new book
	r.books = append(r.books, book)
	return nil
}

func (r *MockBookRepository) FindAll(ctx context.Context) ([]*models.Book, error) {
	return r.books, nil
}

func (r *MockBookRepository) FindByISBN(ctx context.Context, isbn string) (*models.Book, error) {
	if r.findByIDError != nil {
		return nil, r.findByIDError
	}

	for _, book := range r.books {
		if book.ISBN == isbn {
			// Return a copy to prevent external modifications
			return &models.Book{
				ISBN:        book.ISBN,
				Title:       book.Title,
				Author:      book.Author,
				PublishedAt: book.PublishedAt,
			}, nil
		}
	}
	return nil, ErrBookNotFound
}

func (r *MockBookRepository) Delete(ctx context.Context, isbn string) error {
	if r.deleteError != nil {
		return r.deleteError
	}

	for i, book := range r.books {
		if book.ISBN == isbn {
			r.books = append(r.books[:i], r.books[i+1:]...)
			return nil
		}
	}
	return ErrBookNotFound
}

func TestAddBookCommandHandler_Handle(t *testing.T) {
	// Test cases
	tests := []struct {
		name        string
		command     interface{}
		mockSaveErr error
		wantErr     bool
		errType     error
	}{
		{
			name: "Valid AddBookCommand",
			command: AddBookCommand{
				Title:  "The Great Gatsby",
				Author: "F. Scott Fitzgerald",
				ISBN:   "9780743273565",
			},
			wantErr: false,
		},
		{
			name: "Empty Title",
			command: AddBookCommand{
				Title:  "",
				Author: "F. Scott Fitzgerald",
				ISBN:   "9780743273565",
			},
			wantErr: true,
		},
		{
			name: "Empty Author",
			command: AddBookCommand{
				Title:  "The Great Gatsby",
				Author: "",
				ISBN:   "9780743273565",
			},
			wantErr: true,
		},
		{
			name: "Empty ISBN",
			command: AddBookCommand{
				Title:  "The Great Gatsby",
				Author: "F. Scott Fitzgerald",
				ISBN:   "",
			},
			wantErr: true,
		},
		{
			name:    "Invalid command type",
			command: "not a command",
			wantErr: true,
			errType: ErrInvalidCommandType,
		},
		{
			name: "Repository save error",
			command: AddBookCommand{
				Title:  "The Great Gatsby",
				Author: "F. Scott Fitzgerald",
				ISBN:   "9780743273565",
			},
			mockSaveErr: errors.New("database error"),
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockRepo := NewMockBookRepository()
			mockRepo.saveError = tt.mockSaveErr

			handler := NewAddBookCommandHandler(mockRepo)

			// Execute
			err := handler.Handle(context.Background(), tt.command)

			// Verify
			if (err != nil) != tt.wantErr {
				t.Errorf("Handle() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.errType != nil && err != tt.errType {
				t.Errorf("Handle() error = %v, expected error type %v", err, tt.errType)
				return
			}

			// For valid cases, verify the book was saved
			if !tt.wantErr && !mockRepo.saveCalled {
				t.Errorf("Save() was not called")
			}

			// For valid cases, verify we have 1 book in the repo
			if !tt.wantErr && len(mockRepo.books) != 1 {
				t.Errorf("Expected 1 book in repository, got %d", len(mockRepo.books))
			}

			// For valid cases, verify book properties
			if !tt.wantErr && len(mockRepo.books) > 0 {
				book := mockRepo.books[0]
				cmd, _ := tt.command.(AddBookCommand)

				if book.Title != cmd.Title {
					t.Errorf("Book title = %v, want %v", book.Title, cmd.Title)
				}

				if book.Author != cmd.Author {
					t.Errorf("Book author = %v, want %v", book.Author, cmd.Author)
				}

				if book.ISBN != cmd.ISBN {
					t.Errorf("Book ISBN = %v, want %v", book.ISBN, cmd.ISBN)
				}

			}
		})
	}
}