package store

import (
	"context"
	"database/sql"
)

// Defining the Store Dependencies
type PostsStore struct {
	db *sql.DB
}

// Methods
func (s *PostsStore) Create(ctx context.Context) error {
	return nil
}
