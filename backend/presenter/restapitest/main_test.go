package restapitest

import (
	"context"
	"fmt"
	"log"
	"net"
	"testing"
	"time"

	"github.com/arumakan1727/todo-app-go-react/clock"
	"github.com/arumakan1727/todo-app-go-react/config"
	"github.com/arumakan1727/todo-app-go-react/domain"
	"github.com/arumakan1727/todo-app-go-react/presenter/restapi"
	"github.com/arumakan1727/todo-app-go-react/presenter/restapi/client"
	"github.com/arumakan1727/todo-app-go-react/repository/pgsql"
	"github.com/arumakan1727/todo-app-go-react/repository/redis"
	"github.com/arumakan1727/todo-app-go-react/usecase"
	"golang.org/x/sync/errgroup"
)

type UserData struct {
	domain.User
	Passwd string
	Token  domain.AuthToken
}

var (
	gClient *client.Client
	gRepo   domain.Repository
	gKVS    domain.KVS

	gUser  UserData
	gAdmin UserData

	gClock = clock.GetFixedClocker()
	gCtx   = context.Background()
)

func TestMain(m *testing.M) {
	s, l, err := NewServerForTesting(gClock)
	if err != nil {
		log.Fatal(err)
	}

	{
		gClient, err = client.NewClient(fmt.Sprintf("http://%s/v1", l.Addr()))
		if err != nil {
			l.Close()
			log.Fatal(err)
		}
	}

	ctx, cancel := context.WithCancel(context.Background())
	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		s.HideBanner(true)
		return s.Serve(ctx, l)
	})

	if err = gRepo.TruncateAll(ctx); err != nil {
		s.Close()
		log.Fatalf("TestMain: %+v", err)
	}

	if gUser, gAdmin, err = PrepareUserDatum(ctx, gRepo, gKVS); err != nil {
		s.Close()
		log.Fatalf("TestMain: PrepareUserDatum: %+v", err)
	}

	m.Run()

	cancel()
	if err := eg.Wait(); err != nil {
		log.Fatal(err)
	}
}

func NewServerForTesting(clk clock.Clocker) (*restapi.Server, net.Listener, error) {
	ctx := context.Background()
	cfg := config.ForTesting()

	var err error
	gRepo, err = pgsql.NewRepository(ctx, cfg, clk)
	if err != nil {
		return nil, nil, err
	}

	gKVS, err = redis.NewKVS(ctx, cfg)
	if err != nil {
		gRepo.Close()
		return nil, nil, err
	}

	// 空いているポートを自動で選択してリッスンする
	l, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		gRepo.Close()
		gKVS.Close()
		return nil, nil, fmt.Errorf("failed to listen using automatically chosen port: %v", err)
	}

	s := restapi.NewServer(cfg, gRepo, gKVS)
	return s, l, nil
}

func PrepareUserDatum(ctx context.Context, r domain.Repository, kvs domain.KVS) (user UserData, admin UserData, err error) {
	user.Passwd = "passwd-user"
	admin.Passwd = "passwd-admin"

	uu := usecase.NewUserUsecase(r)
	if user.User, err = uu.Store(ctx, "user@example.com", user.Passwd, "displayName-user", "user"); err != nil {
		log.Printf("Cannot create normal user: %+v", err)
		return
	}
	if admin.User, err = uu.Store(ctx, "admin@example.com", admin.Passwd, "displayName-admin", "admin"); err != nil {
		log.Printf("Cannot create admin user: %+v", err)
		return
	}

	au := usecase.NewAuthUsecase(r, kvs, time.Minute)
	if user.Token, err = au.IssueAuthToken(ctx, user.Email, user.Passwd); err != nil {
		return
	}
	if admin.Token, err = au.IssueAuthToken(ctx, admin.Email, admin.Passwd); err != nil {
		return
	}

	err = nil
	return
}
