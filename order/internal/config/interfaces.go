package config

type LoggerConfig interface {
	Level() string
	AsJson() bool
}

type OrderGRPCConfig interface {
	Address() string
}

type PostgresConfig interface {
	URI() string
	DatabaseName() string
}

type PaymentClientConfig interface {
	Address() string
}

type InventoryClientConfig interface {
	Address() string
}
