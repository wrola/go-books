package core

import (
	"books/core/commands"
	"books/core/queries"
)

type Core struct {
	Commands commands.Commands
	Queries  queries.Queries
}



