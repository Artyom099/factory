package model

import "time"

type RepoOrder struct {
	OrderUUID       string             `db:"order_uuid"`
	UserUUID        string             `db:"user_uuid"`
	PartUuids       []string           `db:"part_uuids"`
	TotalPrice      float32            `db:"total_price"`
	TransactionUUID *string            `db:"transaction_uuid"`
	PaymentMethod   OrderPaymentMethod `db:"payment_method"`
	Status          OrderStatus        `db:"status"`
	CreatedAt       time.Time          `db:"created_at"`
	UpdatedAt       *time.Time         `db:"updated_at"`
}

type OrderPaymentMethod string

const (
	OrderPaymentMethodUNSPECIFIED   OrderPaymentMethod = "UNSPECIFIED"
	OrderPaymentMethodCARD          OrderPaymentMethod = "CARD"
	OrderPaymentMethodSBP           OrderPaymentMethod = "SBP"
	OrderPaymentMethodCREDITCARD    OrderPaymentMethod = "CREDIT_CARD"
	OrderPaymentMethodINVESTORMONEY OrderPaymentMethod = "INVESTOR_MONEY"
)

type OrderStatus string

const (
	OrderStatusPAID           OrderStatus = "PAID"
	OrderStatusPENDINGPAYMENT OrderStatus = "PENDING_PAYMENT"
	OrderStatusCANCELLED      OrderStatus = "CANCELLED"
)
