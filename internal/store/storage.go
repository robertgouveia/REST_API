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
		Create(context.Context, *Post) error
	}

	Users interface {
		Create(context.Context, *User) error
	}
}

// Defining a Store and supplying the dependencies
func NewStorage(db *sql.DB) Storage {
	//Creating and returning a Storage object with Repository References
	return Storage{
		Posts: &PostStore{db},
		Users: &UserStore{db},
	}
}
