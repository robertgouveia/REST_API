package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/robertgouveia/social/internal/store"
)

const UserExpTime = time.Minute

type UserStore struct {
	db *redis.Client
}

func (s *UserStore) Get(ctx context.Context, userID int64) (*store.User, error) {
	if s.db == nil {
		return nil, fmt.Errorf("database is not initialized")
	}

	cacheKey := fmt.Sprintf("user-%v", userID)
	data, err := s.db.Get(ctx, cacheKey).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // Cache miss
		}
		return nil, err // Other errors
	}

	if data == "" {
		return nil, nil // No user data in cache
	}

	var user store.User
	if err := json.Unmarshal([]byte(data), &user); err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *UserStore) Set(ctx context.Context, user *store.User) error {
	cacheKey := fmt.Sprintf("user-%v", user.ID)

	json, err := json.Marshal(user)
	if err != nil {
		return err
	}

	return s.db.SetEX(ctx, cacheKey, json, UserExpTime).Err()
}
