package sessions

import (
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
)

type Redis struct {
	Client *redis.Client
}

func NewSessions() *Redis {
	return &Redis{
		Client: redis.NewClient(&redis.Options{
			Addr:     "localhost:6379",
			Password: "",
			DB:       1,
		}),
	}
}

func (r *Redis) GetSession(sid string) (string, error) {
	// Returns whether the user exists in the database or not
	_, err := r.Client.Get(context.Background(), sid).Result()
	if err != nil {
		return "", err
	}
	return sid, err
}

func (r *Redis) CreateSession(username string) (string, error) {
	// nil error means it works
	sid := uuid.New().String()
	err := r.Client.Set(context.Background(), sid, username, 0).Err()
	if err != nil {
		return sid, err
	}

	return sid, nil
}
