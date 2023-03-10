package restapi

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"os/signal"
	"syscall"

	"github.com/arumakan1727/todo-app-go-react/config"
	"github.com/arumakan1727/todo-app-go-react/domain"
	"github.com/arumakan1727/todo-app-go-react/usecase"
	"github.com/labstack/echo/v4"
	elog "github.com/labstack/gommon/log"
	"golang.org/x/sync/errgroup"
)

type Server struct {
	echo    *echo.Echo
	runMode config.RunMode
	repo    domain.Repository
	kvs     domain.KVS
}

func NewServer(
	cfg *config.Config, repo domain.Repository, kvs domain.KVS,
) *Server {
	e := echo.New()
	au := usecase.NewAuthUsecase(repo, kvs, cfg.AuthTokenMaxAge)
	s := &Server{
		echo:    e,
		runMode: cfg.RunMode,
		repo:    repo,
		kvs:     kvs,
	}

	h := NewHandler(cfg.RunMode, repo, au)
	registerRoutes(s.echo, h, AuthMiddleware(cfg.RunMode, au))
	s.setupGlobalMiddleware(cfg)

	switch cfg.RunMode {
	case config.ModeDebug:
		s.echo.Debug = true
		s.echo.Logger.SetLevel(elog.DEBUG)
	case config.ModeRelease:
		s.echo.Debug = false
		s.echo.Logger.SetLevel(elog.INFO)
	}
	s.echo.HTTPErrorHandler = s.errorHandler
	return s
}

func (s *Server) HideBanner(hide bool) {
	s.echo.HideBanner = hide
}

func (s *Server) Close() {
	if s.repo != nil {
		s.repo.Close()
	}
	if s.kvs != nil {
		s.kvs.Close()
	}
	if s.echo.Listener != nil {
		_ = s.echo.Listener.Close()
	}
}

// Serve は Graceful shutdown を有効にしてサーバを起動する。
func (s *Server) Serve(ctx context.Context, l net.Listener) error {
	defer s.Close()

	ctx, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		s.echo.Listener = l
		err := s.echo.Start("")
		if err != nil && err != http.ErrServerClosed {
			log.Printf("(*restapi.Server).Serve: failed to close: %#v", err)
			return err
		}
		return nil
	})

	<-ctx.Done()
	if err := s.echo.Shutdown(context.Background()); err != nil {
		log.Printf("failed to shutdown: %+v", err)
	}

	// wait for graceful shutdown
	return eg.Wait()
}

func (s *Server) setupGlobalMiddleware(cfg *config.Config) {
	s.echo.Use(
		CORSMiddleware(cfg),
	)
}

// domain.Err* がHTTPステータスコードの何番に対応するか
var pairsDomainErrAndHTTPStatus = []struct {
	e      error
	status int
}{
	{e: domain.ErrAlreadyExits, status: http.StatusConflict},
	{e: domain.ErrEmptyPatch, status: http.StatusBadRequest},
	{e: domain.ErrIncorrectEmailOrPasswd, status: http.StatusUnauthorized},
	{e: domain.ErrInvalidInput, status: http.StatusBadRequest},
	{e: domain.ErrNotFound, status: http.StatusNotFound},
	{e: domain.ErrNotInTransaction, status: http.StatusInternalServerError},
	{e: domain.ErrUnauthorized, status: http.StatusUnauthorized},
}

func (s *Server) errorHandler(err error, c echo.Context) {
	if c.Response().Committed {
		return
	}

	if he, ok := err.(*echo.HTTPError); ok {
		if c.Request().Method == http.MethodHead {
			if err := c.NoContent(he.Code); err != nil {
				s.echo.Logger.Error(err)
			}
			return
		}
		if err := c.JSON(he.Code, &he); err != nil {
			s.echo.Logger.Error(err)
		}
		return
	}

	var resp struct {
		status int
		body   RespSimpleError
	}
	for _, pair := range pairsDomainErrAndHTTPStatus {
		if errors.Is(err, pair.e) {
			resp.status = pair.status
			resp.body = RespSimpleError{
				Message: err.Error(),
			}
			break
		}
	}
	if resp.status == 0 {
		resp.status = http.StatusInternalServerError
		resp.body = RespSimpleError{
			Message: http.StatusText(http.StatusInternalServerError),
		}
		s.echo.Logger.Error(err)
	}
	if err := c.JSON(resp.status, &resp.body); err != nil {
		s.echo.Logger.Error(err)
	}
}
