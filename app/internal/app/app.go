package app

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"avito/internal/repository"
	"avito/internal/service"
	"avito/internal/transport/thttp"
	pc "avito/pkg/context"
	"avito/pkg/errors"
	"avito/pkg/postgresql"

	"github.com/rs/cors"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

type App struct {
	router     *thttp.Router
	httpServer *http.Server
	pgClient   *postgresql.Postgres
	service    *service.Service
}

func New(ctx context.Context) (*App, error) {

	logrus.Info("Инициализация подключения к базе данных")
	pgClient, err := postgresql.New(ctx)
	if err != nil {
		return nil, errors.Wrap(err)
	}

	repository := repository.NewRepository(pgClient)
	service := service.NewService(repository)

	router, err := thttp.NewRouter(ctx, service)
	if err != nil {
		return nil, errors.Wrap(err)
	}

	return &App{
		router:   router,
		pgClient: pgClient,
		service:  service,
	}, nil

}

func (a *App) Run(ctx context.Context) error {

	g, ctx := errgroup.WithContext(ctx)

	logrus.Info("Запуск HTTP сервера")

	g.Go(func() error {
		return a.startHTTP(ctx)
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
	logrus.Info("pgClient.Close")

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
		Debug:              cfg.HTTP.CORS.Debug,
	})

	a.httpServer = &http.Server{
		Handler: c.Handler(a.router.Router),
		//WriteTimeout:      cfg.HTTP.WriteTimeout, // Для SSE его не установить
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
