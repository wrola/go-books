package commands

import "context"

type AddBookCommand struct {
	Title  string
	Author string
	ISBN   string
}

type AddBookHandler func(ctx context.Context, command AddBookCommand) error

type Commands struct {
	AddBook AddBookHandler
}