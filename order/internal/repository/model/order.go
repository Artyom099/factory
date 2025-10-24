package model

type RepoOrder struct {
	OrderUUID       string
	UserUUID        string
	PartUuids       []string
	TotalPrice      float32
	TransactionUUID string
	PaymentMethod   OrderPaymentMethod
	Status          OrderStatus
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
