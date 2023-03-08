package redis

import (
	"bytes"
	"context"
	"encoding/gob"
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

func (kvs *kvs) Close() {
	if kvs == nil || kvs.cli == nil {
		return
	}
	_ = kvs.cli.Close()
}

func (kvs *kvs) SaveAuth(
	ctx context.Context, key domain.AuthToken, am *domain.AuthMaterial, expiration time.Duration,
) error {
	buf := bytes.NewBuffer(nil)
	if err := gob.NewEncoder(buf).Encode(am); err != nil {
		return fmt.Errorf("redis.KVS.SaveAuth: failed to serialize KVSAuthRecord: %w", err)
	}
	if err := kvs.cli.Set(ctx, string(key), buf.Bytes(), expiration).Err(); err != nil {
		return fmt.Errorf("redis.KVS.SaveAuth: failed to set auth: %w", err)
	}
	return nil
}

func (kvs *kvs) FetchAuth(ctx context.Context, key domain.AuthToken) (domain.AuthMaterial, error) {
	var res domain.AuthMaterial
	bs, err := kvs.cli.Get(ctx, string(key)).Bytes()
	if err != nil {
		return res, fmt.Errorf("redis.KVS.FetchAuth: failed to get auth by key='%s': %w", key, domain.ErrNotFound)
	}

	if err := gob.NewDecoder(bytes.NewReader(bs)).Decode(&res); err != nil {
		return res, fmt.Errorf("redis.KVS.FetchAuth: failed to deserialize: %w", err)
	}
	return res, nil
}

func (kvs *kvs) DeleteAuth(ctx context.Context, a domain.AuthToken) error {
	return kvs.cli.Del(ctx, string(a)).Err()
}
