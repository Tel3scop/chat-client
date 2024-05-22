package app

import (
	"context"
	"sync"

	"github.com/Tel3scop/chat-client/internal/config"
	"github.com/Tel3scop/helpers/closer"
	"github.com/Tel3scop/helpers/logger"
)

// App структура приложения с сервис-провайдером и GRPC-сервером
type App struct {
	serviceProvider *serviceProvider
}

// NewApp вернуть новый экземпляр приложения с зависимостями
func NewApp(ctx context.Context) (*App, error) {
	a := &App{}

	err := a.initDeps(ctx)
	if err != nil {
		return nil, err
	}

	return a, nil
}

// Run запуск приложения
func (a *App) Run(ctx context.Context) error {
	defer func() {
		closer.CloseAll()
		closer.Wait()
	}()

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()
		a.serviceProvider.Cron().StartCron()
	}()

	a.serviceProvider.Console(ctx).Run()
	wg.Wait()
	return nil
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initConfig,
		a.initServiceProvider,
		a.initLogger,
	}

	for _, f := range inits {
		err := f(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *App) initConfig(_ context.Context) error {
	_, err := config.New()
	if err != nil {
		return err
	}

	return nil
}

func (a *App) initLogger(_ context.Context) error {
	loggerConfig := logger.Config{
		Filename:   a.serviceProvider.Config().Log.FileName,
		Level:      a.serviceProvider.Config().Log.Level,
		MaxSize:    a.serviceProvider.Config().Log.MaxSize,
		MaxBackups: a.serviceProvider.Config().Log.MaxBackups,
		MaxAge:     a.serviceProvider.Config().Log.MaxAge,
		Compress:   a.serviceProvider.Config().Log.Compress,
		StdOut:     false,
	}
	logger.InitByParams(loggerConfig)

	return nil
}

func (a *App) initServiceProvider(_ context.Context) error {
	a.serviceProvider = newServiceProvider()
	return nil
}
