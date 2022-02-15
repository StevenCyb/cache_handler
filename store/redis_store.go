package store

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

// RedisStore uses Redis
type RedisStore struct {
	Client     *redis.Client
	expiration time.Duration
}

// NewRedisStore creates a new RedisStore.
// Use empty string for username & password if not password required.
// Set port to 0 to use default.
func NewRedisStore(endpoint string, port int, username, password string, expiration time.Duration) *RedisStore {
	return &RedisStore{
		Client: redis.NewClient(&redis.Options{
			Addr:     endpoint,
			Username: username,
			Password: password,
			DB:       port,
		}),
		expiration: expiration,
	}
}

// Get data from store with given key
func (store RedisStore) Get(key string) ([]byte, error) {
	return store.Client.Get(context.TODO(), key).Bytes()
}

// Put data to store fore given key
func (store RedisStore) Set(key string, data []byte) error {
	statusCmd := store.Client.Set(context.TODO(), key, data, store.expiration)
	_, err := statusCmd.Result()
	return err
}
