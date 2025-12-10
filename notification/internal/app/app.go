package app

import (
	"context"

	"github.com/go-faster/errors"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"

	"github.com/Artyom099/factory/notification/internal/config"
	"github.com/Artyom099/factory/platform/pkg/closer"
	"github.com/Artyom099/factory/platform/pkg/grpc/health"
	"github.com/Artyom099/factory/platform/pkg/logger"
	"github.com/Artyom099/factory/platform/pkg/tracing"
)

type App struct {
	diContainer *diContainer
	grpcServer  *grpc.Server
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
	// Канал для ошибок от компонентов
	errCh := make(chan error, 2)

	// Контекст для остановки всех горутин
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// OrderPaid Консьюмер
	go func() {
		if err := a.runOrderPaidConsumer(ctx); err != nil {
			errCh <- errors.Errorf("consumer crashed: %v", err)
		}
	}()

	// Order Assembled Консьюмер
	go func() {
		if err := a.runOrderAssembledConsumer(ctx); err != nil {
			errCh <- errors.Errorf("consumer crashed: %v", err)
		}
	}()

	// GRPC сервер - только для health check
	// go func() {
	// 	if err := a.runGRPCServer(ctx); err != nil {
	// 		errCh <- errors.Errorf("grpc server crashed: %v", err)
	// 	}
	// }()

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
		a.initGRPCServer,
		a.initTelegramBot,
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

func (a *App) initGRPCServer(ctx context.Context) error {
	a.grpcServer = grpc.NewServer(grpc.Creds(insecure.NewCredentials()))

	reflection.Register(a.grpcServer)

	// Только healthcheck, никаких API
	health.RegisterService(a.grpcServer)

	closer.AddNamed("gRPC server", func(ctx context.Context) error {
		a.grpcServer.GracefulStop()
		return nil
	})

	return nil
}

// func (a *App) runGRPCServer(ctx context.Context) error {
// 	lis, err := net.Listen("tcp", config.AppConfig().NotificationGRPC.Address())
// 	if err != nil {
// 		return err
// 	}

// 	logger.Info(ctx, "Notification Healthcheck gRPC running on "+config.AppConfig().NotificationGRPC.Address())

// 	return a.grpcServer.Serve(lis)
// }

func (a *App) runOrderPaidConsumer(ctx context.Context) error {
	logger.Info(ctx, "🚀 OrderPaid Kafka consumer running")

	err := a.diContainer.OrderPaidConsumerService(ctx).RunConsumer(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (a *App) runOrderAssembledConsumer(ctx context.Context) error {
	logger.Info(ctx, "🚀 OrderAssembled Kafka consumer running")

	err := a.diContainer.OrderAssembledConsumerService(ctx).RunConsumer(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (a *App) initTelegramBot(ctx context.Context) error {
	// Получаем бота из DI контейнера
	telegramBot := a.diContainer.TelegramBot(ctx)

	// Регистрируем обработчик для активации бота
	telegramBot.RegisterHandler(bot.HandlerTypeMessageText, "/start", bot.MatchTypeExact, func(ctx context.Context, b *bot.Bot, update *models.Update) {
		logger.Info(ctx, "chat id", zap.Int64("chat_id", update.Message.Chat.ID))

		_, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Notification Factory активирован! Теперь вы будете получать уведомления о новых заказах.",
		})
		if err != nil {
			logger.Error(ctx, "Failed to send activation message", zap.Error(err))
		}
	})

	// Запускаем бота в фоне
	go func() {
		logger.Info(ctx, "🤖 Telegram bot started...")
		telegramBot.Start(ctx)
	}()

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
