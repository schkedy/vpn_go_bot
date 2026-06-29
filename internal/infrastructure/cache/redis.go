package cache

import (
	"context"
	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	client    *redis.Client
	keyPrefix string
}

func NewRedisClient(addr string) *RedisClient {
	client := redis.NewClient(&redis.Options{
		Addr: addr,
	})
	return &RedisClient{client: client, keyPrefix: "go-dialog"}
}

func (r *RedisClient) Get(ctx context.Context, key string) (interface{}, error) {
	return r.client.Get(ctx, r.prefixedKey(key)).Result()
}
func (r *RedisClient) Set(ctx context.Context, key string, value interface{}) error {
	return r.client.Set(ctx, r.prefixedKey(key), value, 0).Err()
}

func (r *RedisClient) Delete(ctx context.Context, key string) error {
	return r.client.Del(ctx, r.prefixedKey(key)).Err()
}

func (r *RedisClient) HSet(ctx context.Context, key string, field string, value interface{}) error {
	return r.client.HSet(ctx, r.prefixedKey(key), field, value).Err()
}

func (r *RedisClient) HGet(ctx context.Context, key string, field string) (string, error) {
	return r.client.HGet(ctx, r.prefixedKey(key), field).Result()
}

func (r *RedisClient) HGetAll(ctx context.Context, key string) (map[string]string, error) {
	return r.client.HGetAll(ctx, r.prefixedKey(key)).Result()
}
func (r *RedisClient) HDel(ctx context.Context, key string, field string) error {
	return r.client.HDel(ctx, r.prefixedKey(key), field).Err()
}

func (r *RedisClient) prefixedKey(key string) string {
	return r.keyPrefix + ":" + key
}



func (r *RedisClient) Close() error {
	return r.client.Close()
}
