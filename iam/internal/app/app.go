package app

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"

	"github.com/jackc/pgx/v5/stdlib"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"

	"github.com/Artyom099/factory/iam/internal/config"
	"github.com/Artyom099/factory/platform/pkg/closer"
	"github.com/Artyom099/factory/platform/pkg/grpc/health"
	"github.com/Artyom099/factory/platform/pkg/logger"
	"github.com/Artyom099/factory/platform/pkg/migrator/pg"
	"github.com/Artyom099/factory/platform/pkg/tracing"
	authV1 "github.com/Artyom099/factory/shared/pkg/proto/auth/v1"
	userV1 "github.com/Artyom099/factory/shared/pkg/proto/user/v1"
)

type App struct {
	diContainer *diContainer
	grpcServer  *grpc.Server
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
	return a.runGRPCServer(ctx)
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initDI,
		a.initLogger,
		a.initCloser,
		a.initListener,
		a.initGRPCServer,
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

func (a *App) initLogger(_ context.Context) error {
	return logger.Init(
		config.AppConfig().Logger.Level(),
		config.AppConfig().Logger.AsJson(),
	)
}

func (a *App) initCloser(_ context.Context) error {
	closer.SetLogger(logger.Logger())
	return nil
}

func (a *App) initListener(_ context.Context) error {
	listener, err := net.Listen("tcp", config.AppConfig().IamGRPC.Address())
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

func (a *App) initGRPCServer(ctx context.Context) error {
	a.grpcServer = grpc.NewServer(grpc.Creds(insecure.NewCredentials()))
	closer.AddNamed("gRPC server", func(ctx context.Context) error {
		a.grpcServer.GracefulStop()
		return nil
	})

	reflection.Register(a.grpcServer)

	// Регистрируем health service для проверки работоспособности
	health.RegisterService(a.grpcServer)

	authV1.RegisterAuthServiceServer(a.grpcServer, a.diContainer.AuthV1API(ctx))
	userV1.RegisterUserServiceServer(a.grpcServer, a.diContainer.UserV1API(ctx))

	return nil
}

func (a *App) runGRPCServer(ctx context.Context) error {
	logger.Info(ctx, fmt.Sprintf("🚀 gRPC IamService server listening on %s", config.AppConfig().IamGRPC.Address()))

	err := a.grpcServer.Serve(a.listener)
	if err != nil {
		return err
	}

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

func (a *App) initTracing(ctx context.Context) error {
	err := tracing.InitTracer(ctx, config.AppConfig().Tracing)
	if err != nil {
		return err
	}

	closer.AddNamed("tracer", tracing.ShutdownTracer)

	return nil
}
