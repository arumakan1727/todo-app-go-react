package restapi

import (
	"github.com/arumakan1727/todo-app-go-react/config"
	"github.com/arumakan1727/todo-app-go-react/domain"
	"github.com/arumakan1727/todo-app-go-react/usecase"
	"github.com/labstack/echo/v4"
)

func NewHandler(
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

type EchoRouter interface {
	GET(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	HEAD(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	POST(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PUT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PATCH(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	DELETE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	OPTIONS(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	Group(prefix string, m ...echo.MiddlewareFunc) *echo.Group
}

func registerRoutes(
	e EchoRouter, h *ServerInterfaceWrapper, authMiddleware echo.MiddlewareFunc,
) {
	with := func(e EchoRouter, prefix string, f func(*echo.Group)) {
		g := e.Group(prefix)
		f(g)
	}

	// GET /ping は /v1 ルート外でも利用可能とする
	e.GET("/ping", h.GetPing)

	e = e.Group("/v1")

	//-----------------------------------------------
	// No auth group
	{
		e.GET("/ping", h.GetPing)
		e.POST("/authtoken/new", h.CreateAuthToken)
		e.POST("/users", h.CreateUser)
	}

	//-----------------------------------------------
	// Normal user auth group
	e = e.Group("", authMiddleware)
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

func Routes() []*echo.Route {
	e := echo.New()
	h := NewHandler(config.ModeDebug, nil, nil)
	fakeAuthMiddleware := echo.MiddlewareFunc(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			return next(c)
		}
	})
	registerRoutes(e, h, fakeAuthMiddleware)
	return e.Routes()
}
