package model

// EventUUID - Уникальный идентификатор события (для идемпотентности)

type OrderPaidInEvent struct {
	EventUUID       string
	OrderUUID       string
	UserUUID        string
	PaymentMethod   string
	TransactionUUID string
}

type OrderAssembledOutEvent struct {
	EventUUID    string
	OrderUUID    string
	UserUUID     string
	BuildTimeSec int64
}
