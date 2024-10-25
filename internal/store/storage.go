package store

import (
	"context"
	"database/sql"
	"errors"
)

var (
	ErrNotFound = errors.New("record not found")
)

// Repository Pattern for decoupling
type Storage struct {
	//Interfaces determine Tables / Store Points
	Posts interface {
		//Defining Methods
		Create(context.Context, *Post) error
		GetByID(context.Context, int64) (*Post, error)
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
