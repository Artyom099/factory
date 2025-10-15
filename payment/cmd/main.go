package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	paymentApiV1 "github.com/Artyom099/factory/payment/internal/api/payment/v1"
	paymentRepository "github.com/Artyom099/factory/payment/internal/repository/payment"
	paymentService "github.com/Artyom099/factory/payment/internal/service/payment"
	paymentV1 "github.com/Artyom099/factory/shared/pkg/proto/payment/v1"
)

const grpcPort = 50052

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

	s := grpc.NewServer()

	repo := paymentRepository.NewRepository()
	service := paymentService.NewService(repo)
	api := paymentApiV1.NewAPI(service)

	paymentV1.RegisterPaymentServiceServer(s, api)

	reflection.Register(s)

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
