package v1

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/Artyom099/factory/payment/internal/api/converter"
	paymentV1 "github.com/Artyom099/factory/shared/pkg/proto/payment/v1"
)

func (a *api) PayOrder(ctx context.Context, req *paymentV1.PayOrderRequest) (*paymentV1.PayOrderResponse, error) {
	if err := req.Validate(); err != nil {
		return &paymentV1.PayOrderResponse{}, status.Errorf(codes.InvalidArgument, "validation error: %v", err)
	}

	switch req.GetPaymentMethod() {
	case paymentV1.PaymentMethod_PAYMENT_METHOD_CARD,
		paymentV1.PaymentMethod_PAYMENT_METHOD_SBP,
		paymentV1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD,
		paymentV1.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY:
	default:
		return &paymentV1.PayOrderResponse{}, status.Error(codes.InvalidArgument, "unsupported payment_method")
	}

	transactionUuid, err := a.paymentService.PayOrder(ctx, converter.PayOrderApiRequestToPayOrderServiceRequest(req))
	if err != nil {
		return nil, err
	}

	return &paymentV1.PayOrderResponse{TransactionUuid: transactionUuid}, nil
}
