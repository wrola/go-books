package application

import (
	"books/core"
	"books/core/commands"
	"books/core/models"
	"books/core/queries"
	"context"
)

func NewApplication(ctx context.Context) *core.Core {
	application := &core.Core{
		Commands: commands.Commands{
			AddBook: func(ctx context.Context, command commands.AddBookCommand) error {
				// Dummy implementation
				return nil
			},
		},
		Queries: queries.Queries{
			GetBook: func(ctx context.Context, query queries.GetBookQuery) (*models.Book, error) {
				// Dummy implementation
				return &models.Book{
					ID:     "dummy-id",
					Title:  "Dummy Book",
					Author: "Dummy Author",
					ISBN:   "1234567890",
				}, nil
			},
			GetBookByISBN: func(ctx context.Context, query queries.GetBookByISBNQuery) (*models.Book, error) {
				// Dummy implementation
				return &models.Book{
					ID:     "dummy-id",
					Title:  "Dummy Book",
					Author: "Dummy Author",
					ISBN:   query.ISBN,
				}, nil
			},
		},
	}

	return application
}
