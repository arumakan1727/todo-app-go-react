package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/arumakan1727/todo-app-go-react/clock"
	"github.com/arumakan1727/todo-app-go-react/config"
	"github.com/arumakan1727/todo-app-go-react/presenter/restapi"
	"github.com/arumakan1727/todo-app-go-react/repository/pgsql"
	"github.com/arumakan1727/todo-app-go-react/repository/redis"
)

func CmdServe(host string, port uint) {
	launch := func(ctx context.Context) error {
		cfg, err := config.NewFromEnv()
		if err != nil {
			return fmt.Errorf("cannot configure from env: %v", err)
		}

		listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", host, port))
		if err != nil {
			return fmt.Errorf("failed to listen at TCP port %d on %s: %v", port, host, err)
		}

		clk := clock.GetRealClocker(time.UTC)
		repo, err := pgsql.NewRepository(ctx, cfg, clk)
		if err != nil {
			return err
		}
		defer repo.Close()

		kvs, err := redis.NewKVS(ctx, cfg)
		if err != nil {
			return err
		}
		defer kvs.Close()

		s := restapi.NewServer(cfg, repo, kvs)
		return s.Serve(ctx, listener)
	}

	if err := launch(context.Background()); err != nil {
		log.Fatal(err)
	}
}
