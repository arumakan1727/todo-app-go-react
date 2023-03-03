package restapi

import (
	"github.com/labstack/echo/v4"
)

func newHandler() ServerInterfaceWrapper {
	type allHandler struct {
		PingHandler
		UserHandler
		TaskHandler
	}

	return ServerInterfaceWrapper{
		Handler: allHandler{
			PingHandler{},
			UserHandler{},
			TaskHandler{},
		},
	}
}

func RegisterRoutes(e *echo.Echo) {
	type echoGrouper interface {
		Group(prefix string, m ...echo.MiddlewareFunc) *echo.Group
	}
	with := func(e echoGrouper, prefix string, f func(e *echo.Group)) {
		g := e.Group(prefix)
		f(g)
	}

	h := newHandler()

	e.GET("/ping", h.GetPing)

	with(e, "/authtoken", func(e *echo.Group) {
		e.POST("/new", h.CreateAuthToken)
	})

	with(e, "/users", func(e *echo.Group) {
		e.POST("/", h.CreateUser)
	})

	with(e, "/tasks", func(e *echo.Group) {
		e.GET("/", h.ListTasks)
		e.POST("/", h.CreateTask)

		with(e, "/:taskID", func(e *echo.Group) {
			e.GET("/", h.GetTask)
			e.PATCH("/", h.PatchTask)
			e.DELETE("/", h.DeleteTask)
		})
	})

	with(e, "/__", func(e *echo.Group) {
		e.GET("/users", h.ListUsersForAdmin)
	})
}
