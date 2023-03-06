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
		GetClientAuthFromCtx: getAuthMaterialFromCtx,
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

	// GET /ping は /v1 ルート外でも利用可能とする
	s.echo.GET("/ping", h.GetPing)

	e := s.echo.Group("/v1")

	//-----------------------------------------------
	// No auth group
	{
		e.GET("/ping", h.GetPing)
		e.POST("/authtoken/new", h.CreateAuthToken)
		e.POST("/users", h.CreateUser)
	}

	//-----------------------------------------------
	// Normal user auth group
	e = e.Group("", AuthMiddleware(s.runMode, s.authUc))
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
