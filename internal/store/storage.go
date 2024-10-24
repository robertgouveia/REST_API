package store

import (
	"context"
	"database/sql"
)

// Repository Pattern for decoupling
type Storage struct {
	//Interfaces determine Tables / Store Points
	Posts interface {
		//Defining Methods
		Create(context.Context) error
	}

	Users interface {
		Create(context.Context) error
	}
}

// Defining a Store and supplying the dependencies
func NewStorage(db *sql.DB) Storage {
	//Creating and returning a Storage object with Repository References
	return Storage{
		Posts: &PostsStore{db},
		Users: &UsersStore{db},
	}
}
