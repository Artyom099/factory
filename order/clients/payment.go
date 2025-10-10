package clients

import (
	"log"

	paymentV1 "github.com/Artyom099/factory/shared/pkg/proto/payment/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const paymentServerAddress = "localhost:50052"

func CreatePaymentClient() (paymentV1.PaymentV1ServiceClient, error) {
	// ctx := context.Background()

	connection, err := grpc.NewClient(
		paymentServerAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Printf("failed to connect: %v\n", err)
		return nil, err
	}
	defer func() {
		if cerr := connection.Close(); cerr != nil {
			log.Printf("failed to close connect: %v", cerr)
		}
	}()

	// Создаем gRPC клиент
	return paymentV1.NewPaymentV1ServiceClient(connection), nil
}
