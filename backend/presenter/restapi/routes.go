package restapi

import (
	"github.com/arumakan1727/todo-app-go-react/config"
	"github.com/arumakan1727/todo-app-go-react/domain"
	"github.com/arumakan1727/todo-app-go-react/usecase"
	"github.com/labstack/echo/v4"
)

func newHandler(
	runMode config.RunMode, repo domain.Repository, authUc domain.AuthUsecase,
) *ServerInterfaceWrapper {
	userUc := usecase.NewUserUsecase(repo)
	taskUc := usecase.NewTaskUsecase(repo)

	type allHandler struct {
		PingHandler
		*AuthTokenHandler
		*UserHandler
		*TaskHandler
	}
	return &ServerInterfaceWrapper{
		Handler: allHandler{
			PingHandler{},
			&AuthTokenHandler{runMode, authUc},
			&UserHandler{userUc},
			&TaskHandler{taskUc},
		},
		GetClientUIDFromCtx: getUserIDFromCtx,
	}
}

func (s *Server) registerRoutes(h *ServerInterfaceWrapper) {
	type echoGrouper interface {
		Group(prefix string, m ...echo.MiddlewareFunc) *echo.Group
	}
	with := func(e echoGrouper, prefix string, f func(*echo.Group)) {
		g := e.Group(prefix)
		f(g)
	}

	//-----------------------------------------------
	// No auth group
	{
		e := s.srv
		e.GET("/ping", h.GetPing)
		e.POST("/authtoken/new", h.CreateAuthToken)
		e.POST("/users", h.CreateUser)
	}

	//-----------------------------------------------
	// Normal user auth group
	e := s.srv.Group("", AuthMiddleware(s.authUc))
	with(e, "/tasks", func(e *echo.Group) {
		e.GET("", h.ListTasks)
		e.POST("", h.CreateTask)

		with(e, "/:taskID", func(e *echo.Group) {
			e.GET("", h.GetTask)
			e.PATCH("", h.PatchTask)
			e.DELETE("", h.DeleteTask)
		})
	})

	//-----------------------------------------------
	// admin only group
	e = e.Group("/__", AdminOnlyMiddleware)
	e.GET("/users", h.ListUsersForAdmin)
}
