package main

import (
	"fmt"
	"os"

	"github.com/arumakan1727/todo-app-go-react/config"
	"github.com/k0kubun/pp/v3"
)

func CmdDumpConfig() {
	cfg, err := config.LoadFromEnv()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load config from env: %v", err)
	}

	pp.Println(cfg)
}
