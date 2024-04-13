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

	"os"
	"testing"

	"avito/app/internal/repository"
	"avito/app/internal/service"
	"github.com/stretchr/testify/suite"
)

type APITestSuite struct {
	suite.Suite
	router     *thttp.Router
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
		s.FailNow("Ошибка инициализации Postgres", err)
	}
	s.redClient, err = rediscl.New(ctx)
	if err != nil {
		s.FailNow("Ошибка инициализации Redis", err)
	}

	s.initDeps()

}

func (s *APITestSuite) TearDownSuite() {
	s.pgClient.DB.Close()
	err := s.redClient.DB.Close()
	if err != nil {
		logrus.Errorln("Redis закрылся с ошибкой:", err)

	}
	s.cancel()
}

func (s *APITestSuite) initDeps() {

	repositoryField := repository.NewRepository(s.pgClient, s.redClient)
	serviceField := service.NewService(repositoryField)

	router, err := thttp.NewRouter(s.ctx, serviceField)
	if err != nil {
		s.FailNow("Failed to initialize", err)
	}

	s.router = router
	s.service = serviceField
	s.repository = repositoryField

}

func TestMain(m *testing.M) {
	rc := m.Run()
	os.Exit(rc)
}
