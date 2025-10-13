package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	orderV1 "github.com/Artyom099/factory/shared/pkg/openapi/order/v1"
	inventoryV1 "github.com/Artyom099/factory/shared/pkg/proto/inventory/v1"
	paymentV1 "github.com/Artyom099/factory/shared/pkg/proto/payment/v1"
)

type OrderHandler struct {
	storage         *OrderStorage
	inventoryClient inventoryV1.InventoryServiceClient
	paymentClient   paymentV1.PaymentServiceClient
}

func NewOrderHandler(storage *OrderStorage, inventoryClient inventoryV1.InventoryServiceClient, paymentClient paymentV1.PaymentServiceClient) *OrderHandler {
	return &OrderHandler{
		storage:         storage,
		inventoryClient: inventoryClient,
		paymentClient:   paymentClient,
	}
}

func (h *OrderHandler) CancelOrder(ctx context.Context, params orderV1.CancelOrderParams) (orderV1.CancelOrderRes, error) {
	if _, err := uuid.Parse(params.OrderUUID.String()); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid order_uuid: %v", err)
	}

	order, err := h.storage.GetOrder(params.OrderUUID.String())
	if err != nil {
		return &orderV1.NotFoundError{
			Code:    404,
			Message: fmt.Sprintf("Order %s not found", params.OrderUUID.String()),
		}, nil
	}

	if order.Status == orderV1.OrderStatusPAID {
		return &orderV1.ConflictError{
			Code:    409,
			Message: fmt.Sprintf("Order %s already paid, cannot be cancelled", params.OrderUUID.String()),
		}, nil
	}

	if order.Status == orderV1.OrderStatusPENDINGPAYMENT {
		h.storage.CancelOrder(params.OrderUUID.String())
	}

	return nil, nil
}

func (h *OrderHandler) CreateOrder(ctx context.Context, req *orderV1.OrderCreateRequest) (orderV1.CreateOrderRes, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "validation error: %v", err)
	}

	listPartsReq := inventoryV1.ListPartsRequest{
		Filter: &inventoryV1.PartsFilter{
			Uuids: req.GetPartUuids(),
		},
	}
	parts, err := h.inventoryClient.ListParts(ctx, &listPartsReq)
	if err != nil {
		return &orderV1.InternalServerError{
			Code:    500,
			Message: "Failed to call ListParts method",
		}, err
	}

	if len(parts.GetParts()) != len(req.GetPartUuids()) {
		return &orderV1.InternalServerError{
			Code:    400,
			Message: "Not all parts exist",
		}, nil
	}

	var totalPrice float32
	for _, part := range parts.GetParts() {
		totalPrice += float32(part.GetPrice())
	}

	orderUuid := uuid.New().String()
	userUuid := req.GetUserUUID()

	dto := orderV1.Order{
		OrderUUID:       orderUuid,
		TotalPrice:      totalPrice,
		UserUUID:        userUuid,
		TransactionUUID: "",
		PaymentMethod:   orderV1.OrderPaymentMethodCARD,
		Status:          orderV1.OrderStatusPENDINGPAYMENT,
	}
	_, err = h.storage.CreateOrder(&dto)
	if err != nil {
		return &orderV1.BadRequestError{
			Code:    400,
			Message: "Failed to create order",
		}, nil
	}

	return &orderV1.OrderCreateResponse{
		OrderUUID:  orderUuid,
		TotalPrice: totalPrice,
	}, nil
}

func (h *OrderHandler) GetOrder(ctx context.Context, params orderV1.GetOrderParams) (orderV1.GetOrderRes, error) {
	if _, err := uuid.Parse(params.OrderUUID.String()); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid order_uuid: %v", err)
	}

	order, err := h.storage.GetOrder(params.OrderUUID.String())
	if err != nil {
		return &orderV1.NotFoundError{
			Code:    404,
			Message: fmt.Sprintf("Order %s not found", params.OrderUUID.String()),
		}, nil
	}

	return order, nil
}

func (h *OrderHandler) PayOrder(ctx context.Context, req *orderV1.OrderPayRequest, params orderV1.PayOrderParams) (orderV1.PayOrderRes, error) {
	if _, err := uuid.Parse(params.OrderUUID.String()); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid order_uuid: %v", err)
	}

	if err := req.Validate(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "validation error: %v", err)
	}

	order, err := h.storage.GetOrder(params.OrderUUID.String())
	if err != nil {
		return &orderV1.NotFoundError{
			Code:    404,
			Message: fmt.Sprintf("Data for specified order %s not found", params.OrderUUID.String()),
		}, nil
	}

	if order.Status == orderV1.OrderStatusPAID {
		return &orderV1.BadRequestError{
			Code:    400,
			Message: fmt.Sprintf("Order %s already paid, cannot be paid again", params.OrderUUID.String()),
		}, nil
	}

	if order.Status == orderV1.OrderStatusCANCELLED {
		return &orderV1.BadRequestError{
			Code:    400,
			Message: fmt.Sprintf("Order %s cancelled, cannot be paid", params.OrderUUID.String()),
		}, nil
	}

	userUuid := uuid.New().String()
	dto := paymentV1.PayOrderRequest{
		UserUuid:      userUuid,
		OrderUuid:     params.OrderUUID.String(),
		PaymentMethod: convertPaymentMethodFromOrderToPayment(req.GetPaymentMethod()),
	}

	res, err := h.paymentClient.PayOrder(ctx, &dto)
	if err != nil {
		return &orderV1.InternalServerError{
			Code:    500,
			Message: "Failed to pay order im payment service",
		}, err
	}

	updateDto := orderV1.Order{
		Status:          orderV1.OrderStatusPAID,
		TransactionUUID: res.TransactionUuid,
		PaymentMethod:   orderV1.OrderPaymentMethod(req.PaymentMethod),
	}
	h.storage.SetOrderPaid(params.OrderUUID.String(), &updateDto)

	return &orderV1.OrderPayResponse{
		TransactionUUID: res.GetTransactionUuid(),
	}, nil
}

func (h *OrderHandler) NewError(_ context.Context, err error) *orderV1.GenericErrorStatusCode {
	return &orderV1.GenericErrorStatusCode{
		StatusCode: http.StatusInternalServerError,
		Response: orderV1.GenericError{
			Code:    orderV1.NewOptInt(http.StatusInternalServerError),
			Message: orderV1.NewOptString(err.Error()),
		},
	}
}

func convertPaymentMethodFromOrderToPayment(method orderV1.OrderPayRequestPaymentMethod) paymentV1.PaymentMethod {
	switch method {
	case orderV1.OrderPayRequestPaymentMethodCARD:
		return paymentV1.PaymentMethod_PAYMENT_METHOD_CARD
	case orderV1.OrderPayRequestPaymentMethodSBP:
		return paymentV1.PaymentMethod_PAYMENT_METHOD_SBP
	case orderV1.OrderPayRequestPaymentMethodCREDITCARD:
		return paymentV1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD
	case orderV1.OrderPayRequestPaymentMethodINVESTORMONEY:
		return paymentV1.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY
	default:
		return paymentV1.PaymentMethod_PAYMENT_METHOD_UNSPECIFIED
	}
}
