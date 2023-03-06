package restapi

import (
	"context"
	"log"
	"net/http"
	"os/signal"
	"syscall"

	"github.com/arumakan1727/todo-app-go-react/config"
	"github.com/arumakan1727/todo-app-go-react/domain"
	"github.com/arumakan1727/todo-app-go-react/usecase"
	"github.com/labstack/echo/v4"
	"golang.org/x/sync/errgroup"
)

type Server struct {
	srv    *echo.Echo
	authUc domain.AuthUsecase
}

func NewServer(
	cfg *config.Config, repo domain.Repository, kvs domain.KVS,
) Server {
	e := echo.New()
	a := usecase.NewAuthUsecase(repo, kvs, cfg.AuthTokenMaxAge)
	s := Server{
		srv: e,
		authUc: a,
	}

	h := newHandler(cfg.RunMode, repo, a)
	s.registerRoutes(h)

	return s
}

// Run は Graceful shutdown を有効にしてサーバを起動する。
func (s *Server) Run(ctx context.Context, address string) error {
	ctx, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		err := s.srv.Start(address)
		if err != nil && err != http.ErrServerClosed {
			log.Printf("failed to close: %#v", err)
			return err
		}
		return nil
	})

	<-ctx.Done()
	if err := s.srv.Shutdown(context.Background()); err != nil {
		log.Printf("failed to shutdown: %+v", err)
	}

	// wait for graceful shutdown
	return eg.Wait()
}

// Routes は登録済みのルーティング情報の一覧を返す。
func (s *Server) Routes() []*echo.Route {
	return s.srv.Routes()
}
