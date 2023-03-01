package config

import (
	"fmt"

	"github.com/caarlos0/env/v7"
)

type RunMode string

const (
	ModeDebug   RunMode = "debug"
	ModeRelease RunMode = "release"
)

type Config struct {
	RunMode   RunMode `env:"TODO_RUN_MODE"`
	ServePort int     `env:"TODO_SERVE_PORT"`

	PgSQLHost     string `env:"TODO_PGSQL_HOST"`
	PgSQLPort     int    `env:"TODO_PGSQL_PORT" envDefault:"5432"`
	PgSQLUser     string `env:"TODO_PGSQL_USER" envDefault:"todouser"`
	PgSQLPasswd   string `env:"TODO_PGSQL_PASSWORD" envDefault:"todopass"`
	PgSQLDatabase string `env:"TODO_PGSQL_DATABASE" envDefault:"tododb"`

	RedisHost string `env:"TODO_REDIS_HOST"`
	RedisPort int    `env:"TODO_REDIS_PORT" envDefault:"6379"`
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

var testConfig = Config{}

func ForTesting() *Config {
	return &Config{
		RunMode:       ModeDebug,
		ServePort:     8181,
		PgSQLHost:     "127.0.0.1",
		PgSQLPort:     25432,
		PgSQLUser:     "todouser",
		PgSQLPasswd:   "todopass",
		PgSQLDatabase: "tododb",
		RedisHost:     "127.0.0.1",
		RedisPort:     26379,
	}
}
