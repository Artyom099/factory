package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	paymentV1 "github.com/Artyom099/factory/shared/pkg/proto/payment/v1"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const grpcPort = 50052

type paymentService struct {
	paymentV1.UnimplementedPaymentV1ServiceServer
}

func (s *paymentService) PayOrder(_ context.Context, _ *paymentV1.PayOrderRequest) (*paymentV1.PayOrderResponse, error) {
	transactionUuid := uuid.New().String()

	log.Printf("Оплата прошла успешно, transaction_uuid: %s", transactionUuid)

	return &paymentV1.PayOrderResponse{TransactionUuid: transactionUuid}, nil
}

func main() {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Printf("failed to listen: %v\n", err)
		return
	}
	defer func() {
		if cerr := listener.Close(); cerr != nil {
			log.Printf("failed to close listener: %v\n", cerr)
		}
	}()

	// Создаем gRPC сервер
	s := grpc.NewServer()

	service := &paymentService{}
	paymentV1.RegisterPaymentV1ServiceServer(s, service)

	reflection.Register(s)

	if err := s.Serve(listener); err != nil {
		log.Printf("failed to serve: %v\n", err)
	}

	go func() {
		log.Printf("Payment gRPC server listening on %d\n", grpcPort)
		err = s.Serve(listener)
		if err != nil {
			log.Printf("failed to serve: %v\n", err)
			return
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down payment gRPC server...")
	s.GracefulStop()
	log.Println("Server stopped")
}
