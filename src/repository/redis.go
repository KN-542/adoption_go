package repository

import (
	"context"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

type IRedisRepository interface {
	Set(
		ctx context.Context,
		hashKey string,
		key string,
		value *string,
		ttl time.Duration,
	) error
	Get(ctx context.Context, hashKey string, key string) (*string, error)
}

type RedisRepository struct {
	redis *redis.Client
}

func NewRedisRepository(redis *redis.Client) IRedisRepository {
	return &RedisRepository{redis}
}

func (r *RedisRepository) Set(
	ctx context.Context,
	hashKey string,
	key string,
	value *string,
	ttl time.Duration,
) error {
	_, err := r.redis.HSet(ctx, hashKey, key, *value).Result()
	if err != nil {
		log.Printf("%v", err)
		return err
	}

	if ttl > 0 {
		_, err = r.redis.Expire(ctx, hashKey, ttl).Result()
		if err != nil {
			log.Printf("%v", err)
			return err
		}
	}

	return nil
}

func (r *RedisRepository) Get(ctx context.Context, hashKey string, key string) (*string, error) {
	value, err := r.redis.HGet(ctx, hashKey, key).Result()
	if err != nil {
		log.Printf("%v", err)
		return nil, err
	}
	return &value, nil

}
