package integration

import (
	"context"
	"os"
	"time"

	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.uber.org/zap"

	"github.com/Artyom099/factory/platform/pkg/logger"
	"github.com/Artyom099/factory/platform/pkg/testcontainers"
	"github.com/Artyom099/factory/platform/pkg/testcontainers/app"
	"github.com/Artyom099/factory/platform/pkg/testcontainers/network"
	"github.com/Artyom099/factory/platform/pkg/testcontainers/path"
	"github.com/Artyom099/factory/platform/pkg/testcontainers/postgres"
)

const (
	// Параметры для контейнеров
	orderAppName    = "order-app"
	orderDockerfile = "deploy/docker/order/Dockerfile"

	// Переменные окружения приложения
	httpPortKey = "HTTP_PORT"

	// Значения переменных окружения
	loggerLevelValue = "debug"
	startupTimeout   = 5 * time.Minute // было 3
)

// TestEnvironment — структура для хранения ресурсов тестового окружения
type TestEnvironment struct {
	Network  *network.Network
	Postgres *postgres.Container
	App      *app.Container
}

func setupTestEnvironment(ctx context.Context) *TestEnvironment {
	logger.Info(ctx, "🚀 Подготовка тестового окружения...")

	// Шаг 1: Создаём общую Docker-сеть
	generatedNetwork, err := network.NewNetwork(ctx, projectName)
	if err != nil {
		logger.Fatal(ctx, "не удалось создать общую сеть", zap.Error(err))
	}
	logger.Info(ctx, "✅ Сеть успешно создана")

	// Получаем переменные окружения для Postgres с проверкой на наличие
	postgresUsername := getEnvWithLogging(ctx, testcontainers.PostgresUsernameKey)
	postgresPassword := getEnvWithLogging(ctx, testcontainers.PostgresPasswordKey)
	postgresImageName := getEnvWithLogging(ctx, testcontainers.PostgresImageNameKey)
	postgresDatabase := getEnvWithLogging(ctx, testcontainers.PostgresDatabaseKey)

	// Получаем порт httpPort для waitStrategy
	httpPort := getEnvWithLogging(ctx, httpPortKey)

	// Шаг 2: Запускаем контейнер с Postgres
	generatedPostgres, err := postgres.NewContainer(ctx,
		postgres.WithNetworkName(generatedNetwork.Name()),
		postgres.WithContainerName(testcontainers.PostgresContainerName),
		postgres.WithImageName(postgresImageName),
		postgres.WithDatabase(postgresDatabase),
		postgres.WithAuth(postgresUsername, postgresPassword),
		postgres.WithLogger(logger.Logger()),
	)
	if err != nil {
		cleanupTestEnvironment(ctx, &TestEnvironment{Network: generatedNetwork})
		logger.Fatal(ctx, "не удалось запустить контейнер Postgres", zap.Error(err))
	}
	logger.Info(ctx, "✅ Контейнер Postgres успешно запущен")

	// Шаг 3: Запускаем контейнер с приложением
	projectRoot := path.GetProjectRoot()

	appEnv := map[string]string{
		// Переопределяем хост Postgres для подключения к контейнеру из testcontainers
		testcontainers.PostgresHostKey: generatedPostgres.Config().ContainerName,
	}

	// Создаем настраиваемую стратегию ожидания с увеличенным таймаутом
	waitStrategy := wait.ForListeningPort(nat.Port(httpPort + "/tcp")).
		WithStartupTimeout(startupTimeout)

	appContainer, err := app.NewContainer(ctx,
		app.WithName(orderAppName),
		app.WithPort(httpPort),
		app.WithDockerfile(projectRoot, orderDockerfile),
		app.WithNetwork(generatedNetwork.Name()),
		app.WithEnv(appEnv),
		app.WithLogOutput(os.Stdout),
		app.WithStartupWait(waitStrategy),
		app.WithLogger(logger.Logger()),
	)
	if err != nil {
		cleanupTestEnvironment(ctx, &TestEnvironment{Network: generatedNetwork, Postgres: generatedPostgres})
		logger.Fatal(ctx, "не удалось запустить контейнер приложения", zap.Error(err))
	}
	logger.Info(ctx, "✅ Контейнер приложения успешно запущен")

	logger.Info(ctx, "🎉 Тестовое окружение готово")
	return &TestEnvironment{
		Network:  generatedNetwork,
		Postgres: generatedPostgres,
		App:      appContainer,
	}
}

func getEnvWithLogging(ctx context.Context, key string) string {
	value := os.Getenv(key)
	if value == "" {
		logger.Warn(ctx, "Переменная окружения не установлена", zap.String("key", key))
	}

	return value
}
