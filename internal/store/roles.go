package store

import (
	"context"
	"database/sql"
)

type Role struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Level       int    `json:"level"`
}

type RoleStore struct {
	db *sql.DB
}

func (s *RoleStore) GetByName(ctx context.Context, role string) (*Role, error) {
	query := `
		SELECT id, name, description, level FROM roles WHERE name = $1
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	Role := &Role{}
	err := s.db.QueryRowContext(ctx, query, role).Scan(&Role.ID, &Role.Name, &Role.Description, &Role.Level)
	if err != nil {
		return nil, err
	}

	return Role, nil
}
