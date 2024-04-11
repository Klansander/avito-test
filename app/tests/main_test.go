package tests

import (
	"avito/app/internal/transport/thttp"
	"avito/app/pkg/config"
	pc "avito/app/pkg/context"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"log"

	"avito/app/pkg/postgresql"
	rediscl "avito/app/pkg/redis"
	"context"

	"net/http"
	"os"
	"testing"

	"avito/app/internal/repository"
	"avito/app/internal/service"
	"github.com/stretchr/testify/suite"
)

type APITestSuite struct {
	suite.Suite

	router     *thttp.Router
	httpServer *http.Server
	pgClient   *postgresql.Postgres
	service    *service.Service
	redClient  *rediscl.Redis
	repository *repository.Repository
	ctx        context.Context
	cancel     context.CancelFunc
}

func TestAPISuite(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	suite.Run(t, new(APITestSuite))
}

func (s *APITestSuite) SetupSuite() {

	ctx, cancel := context.WithCancel(context.Background())

	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
	logrus.Info("Инициализация конфигурации")

	ctx = pc.AddConfig(ctx, config.New())

	s.ctx = ctx
	s.cancel = cancel

	var err error
	s.pgClient, err = postgresql.New(ctx)
	if err != nil {
		s.FailNow("Failed to initialize token manager", err)
	}
	s.redClient, err = rediscl.New(ctx)
	if err != nil {
		s.FailNow("Failed to initialize token manager", err)
	}

	s.initDeps()

}

func (s *APITestSuite) TearDownSuite() {
	s.pgClient.DB.Close()
	s.redClient.DB.Close()
	s.cancel()
	//nolint:errcheck
}

func (s *APITestSuite) initDeps() {

	repository := repository.NewRepository(s.pgClient, s.redClient)
	service := service.NewService(repository)

	router, err := thttp.NewRouter(s.ctx, service)
	if err != nil {
		s.FailNow("Failed to initialize token manager", err)
	}

	s.router = router
	s.service = service
	s.repository = repository

}

func TestMain(m *testing.M) {
	rc := m.Run()
	os.Exit(rc)
}
