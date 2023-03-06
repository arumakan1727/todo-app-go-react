package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/arumakan1727/todo-app-go-react/config"
	"github.com/arumakan1727/todo-app-go-react/domain"
	"github.com/redis/go-redis/v9"
)

type kvs struct {
	cli *redis.Client
}

func NewKVS(ctx context.Context, cfg *config.Config) (domain.KVS, error) {
	cli := redis.NewClient(&redis.Options{
		Addr: cfg.RedisAddr,
	})

	if err := cli.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("redis.NewKVS: failed to connect to redis: %w", err)
	}
	return &kvs{cli: cli}, nil
}

func (kvs *kvs) SaveAuth(
	ctx context.Context, a domain.AuthToken, uid domain.UserID, expiration time.Duration,
) error {
	return kvs.cli.Set(ctx, string(a), int(uid), expiration).Err()
}

func (kvs *kvs) FetchAuth(ctx context.Context, a domain.AuthToken) (domain.UserID, error) {
	uid, err := kvs.cli.Get(ctx, string(a)).Uint64()
	if err != nil {
		return 0, fmt.Errorf("failed to get userID by key='%s': %w", a, domain.ErrNotFound)
	}
	return domain.UserID(uid), nil
}

func (kvs *kvs) DeleteAuth(ctx context.Context, a domain.AuthToken) error {
	return kvs.cli.Del(ctx, string(a)).Err()
}
