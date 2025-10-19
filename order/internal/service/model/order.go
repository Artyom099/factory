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
}

// type Category int32

// const (
// 	UNSPECIFIED Category = 0
// 	ENGINE      Category = 1
// 	FUEL        Category = 2
// 	PORTHOLE    Category = 3
// 	WING        Category = 4
// )

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

type Category int32

const (
	Category_CATEGORY_UNSPECIFIED Category = 0
	Category_CATEGORY_ENGINE      Category = 1
	Category_CATEGORY_FUEL        Category = 2
	Category_CATEGORY_PORTHOLE    Category = 3
	Category_CATEGORY_WING        Category = 4
)

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
