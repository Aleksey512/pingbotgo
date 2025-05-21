package redis

import (
	"context"
	"errors"
	"fmt"
	"misbotgo/internal/config"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisStorage struct {
	client *redis.Client
}

func NewRedisStorage(cfg *config.Settings) (*RedisStorage, error) {
	options, err := redis.ParseURL(cfg.RedisURL())
	if err != nil {
		return nil, fmt.Errorf("failed to parse Redis URL: %w", err)
	}

	client := redis.NewClient(options)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if _, err := client.Ping(ctx).Result(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return &RedisStorage{client: client}, nil
}

func (r *RedisStorage) Close() error {
	if r.client != nil {
		return r.client.Close()
	}
	return nil
}

func (r *RedisStorage) AddChatID(ctx context.Context, chatID string) error {
	if r.client == nil {
		return errors.New("Redis client is not initialized")
	}
	return r.client.SAdd(ctx, "CHAT_IDS", chatID).Err()
}

func (r *RedisStorage) RemoveChatID(ctx context.Context, chatID string) error {
	if r.client == nil {
		return errors.New("Redis client is not initialized")
	}
	return r.client.SRem(ctx, "CHAT_IDS", chatID).Err()
}

func (r *RedisStorage) GetChatIDs(ctx context.Context) ([]string, error) {
	if r.client == nil {
		return nil, errors.New("Redis client is not initialized")
	}
	return r.client.SMembers(ctx, "CHAT_IDS").Result()
}
