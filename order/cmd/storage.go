package main

import (
	"sync"

	"github.com/Artyom099/factory/order/intermal/model"
	orderV1 "github.com/Artyom099/factory/shared/pkg/openapi/order/v1"
)

type OrderStorage struct {
	mu     sync.RWMutex
	orders map[string]*orderV1.Order
}

func NewOrderStorage() *OrderStorage {
	return &OrderStorage{
		orders: make(map[string]*orderV1.Order),
	}
}

func (s *OrderStorage) GetOrder(uuid string) (*orderV1.Order, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	order, ok := s.orders[uuid]
	if !ok {
		return nil, model.ErrOrderNotFound
	}

	return order, nil
}

func (s *OrderStorage) CreateOrder(order *orderV1.Order) (*orderV1.OrderCreateResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.orders[order.OrderUUID] = order

	return &orderV1.OrderCreateResponse{
		OrderUUID:  order.OrderUUID,
		TotalPrice: order.TotalPrice,
	}, nil
}

func (s *OrderStorage) CancelOrder(uuid string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	order, ok := s.orders[uuid]
	if !ok {
		return
	}

	order.Status = orderV1.OrderStatusCANCELLED
}

func (s *OrderStorage) UpdateOrder(uuid string, dto *orderV1.Order) {
	s.mu.Lock()
	defer s.mu.Unlock()

	order, ok := s.orders[uuid]
	if !ok {
		return
	}

	order.Status = dto.Status
	order.TransactionUUID = dto.TransactionUUID
	order.PaymentMethod = dto.PaymentMethod
}
