package v1

import (
	def "github.com/Artyom099/factory/order/internal/client/grpc"
	paymentV1 "github.com/Artyom099/factory/shared/pkg/proto/payment/v1"
)

var _ def.IPaymentClient = (*client)(nil)

type client struct {
	generatedClient paymentV1.PaymentServiceClient
}

func NewClient(generatedClient paymentV1.PaymentServiceClient) *client {
	return &client{
		generatedClient: generatedClient,
	}
}
