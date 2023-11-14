package repository

import (
	"context"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

type IRedisRepository interface {
	Set(ctx context.Context, key string, value *string, time time.Duration) error
	Get(ctx context.Context, key string) (*string, error)
}

type RedisRepository struct {
	redis *redis.Client
}

func NewRedisRepository(redis *redis.Client) IRedisRepository {
	return &RedisRepository{redis}
}

func (r *RedisRepository) Set(ctx context.Context, key string, value *string, time time.Duration) error {
	_, err := r.redis.Set(ctx, key, value, time).Result()
	if err != nil {
		log.Printf("%v", err)
		return err
	}
	return nil
}

func (r *RedisRepository) Get(ctx context.Context, key string) (*string, error) {
	value, err := r.redis.Get(ctx, key).Result()
	if err != nil {
		log.Printf("%v", err)
		return nil, err
	}
	return &value, nil

}
