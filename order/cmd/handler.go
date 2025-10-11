package main

import (
	"context"
	"net/http"

	"github.com/google/uuid"

	"github.com/Artyom099/factory/order/clients"
	orderV1 "github.com/Artyom099/factory/shared/pkg/openapi/order/v1"
	inventoryV1 "github.com/Artyom099/factory/shared/pkg/proto/inventory/v1"
	paymentV1 "github.com/Artyom099/factory/shared/pkg/proto/payment/v1"
)

// OrderHandler реализует интерфейс orderV1.Handler для обработки запросов к API погоды
type OrderHandler struct {
	storage *OrderStorage
}

// NewOrderHandler создает новый обработчик запросов к API заказов
func NewOrderHandler(storage *OrderStorage) *OrderHandler {
	return &OrderHandler{
		storage: storage,
	}
}

func (h *OrderHandler) CancelOrder(ctx context.Context, params orderV1.CancelOrderParams) (orderV1.CancelOrderRes, error) {
	order := h.storage.GetOrder(params.OrderUUID.String())

	if order.Status == orderV1.OrderStatusPAID {
		return &orderV1.ConflictError{
			Code:    409,
			Message: "Order '" + params.OrderUUID.String() + "' already paid, cannot be cancelled",
		}, nil
	}

	if order.Status == orderV1.OrderStatusPENDINGPAYMENT {
		h.storage.UpdateOrderStatus(params.OrderUUID.String())
	}

	return nil, nil
}

func (h *OrderHandler) CreateOrder(ctx context.Context, params *orderV1.OrderCreateRequest) (orderV1.CreateOrderRes, error) {
	// Получает детали через `InventoryService.ListParts`.
	// Проверяет, что все детали существуют. Если хотя бы одной нет — возвращает ошибку.
	// Считает `total_price`.

	inventoryClient, err := clients.CreateInventoryClient()
	if err != nil {
		return &orderV1.InternalServerError{
			Code:    500,
			Message: "Failed to create inventory client",
		}, err
	}

	// Вызываем gRPC метод ListParts
	req := inventoryV1.ListPartsRequest{
		Filter: &inventoryV1.PartsFilter{
			Uuids: params.GetPartUuids(),
		},
	}
	parts, err := inventoryClient.ListParts(ctx, &req)
	if err != nil {
		return &orderV1.InternalServerError{
			Code:    500,
			Message: "Failed to call ListParts method",
		}, err
	}

	if len(parts.GetParts()) != len(params.GetPartUuids()) {
		return &orderV1.InternalServerError{
			Code:    500,
			Message: "Not all parts exist",
		}, nil
	}

	var totalPrice float32
	for _, part := range parts.GetParts() {
		totalPrice += float32(part.GetPrice())
	}

	orderUuid := uuid.New().String()

	dto := orderV1.Order{
		OrderUUID:       orderUuid,
		TotalPrice:      totalPrice,
		UserUUID:        "UserUUID",
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
	order := h.storage.GetOrder(params.OrderUUID.String())
	if order == nil {
		return &orderV1.NotFoundError{
			Code:    404,
			Message: "Data for specified order '" + params.OrderUUID.String() + "' not found",
		}, nil
	}

	return order, nil
}

func (h *OrderHandler) PayOrder(ctx context.Context, req *orderV1.OrderPayRequest, params orderV1.PayOrderParams) (orderV1.PayOrderRes, error) {
	// Находит заказ по `order_uuid`. Если не существует — возвращает 404 Not Found.
	order := h.storage.GetOrder(params.OrderUUID.String())
	if order == nil {
		return &orderV1.NotFoundError{
			Code:    404,
			Message: "Data for specified order '" + params.OrderUUID.String() + "' not found",
		}, nil
	}

	// Вызывает `PaymentService.PayOrder`, передаёт `user_uuid`, `order_uuid` и `payment_method`. Получает`transaction_uuid`.
	paymentClient, err := clients.CreatePaymentClient()
	if err != nil {
		return &orderV1.InternalServerError{
			Code:    500,
			Message: "Failed to create payment client",
		}, err
	}

	userUuid := uuid.New().String()
	dto := paymentV1.PayOrderRequest{
		UserUuid:      userUuid,
		OrderUuid:     params.OrderUUID.String(),
		PaymentMethod: convertPaymentMethod(req.GetPaymentMethod()),
	}

	res, err := paymentClient.PayOrder(ctx, &dto)
	if err != nil {
		return &orderV1.InternalServerError{
			Code:    500,
			Message: "Failed to pay order",
		}, err
	}

	// Обновляет заказ: статус → `PAID`, сохраняет `transaction_uuid`, `payment_method`
	updateDto := orderV1.Order{
		Status:          orderV1.OrderStatusPAID,
		TransactionUUID: res.TransactionUuid,
		PaymentMethod:   orderV1.OrderPaymentMethod(req.PaymentMethod),
	}
	h.storage.UpdateOrder(params.OrderUUID.String(), &updateDto)

	return &orderV1.OrderPayResponse{
		TransactionUUID: res.GetTransactionUuid(),
	}, nil
}

// NewError создает новую ошибку в формате GenericError
func (h *OrderHandler) NewError(_ context.Context, err error) *orderV1.GenericErrorStatusCode {
	return &orderV1.GenericErrorStatusCode{
		StatusCode: http.StatusInternalServerError,
		Response: orderV1.GenericError{
			Code:    orderV1.NewOptInt(http.StatusInternalServerError),
			Message: orderV1.NewOptString(err.Error()),
		},
	}
}

// преобразование метода оплаты из типа order сервиса в тип payment сервиса
func convertPaymentMethod(method orderV1.OrderPayRequestPaymentMethod) paymentV1.PaymentMethod {
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
