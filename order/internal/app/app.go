package app

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/Artyom099/factory/order/internal/config"
	"github.com/Artyom099/factory/order/internal/metrics"
	"github.com/Artyom099/factory/platform/pkg/closer"
	"github.com/Artyom099/factory/platform/pkg/logger"
	"github.com/Artyom099/factory/platform/pkg/migrator/pg"
	"github.com/Artyom099/factory/platform/pkg/tracing"
	orderV1 "github.com/Artyom099/factory/shared/pkg/openapi/order/v1"
)

type App struct {
	diContainer *diContainer
	httpServer  *http.Server
	listener    net.Listener
}

func New(ctx context.Context) (*App, error) {
	a := &App{}

	err := a.initDeps(ctx)
	if err != nil {
		return nil, err
	}

	return a, nil
}

func (a *App) Run(ctx context.Context) error {
	errCh := make(chan error, 2)

	// Контекст для остановки всех горутин
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// HTTP сервер
	go func() {
		if err := a.runHTTPServer(ctx); err != nil {
			errCh <- errors.Errorf("http server crashed: %v", err)
		}
	}()

	// Order Assembled Консьюмер
	go func() {
		if err := a.runOrderAssembledConsumer(ctx); err != nil {
			errCh <- errors.Errorf("consumer crashed: %v", err)
		}
	}()

	// Ожидание либо ошибки, либо завершения контекста (например, сигнал SIGINT/SIGTERM)
	select {
	case <-ctx.Done():
		logger.Info(ctx, "Shutdown signal received")
	case err := <-errCh:
		logger.Error(ctx, "Component crashed, shutting down", zap.Error(err))
		// Триггерим cancel, чтобы остановить второй компонент
		cancel()
		// Дождись завершения всех задач (если есть graceful shutdown внутри)
		<-ctx.Done()
		return err
	}

	return nil
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initDI,
		a.initLogger,
		a.initCloser,
		// a.initTracing, //todo - надо подключить как в примере
		a.initMetrics,
		a.initListener,
		a.initHTTPServer,
		a.initMigrator,
		a.initTracing,
	}

	for _, f := range inits {
		err := f(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *App) initDI(_ context.Context) error {
	a.diContainer = NewDiContainer()
	return nil
}

func (a *App) initLogger(ctx context.Context) error {
	return logger.Init(
		ctx,
		config.AppConfig().Logger.Level(),
		config.AppConfig().Logger.AsJson(),
		config.AppConfig().Logger.EnableOTLP(),
	)
}

func (a *App) initCloser(_ context.Context) error {
	closer.SetLogger(logger.Logger())
	return nil
}

func (a *App) initListener(_ context.Context) error {
	listener, err := net.Listen("tcp", config.AppConfig().OrderHTTP.Address())
	if err != nil {
		return err
	}
	closer.AddNamed("TCP listener", func(ctx context.Context) error {
		lerr := listener.Close()
		if lerr != nil && !errors.Is(lerr, net.ErrClosed) {
			return lerr
		}

		return nil
	})

	a.listener = listener

	return nil
}

func (a *App) initHTTPServer(ctx context.Context) error {
	orderServer, err := orderV1.NewServer(a.diContainer.OpderV1API(ctx))
	if err != nil {
		return err
	}

	mux := chi.NewRouter()
	mux.Use(a.diContainer.AuthMiddleware(ctx).Handle)
	mux.Use(tracing.HTTPHandlerMiddleware("order-service"))
	mux.Use(render.SetContentType(render.ContentTypeJSON))
	mux.Mount("/", orderServer)

	readHeaderTimeout := 5 * time.Second

	a.httpServer = &http.Server{
		Addr: config.AppConfig().OrderHTTP.Address(),
		// todo - в конфиг
		ReadHeaderTimeout: readHeaderTimeout, // Защита от Slowloris атак - тип DDoS-атаки
		Handler:           mux,
	}

	closer.AddNamed("HTTP server", func(ctx context.Context) error {
		return a.httpServer.Shutdown(ctx)
	})

	return nil
}

func (a *App) initMigrator(ctx context.Context) error {
	migrationsDir := os.Getenv("MIGRATIONS_DIR")
	if migrationsDir == "" {
		return errors.New("MIGRATIONS_DIR env is not set")
	}

	if _, err := os.Stat(migrationsDir); err != nil {
		return fmt.Errorf("migrations directory %s: %w", migrationsDir, err)
	}

	pool := a.diContainer.PostgresHandle(ctx)
	connConfig := pool.Config().ConnConfig.Copy()
	sqlDB := stdlib.OpenDB(*connConfig)

	migratorRunner := pg.NewMigrator(sqlDB, migrationsDir)
	if err := migratorRunner.Up(); err != nil {
		return fmt.Errorf("apply migrations: %w", err)
	}

	logger.Info(ctx, fmt.Sprintf("✅ Миграции успешно применены из %s", migrationsDir))

	return nil
}

func (a *App) runHTTPServer(ctx context.Context) error {
	logger.Info(ctx, fmt.Sprintf("🚀 HTTP OrderService server listening on %s", config.AppConfig().OrderHTTP.Address()))

	err := a.httpServer.Serve(a.listener)
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}

func (a *App) runOrderAssembledConsumer(ctx context.Context) error {
	logger.Info(ctx, "🚀 OrderAssembled Kafka consumer running")

	err := a.diContainer.OrderConsumerService(ctx).RunConsumer(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (a *App) initTracing(ctx context.Context) error {
	err := tracing.InitTracer(ctx, config.AppConfig().Tracing)
	if err != nil {
		return err
	}

	closer.AddNamed("tracer", tracing.ShutdownTracer)

	return nil
}

func (a *App) initMetrics(ctx context.Context) error {
	err := metrics.InitMetrics(ctx, config.AppConfig().MetricServer)
	if err != nil {
		return err
	}
	closer.AddNamed("metrics", metrics.ShutdownMetrics)
	return nil
}
