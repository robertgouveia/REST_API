package store

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

var (
	ErrNotFound          = errors.New("record not found")
	ErrConflict          = errors.New("client conflict in versions")
	QueryTimeoutDuration = time.Second * 5
	ErrDuplicateEmail    = errors.New("email already exists")
	ErrDuplicateUsername = errors.New("username already exists")
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
		GetUserFeed(context.Context, int64, PaginatedFeedQuery) ([]PostWithMetaData, error)
	}

	Users interface {
		Create(context.Context, *sql.Tx, *User) error
		GetByID(context.Context, int64) (*User, error)
		CreateAndInvite(context.Context, *User, string, time.Duration) error
		Activate(context.Context, string) error
		Delete(context.Context, int64) error
		GetByEmail(context.Context, string) (*User, error)
	}

	Comments interface {
		Create(context.Context, *Comment) error
		GetByPostID(context.Context, int64) ([]Comment, error)
	}

	Followers interface {
		Follow(context.Context, int64, int64) error
		Unfollow(context.Context, int64, int64) error
	}

	Roles interface {
		GetByName(context.Context, string) (*Role, error)
	}
}

// Defining a Store and supplying the dependencies
func NewStorage(db *sql.DB) Storage {
	//Creating and returning a Storage object with Repository References
	return Storage{
		Posts:     &PostStore{db},
		Users:     &UserStore{db},
		Comments:  &CommentStore{db},
		Followers: &FollowerStore{db},
		Roles:     &RoleStore{db},
	}
}

func withTx(db *sql.DB, ctx context.Context, fn func(*sql.Tx) error) error {
	//begin creates a transaction which can be rolled back
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	//send the transaction to the function
	if err := fn(tx); err != nil {
		//if an error occurs we can rollback
		_ = tx.Rollback()
		return err
	}

	//completes the transation
	return tx.Commit()
}
