package store

import (
	"context"
	"database/sql"
	"errors"
)

var (
	ErrNotFound = errors.New("Record Not Found")
	ErrConflict = errors.New("Client Conflict in Versions")
)

// Repository Pattern for decoupling
type Storage struct {
	//Interfaces determine Tables / Store Points
	Posts interface {
		//Defining Methods
		Create(context.Context, *Post) error
		GetByID(context.Context, int64) (*Post, error)
		Delete(context.Context, int64) error
		Update(context.Context, *Post) error
	}

	Users interface {
		Create(context.Context, *User) error
	}

	Comments interface {
		GetByPostID(context.Context, int64) ([]Comment, error)
	}
}

// Defining a Store and supplying the dependencies
func NewStorage(db *sql.DB) Storage {
	//Creating and returning a Storage object with Repository References
	return Storage{
		Posts:    &PostStore{db},
		Users:    &UserStore{db},
		Comments: &CommentStore{db},
	}
}
