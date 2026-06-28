package cache

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	client    *redis.Client
	keyPrefix string
}

var _ DialogCache = (*RedisClient)(nil)

func NewRedisClient(addr string) *RedisClient {
	client := redis.NewClient(&redis.Options{
		Addr: addr,
	})
	return &RedisClient{client: client, keyPrefix: "dialog"}
}

func (r *RedisClient) Get(ctx context.Context, key string) (string, error) {
	return r.client.Get(ctx, key).Result()
}

func (r *RedisClient) Set(ctx context.Context, key string, value interface{}) error {
	return r.client.Set(ctx, key, value, 0).Err()
}

func (r *RedisClient) SaveSession(ctx context.Context, session *DialogSession) error {
	if session == nil {
		return fmt.Errorf("session is nil")
	}

	if session.DialogData == nil {
		session.DialogData = map[string]interface{}{}
	}

	payload, err := json.Marshal(session)
	if err != nil {
		return fmt.Errorf("marshal dialog session: %w", err)
	}

	key := r.buildSessionKey(session.UserID)
	return r.client.Set(ctx, key, payload, 0).Err()
}

func (r *RedisClient) GetSession(ctx context.Context, userID int64) (*DialogSession, error) {
	key := r.buildSessionKey(userID)
	rawValue, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, ErrSessionNotFound
		}
		return nil, fmt.Errorf("get dialog session from redis: %w", err)
	}

	var session DialogSession
	if err := json.Unmarshal([]byte(rawValue), &session); err != nil {
		return nil, fmt.Errorf("unmarshal dialog session: %w", err)
	}

	if session.DialogData == nil {
		session.DialogData = map[string]interface{}{}
	}

	return &session, nil
}

func (r *RedisClient) DeleteSession(ctx context.Context, userID int64) error {
	key := r.buildSessionKey(userID)
	if err := r.client.Del(ctx, key).Err(); err != nil {
		return fmt.Errorf("delete dialog session from redis: %w", err)
	}
	return nil
}

func (r *RedisClient) buildSessionKey(userID int64) string {
	return fmt.Sprintf("%s:%d", r.keyPrefix, userID)
}

func (r *RedisClient) Close() error {
	return r.client.Close()
}
