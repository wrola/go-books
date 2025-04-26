package commands

import (
	"context"
	"testing"
	"time"

	"books/core/storage/models"
	"books/core/storage/repositories"
	"books/core/storage/repositories/interfaces"
	"errors"
)

func TestDeleteBookCommandHandler_Handle(t *testing.T) {
	// Create a test book
	testBook, _ := models.NewBook("1234567890", "Test Book", "Test Author", time.Now())

	// Test cases
	tests := []struct {
		name           string
		setupRepo      func(*repositories.BookStorageInMemoryRepository)
		command        interface{}
		validateResult func(*testing.T, *repositories.BookStorageInMemoryRepository)
		wantErr        bool
		expectedErr    error
	}{
		{
			name: "successful deletion",
			setupRepo: func(repo *repositories.BookStorageInMemoryRepository) {
				repo.Save(context.Background(), testBook)
			},
			command: &DeleteBookCommand{
				ISBN: testBook.ISBN,
			},
			wantErr: false,
			validateResult: func(t *testing.T, repo *repositories.BookStorageInMemoryRepository) {
				_, err := repo.FindByISBN(context.Background(), testBook.ISBN)
				if err != interfaces.ErrBookNotFound {
					t.Errorf("expected book to be deleted, got error: %v", err)
				}
			},
		},
		{
			name: "book not found",
			setupRepo: func(repo *repositories.BookStorageInMemoryRepository) {
				// Empty repository
			},
			command: &DeleteBookCommand{
				ISBN: "non-existent-isbn",
			},
			wantErr:     true,
			expectedErr: interfaces.ErrBookNotFound,
		},
		{
			name: "empty ISBN",
			setupRepo: func(repo *repositories.BookStorageInMemoryRepository) {},
			command: &DeleteBookCommand{
				ISBN: "",
			},
			wantErr:   true,
			expectedErr: errors.New("book ID cannot be empty"),
		},
		{
			name:      "invalid command type",
			setupRepo: func(repo *repositories.BookStorageInMemoryRepository) {},
			command:   &AddBookCommand{ISBN: "1234567890"}, // Pass a different command type
			wantErr:   true,
			expectedErr: ErrInvalidCommandType,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := repositories.NewBookStorageInMemoryRepository()
			if tt.setupRepo != nil {
				tt.setupRepo(mockRepo)
			}

			handler := NewDeleteBookCommandHandler(mockRepo)
			err := handler.Handle(context.Background(), tt.command)

			if tt.wantErr {
				if err == nil {
					t.Error("expected error but got none")
				}
				if err.Error() != tt.expectedErr.Error() {
					t.Errorf("expected error %v but got %v", tt.expectedErr, err)
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if tt.validateResult != nil {
				tt.validateResult(t, mockRepo)
			}
		})
	}
}