package v1

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/Artyom099/factory/payment/internal/api/converter"
	"github.com/Artyom099/factory/payment/internal/service/model"
	paymentV1 "github.com/Artyom099/factory/shared/pkg/proto/payment/v1"
)

func (a *api) PayOrder(ctx context.Context, req *paymentV1.PayOrderRequest) (*paymentV1.PayOrderResponse, error) {
	if err := req.Validate(); err != nil {
		return &paymentV1.PayOrderResponse{}, status.Errorf(codes.InvalidArgument, "validation error: %v", err)
	}

	transactionUuid, err := a.paymentService.PayOrder(ctx, converter.ToModelPayment(req))
	if err != nil {
		if errors.Is(err, model.ErrInvalidPaymentMethod) {
			return nil, status.Error(codes.InvalidArgument, "unsupported payment method")
		}
		return nil, status.Error(codes.Internal, "internal server error")
	}

	return &paymentV1.PayOrderResponse{TransactionUuid: transactionUuid}, nil
}
