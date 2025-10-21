package model

import "time"

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

type OrderCancelServiceResponseDto struct{}

type OrderPayServiceRequestDto struct {
	OrderUUID     string
	PaymentMethod OrderPaymentMethod
}

type ListPartsFilter struct {
	Uuids                 []string
	Names                 []string
	Categories            []Category
	ManufacturerCountries []string
	Tags                  []string
}

type ListPartsResponseDto struct {
	Parts []*Part
}

type Value struct {
	StringValue *string
	Int64Value  *int64
	DoubleValue *float64
	BoolValue   *bool
}

type Part struct {
	Uuid          string
	Name          string
	Description   string
	Price         float64
	StockQuantity int64
	Category      Category
	Dimensions    *Dimensions
	Manufacturer  *Manufacturer
	Tags          []string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Metadata      map[string]*Value
}

type Category int32

const (
	UNSPECIFIED Category = 0
	ENGINE      Category = 1
	FUEL        Category = 2
	PORTHOLE    Category = 3
	WING        Category = 4
)

type Dimensions struct {
	Length float64
	Width  float64
	Height float64
	Weight float64
}

type Manufacturer struct {
	Name    string
	Country string
	Website string
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
