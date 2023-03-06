package main

import (
	"context"
	"log"
	"time"

	"github.com/arumakan1727/todo-app-go-react/clock"
	"github.com/arumakan1727/todo-app-go-react/config"
	"github.com/arumakan1727/todo-app-go-react/presenter/restapi"
	"github.com/arumakan1727/todo-app-go-react/repository/pgsql"
	"github.com/arumakan1727/todo-app-go-react/repository/redis"
)

func CmdServe(address string) {
	serve := func(ctx context.Context, address string) error {
		cfg, err := config.NewFromEnv()
		if err != nil {
			log.Fatalf("cannot configure from env: %v", err)
		}

		clk := clock.GetRealClocker(time.UTC)

		repo, closeRepo, err := pgsql.NewRepository(ctx, cfg, clk)
		if err != nil {
			return err
		}
		defer closeRepo()

		kvs, err := redis.NewKVS(ctx, cfg)
		if err != nil {
			return err
		}

		s := restapi.NewServer(cfg, repo, kvs)
		return s.Run(ctx, address)
	}

	if err := serve(context.Background(), address); err != nil {
		log.Fatal(err)
	}
}
