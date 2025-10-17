package model

type PaymentMethod int

const (
	UNSPECIFIED    PaymentMethod = 0
	CARD           PaymentMethod = 1
	SBP            PaymentMethod = 2
	CREDIT_CARD    PaymentMethod = 3
	INVESTOR_MONEY PaymentMethod = 4
)

type PayOrderServiceRequestDto struct {
	OrderUuid     string
	UserUuid      string
	PaymentMethod PaymentMethod
}
