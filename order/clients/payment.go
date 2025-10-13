package clients

import (
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	paymentV1 "github.com/Artyom099/factory/shared/pkg/proto/payment/v1"
)

const paymentServerAddress = "localhost:50052"

func CreatePaymentClient() (paymentV1.PaymentServiceClient, error) {
	connection, err := grpc.NewClient(
		paymentServerAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Printf("failed to connect: %v\n", err)
		return nil, err
	}

	return paymentV1.NewPaymentServiceClient(connection), nil
}
