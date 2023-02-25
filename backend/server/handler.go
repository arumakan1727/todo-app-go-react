package server

import "github.com/arumakan1727/todo-app-go-react/server/handler"

type Handler struct {
	handler.PingHandler
	handler.AuthTokenHandler
	handler.UserHandler
	handler.TaskHandler
}

// type check for implementation of ServerInterface
var _ ServerInterface = &Handler{}
