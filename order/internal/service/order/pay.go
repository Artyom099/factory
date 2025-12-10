package order

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"github.com/Artyom099/factory/order/internal/service/model"
	"github.com/Artyom099/factory/platform/pkg/tracing"
)

func (s *service) Pay(ctx context.Context, orderUUID string, paymentMethod model.OrderPaymentMethod) (string, error) {
	ctx, getOrderSpan := tracing.StartSpan(ctx, "order.get_order",
		trace.WithAttributes(
			attribute.String("order.uuid", orderUUID),
		),
	)

	order, err := s.orderRepository.Get(ctx, orderUUID)
	if err != nil {
		getOrderSpan.RecordError(err)
		getOrderSpan.End()
		if errors.Is(err, model.ErrOrderNotFound) {
			return "", model.ErrOrderNotFound
		}
		return "", err
	}

	getOrderSpan.End()

	if order.Status == model.OrderStatusPAID {
		return "", model.ErrOrderAlreadyPaid
	}

	if order.Status == model.OrderStatusCANCELLED {
		return "", model.ErrOrderCancelled
	}

	ctx, payOrderSpan := tracing.StartSpan(ctx, "order.call_payment_pay_order",
		trace.WithAttributes(
			attribute.String("order.uuid", orderUUID),
		),
	)

	transactionUUID, err := s.paymentClient.PayOrder(ctx, paymentMethod, orderUUID, order.UserUUID)
	if err != nil {
		payOrderSpan.RecordError(err)
		payOrderSpan.End()
		return "", model.ErrInPaymeentService
	}

	payOrderSpan.SetAttributes(
		attribute.String("order.transactionUUID", transactionUUID),
	)
	payOrderSpan.End()

	ctx, updateOrderSpan := tracing.StartSpan(ctx, "order.update_order",
		trace.WithAttributes(
			attribute.String("order.uuid", orderUUID),
		),
	)

	order.Status = model.OrderStatusPAID
	order.TransactionUUID = transactionUUID
	order.PaymentMethod = paymentMethod

	err = s.orderRepository.Update(ctx, order)
	if err != nil {
		updateOrderSpan.RecordError(err)
		updateOrderSpan.End()
		return "", err
	}

	updateOrderSpan.End()

	ctx, produceOrderPaidSpan := tracing.StartSpan(ctx, "order.produce_order_paid_event",
		trace.WithAttributes(
			attribute.String("order.uuid", orderUUID),
		),
	)
	defer produceOrderPaidSpan.End()

	// отправляем сообщение в кафку, что заказ оплачен
	err = s.orderProducerService.ProduceOrderPaid(ctx, model.OrderPaidOutEvent{
		EventUUID:       uuid.New().String(),
		OrderUUID:       orderUUID,
		UserUUID:        order.UserUUID,
		PaymentMethod:   string(paymentMethod),
		TransactionUUID: transactionUUID,
	})
	if err != nil {
		produceOrderPaidSpan.RecordError(err)
		return "", model.ErrSendOderPaidMessageToKafka
	}

	return transactionUUID, nil
}
