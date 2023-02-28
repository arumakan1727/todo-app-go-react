package config

import (
	"github.com/caarlos0/env/v7"
)

type Config struct {
	ServePort     int    `env:"TODO_SERVE_PORT" envDefault:"8181"`
	PgSQLHost     string `env:"TODO_PGSQL_HOST" envDefault:"127.0.0.1"`
	PgSQLPort     int    `env:"TODO_PGSQL_PORT" emvDefault:"5432"`
	PgSQLUser     string `env:"TODO_PGSQL_USER" envDefault:"todoapp"`
	PgSQLPasswd   string `env:"TODO_PGSQL_PASSWORD" envDefault:"todoapp"`
	PgSQLDatabase string `env:"TODO_PGSQL_DATABASE" envDefault:"todoapp"`
	RedisHost     string `env:"TODO_REDIS_HOST" envDefault:"127.0.0.1"`
	RedisPort     int    `env:"TODO_REDIS_PORT" envDefault:"6379"`
}

func NewFromEnv() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
