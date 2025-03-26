package queries

import (
	"books/core/models"
	"context"
)


type GetBookQuery struct {
	ID string
}

type GetBookByISBNQuery struct {
	ISBN string
}

type GetBookHandler func(ctx context.Context, query GetBookQuery) (*models.Book, error)
type GetBookByISBNHandler func(ctx context.Context, query GetBookByISBNQuery) (*models.Book, error)

type Queries struct {
	GetBook       GetBookHandler
	GetBookByISBN GetBookByISBNHandler
}
