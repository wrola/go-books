package commands

import (
	"context"
	"errors"
	"testing"
	"time"
	"books/core/storage/models"
)

func TestDeleteBookCommandHandler_Handle(t *testing.T) {
	// Create a test book
	testBook, _ := models.NewBook("Test Book", "Test Author", "1234567890", time.Now())

	// Test cases
	tests := []struct {
		name        string
		setupRepo   func(*MockBookRepository)
		command     interface{}
		wantErr     bool
		expectedErr error
	}{
		{
			name: "Delete existing book",
			setupRepo: func(repo *MockBookRepository) {
				repo.books = append(repo.books, testBook)
			},
			command: DeleteBookCommand{
				ISBN: testBook.ISBN,
			},
			wantErr: false,
		},
		{
			name: "Delete non-existent book",
			setupRepo: func(repo *MockBookRepository) {
				// Empty repository
			},
			command: DeleteBookCommand{
				ISBN: "non-existent-isbn",
			},
			wantErr:     true,
			expectedErr: ErrBookNotFound,
		},
		{
			name:      "Invalid command type",
			setupRepo: func(repo *MockBookRepository) {},
			command:   "not a command",
			wantErr:   true,
			expectedErr: ErrInvalidCommandType,
		},
		{
			name: "Empty ID",
			setupRepo: func(repo *MockBookRepository) {},
			command: DeleteBookCommand{
				ISBN: "",
			},
			wantErr:     true,
			expectedErr: errors.New("book ID cannot be empty"),
		},
		{
			name: "Repository error during lookup",
			setupRepo: func(repo *MockBookRepository) {
				repo.findByIDError = errors.New("database error")
			},
			command: DeleteBookCommand{
				ISBN: testBook.ISBN,
			},
			wantErr:     true,
			expectedErr: errors.New("database error"),
		},
		{
			name: "Repository error during delete",
			setupRepo: func(repo *MockBookRepository) {
				repo.books = append(repo.books, testBook)
				repo.deleteError = errors.New("delete error")
			},
			command: DeleteBookCommand{
				ISBN: testBook.ISBN,
			},
			wantErr:     true,
			expectedErr: errors.New("delete error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockRepo := NewMockBookRepository()
			tt.setupRepo(mockRepo)

			handler := NewDeleteBookCommandHandler(mockRepo)

			// Execute
			err := handler.Handle(context.Background(), tt.command)

			// Verify
			if (err != nil) != tt.wantErr {
				t.Errorf("Handle() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Check specific error type if expected
			if tt.expectedErr != nil && err != nil {
				if tt.expectedErr.Error() != err.Error() {
					t.Errorf("Expected error '%v', got '%v'", tt.expectedErr, err)
				}
			}

			// For successful delete, verify the book was removed
			if !tt.wantErr {
				cmd, _ := tt.command.(DeleteBookCommand)
				for _, book := range mockRepo.books {
					if book.ISBN == cmd.ISBN {
						t.Errorf("Book with ISBN %s still exists in repository after delete", cmd.ISBN)
					}
				}
			}
		})
	}
}