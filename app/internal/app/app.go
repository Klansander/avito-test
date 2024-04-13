package app

import (
	"avito/app/internal/repository"
	"avito/app/internal/service"
	"avito/app/internal/transport/thttp"
	pc "avito/app/pkg/context"
	"avito/app/pkg/cron"
	"avito/app/pkg/errors"
	"avito/app/pkg/postgresql"
	"avito/app/pkg/redis"
	"context"
	"fmt"
	"github.com/go-co-op/gocron/v2"
	"github.com/rs/cors"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
	"net"
	"net/http"
)

type App struct {
	router     *thttp.Router
	httpServer *http.Server
	pgClient   *postgresql.Postgres
	service    *service.Service
	redClient  *rediscl.Redis
	Cron       *cron2.Cron
	repository *repository.Repository
}

func New(ctx context.Context) (*App, error) {

	logrus.Info("Инициализация подключения к базе данных")
	pgClient, err := postgresql.New(ctx)
	if err != nil {
		return nil, errors.Wrap(err)
	}

	redClient, err := rediscl.New(ctx)
	if err != nil {
		return nil, errors.Wrap(err)
	}

	repositoryField := repository.NewRepository(pgClient, redClient)
	serviceField := service.NewService(repositoryField)

	cron, err := cron2.NewCron(ctx)
	if err != nil {
		return nil, errors.Wrap(err)
	}

	router, err := thttp.NewRouter(ctx, serviceField)
	if err != nil {
		return nil, errors.Wrap(err)
	}

	return &App{
		router:     router,
		pgClient:   pgClient,
		redClient:  redClient,
		service:    serviceField,
		Cron:       cron,
		repository: repositoryField,
	}, nil

}

func (a *App) Run(ctx context.Context) error {

	g, ctx := errgroup.WithContext(ctx)

	logrus.Info("Запуск HTTP сервера")

	g.Go(func() error {
		return a.startHTTP(ctx)
	})

	g.Go(func() error {
		return a.StartCron(ctx)
	})

	return g.Wait()

}

func (a *App) Stop(ctx context.Context) {

	if err := a.httpServer.Shutdown(ctx); err != nil {
		logrus.Errorf("HTTP shutdown error: %v", err)

		if err = a.httpServer.Close(); err != nil {
			logrus.Errorf("HTTP close error: %v", err)
		}
	}

	logrus.Info("HTTP shutdown")

	a.pgClient.Close()
	a.redClient.Close()

}

func (a *App) startHTTP(ctx context.Context) error {

	cfg := pc.GetConfig(ctx)

	logrus.WithFields(logrus.Fields{
		"Host": cfg.HTTP.Host,
		"Port": cfg.HTTP.Port,
	}).Info("Параметры подключения")

	c := cors.New(cors.Options{
		AllowedMethods:     cfg.HTTP.CORS.AllowedMethods,
		AllowedOrigins:     cfg.HTTP.CORS.AllowedOrigins,
		AllowCredentials:   cfg.HTTP.CORS.AllowCredentials,
		AllowedHeaders:     cfg.HTTP.CORS.AllowedHeaders,
		OptionsPassthrough: cfg.HTTP.CORS.OptionsPassthrough,
		ExposedHeaders:     cfg.HTTP.CORS.ExposedHeaders,
	})

	a.httpServer = &http.Server{
		Handler:           c.Handler(a.router.Router),
		WriteTimeout:      cfg.HTTP.WriteTimeout,
		ReadTimeout:       cfg.HTTP.ReadTimeout,
		IdleTimeout:       cfg.HTTP.IdleTimeout,
		ReadHeaderTimeout: cfg.HTTP.ReadHeaderTimeout,
	}

	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", cfg.HTTP.Host, cfg.HTTP.Port))
	if err != nil {
		return errors.Wrap(err)
	}

	if err = a.httpServer.Serve(listener); err != nil {
		switch {
		case errors.Is(err, http.ErrServerClosed):
			logrus.Info("Сервер остановлен")
		default:
			return errors.Wrap(err)
		}
	}
	return nil

}

func (a *App) StartCron(ctx context.Context) error {

	s := *a.Cron.Scheduler
	interval := a.Cron.Interval

	_, err := s.NewJob(gocron.DurationJob(interval), gocron.NewTask(a.repository.Banner.SaveVersionBanner, ctx))
	if err != nil {
		return errors.Wrap(err)
	}

	s.Start()

	<-ctx.Done()

	return nil
}
