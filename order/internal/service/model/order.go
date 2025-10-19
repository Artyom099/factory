package model

type OrderCreateServiceRequestDto struct {
	UserUUID  string
	PartUuids []string
}

type OrderCreateServiceResponseDto struct {
	OrderUUID  string
	TotalPrice float32
}

type OrderGetServiceRequestDto struct {
	OrderUUID string
}

type OrderGetServiceResponseDto struct {
	OrderUUID       string
	UserUUID        string
	PartUuids       []string
	TotalPrice      float32
	TransactionUUID string
	PaymentMethod   OrderPaymentMethod
	Status          OrderStatus
}

type OrderCancelServiceRequestDto struct {
	OrderUUID string
}
type OrderCancelServiceResponseDto struct{}

type OrderPayServiceRequestDto struct {
	OrderUUID     string
	PaymentMethod OrderPaymentMethod
}
type OrderPayServiceResponseDto struct {
	TransactionUUID string
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
