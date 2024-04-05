package main

import (
	"context"
	"os"
	"os/signal"
	"path"
	"runtime"
	"avito/pkg/config"
	"strings"
	"syscall"

	"avito/internal/app"
	pc "avito/pkg/context"

	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

func init() {

	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
		ForceColors:     true,
		DisableColors:   false,
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			_, fileName := path.Split(f.File)
			arr := strings.Split(f.Function, ".")

			dir := " " + arr[0] + "/"
			funcName := strings.Join(arr[1:], ".")

			return funcName, dir + fileName
		},
	})

	logrus.SetLevel(logrus.TraceLevel)
	logrus.SetOutput(os.Stdout)
	logrus.SetReportCaller(true)

}

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @description	Токен авторизации в формате "Bearer jwt-token"
func main() {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logrus.Info("Инициализация конфигурации")
	ctx = pc.AddConfig(ctx, config.New())

	logrus.Info("Инициализация App")
	a, err := app.New(ctx)
	if err != nil {
		logrus.Fatal(err)
	}

	// Стартуем фронт, если нужно
	//	go initDistributionStatic(ctx)

	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {

		signalChannel := make(chan os.Signal, 1)
		signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)

		select {
		case <-signalChannel:
			cancel()

		case <-ctx.Done():
			return ctx.Err()
		}

		a.Stop(ctx)

		return nil

	})

	g.Go(func() error {

		logrus.Info("Запуск App")

		return a.Run(ctx)

	})

	if err := g.Wait(); err != nil {
		logrus.Errorf("Приложение упало с ошибкой: %v", err)
	}

	logrus.Warn("app stopped")

}

// func initDistributionStatic(ctx context.Context) {

// 	conf := config.GetConfig(ctx)

// 	if conf.HTTP.StartFront {
// 		application := fiber.New(fiber.Config{DisableStartupMessage: true})
// 		application.Static("/", conf.HTTP.DistFolder)
// 		application.Get("/*", func(ctx *fiber.Ctx) error {
// 			return ctx.SendFile(fmt.Sprintf("%s/index.html", conf.HTTP.DistFolder))
// 		})
// 		// TODO Обязательно реализовать Shutdown
// 		_ = application.Listen(fmt.Sprintf("0.0.0.0:%d", conf.HTTP.DistPort))
// 	}

// }
