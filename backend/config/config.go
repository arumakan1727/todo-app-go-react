package config

import (
	"fmt"
	"time"

	"github.com/caarlos0/env/v7"
)

type RunMode string

const (
	ModeDebug   RunMode = "debug"
	ModeRelease RunMode = "release"
)

type Config struct {
	RunMode         RunMode       `env:"TODO_RUN_MODE"`
	AuthTokenMaxAge time.Duration `env:"TODO_AUTH_TOKEN_MAX_AGE" envDefault:"24h"`

	PgSQLURL  string `env:"TODO_PGSQL_URL" envDefault:"postgres://todouser:todopass@127.0.0.1:5432/tododb?sslmode=disable"`
	RedisAddr string `env:"TODO_REDIS_ADDR" envDefault:"127.0.0.1:6379"`
}

func NewFromEnv() (*Config, error) {
	cfg := &Config{}
	opt := env.Options{
		RequiredIfNoDef: true,
	}
	if err := env.Parse(cfg, opt); err != nil {
		return nil, err
	}
	switch cfg.RunMode {
	case ModeDebug, ModeRelease:
		break
	default:
		return nil, fmt.Errorf("NewFromEnv: Invalid run mode: %q", cfg.RunMode)
	}
	return cfg, nil
}

func ForTesting() *Config {
	return &Config{
		RunMode:         ModeDebug,
		AuthTokenMaxAge: time.Minute * 1,
		PgSQLURL:        "postgres://todouser:todopass@127.0.0.1:25432/tododb__test?sslmode=disable",
		RedisAddr:       "127.0.0.1:26379",
	}
}
