package main

import (
	"flag"
	"fmt"
	"os"
)

func printUsage() {
	const s = `USAGE:
    api <command>

COMMANDS:
	serve    Start API server
	routes   Print API routes

OPTIONS:`

	fmt.Println(s)
	flag.PrintDefaults()
}

func main() {
	var (
		flagPort uint
		flagHost string
	)
	flag.Usage = printUsage
	flag.UintVar(&flagPort, "port", 8181, "TCP port of API server")
	flag.StringVar(&flagHost, "host", "localhost", "Host name or IP address of API server")
	flag.Parse()

	switch flag.Arg(0) {
	case "serve":
		CmdServe(flagHost, flagPort)
	case "routes":
		CmdRoutes()
	default:
		printUsage()
		os.Exit(1)
	}
}
