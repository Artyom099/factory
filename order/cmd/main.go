package main

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	orderApiV1 "github.com/Artyom099/factory/order/internal/api/order/v1"
	grpcInventoryV1 "github.com/Artyom099/factory/order/internal/client/grpc/inventory/v1"
	grpcPaymentV1 "github.com/Artyom099/factory/order/internal/client/grpc/payment/v1"
	orderRepository "github.com/Artyom099/factory/order/internal/repository/order"
	orderService "github.com/Artyom099/factory/order/internal/service/order"
	orderV1 "github.com/Artyom099/factory/shared/pkg/openapi/order/v1"
	inventoryV1 "github.com/Artyom099/factory/shared/pkg/proto/inventory/v1"
	paymentV1 "github.com/Artyom099/factory/shared/pkg/proto/payment/v1"
)

const (
	httpPort               = "8080"
	inventoryServerAddress = "localhost:50051"
	paymentServerAddress   = "localhost:50052"
	readHeaderTimeout      = 5 * time.Second
	shutdownTimeout        = 10 * time.Second
)

func main() {
	inventoryConn, err := grpc.NewClient(
		inventoryServerAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("failed to connect inventory: %v", err)
	}

	inventoryClient := inventoryV1.NewInventoryServiceClient(inventoryConn)

	paymentConn, err := grpc.NewClient(
		paymentServerAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("failed to connect payment: %v", err)
	}

	paymentClient := paymentV1.NewPaymentServiceClient(paymentConn)

	grpcInventoryClient := grpcInventoryV1.NewClient(inventoryClient)
	grpcPaymentClient := grpcPaymentV1.NewClient(paymentClient)

	repo := orderRepository.NewRepository()
	service := orderService.NewService(repo, grpcInventoryClient, grpcPaymentClient)
	api := orderApiV1.NewAPI(service)

	orderServer, err := orderV1.NewServer(api)
	if err != nil {
		log.Fatalf("ошибка создания сервера OpenAPI: %v", err)
	}

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(10 * time.Second))
	r.Use(render.SetContentType(render.ContentTypeJSON))

	r.Mount("/", orderServer)

	server := &http.Server{
		Addr:              net.JoinHostPort("localhost", httpPort),
		Handler:           r,
		ReadHeaderTimeout: readHeaderTimeout, // Защита от Slowloris атак - тип DDoS-атаки
	}

	go func() {
		log.Printf("🚀 HTTP-сервер запущен на порту %s\n", httpPort)
		err := server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("❌ Ошибка запуска сервера: %v\n", err)
		}
	}()

	// Graceful shutdown
	quitCh := make(chan os.Signal, 1)
	signal.Notify(quitCh, syscall.SIGINT, syscall.SIGTERM)
	<-quitCh

	log.Println("🛑 Завершение работы сервера...")

	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	err = server.Shutdown(ctx)
	if err != nil {
		log.Printf("❌ Ошибка при остановке сервера: %v\n", err)
	}

	log.Println("✅ Сервер остановлен")
}
