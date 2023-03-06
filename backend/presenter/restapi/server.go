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
	elog "github.com/labstack/gommon/log"
	"golang.org/x/sync/errgroup"
)

type Server struct {
	echo   *echo.Echo
	authUc domain.AuthUsecase
}

func NewServer(
	cfg *config.Config, repo domain.Repository, kvs domain.KVS,
) Server {
	e := echo.New()
	a := usecase.NewAuthUsecase(repo, kvs, cfg.AuthTokenMaxAge)
	s := Server{
		echo:   e,
		authUc: a,
	}

	h := newHandler(cfg.RunMode, repo, a)
	s.registerRoutes(h)
	s.setupGlobalMiddleware(cfg)

	switch cfg.RunMode {
	case config.ModeDebug:
		s.echo.Logger.SetLevel(elog.DEBUG)
	case config.ModeRelease:
		s.echo.Logger.SetLevel(elog.INFO)
	}
	return s
}

// Run は Graceful shutdown を有効にしてサーバを起動する。
func (s *Server) Run(ctx context.Context, address string) error {
	ctx, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		err := s.echo.Start(address)
		if err != nil && err != http.ErrServerClosed {
			log.Printf("failed to close: %#v", err)
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

// Routes は登録済みのルーティング情報の一覧を返す。
func (s *Server) Routes() []*echo.Route {
	return s.echo.Routes()
}

func (s *Server) setupGlobalMiddleware(cfg *config.Config) {
	s.echo.Use(
		CORSMiddleware(cfg),
	)
}
