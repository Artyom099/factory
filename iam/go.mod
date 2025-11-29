module github.com/Artyom099/factory/iam

go 1.24.0

replace github.com/Artyom099/factory/shared => ../shared

replace github.com/Artyom099/factory/platform => ../platform

require (
	github.com/Artyom099/factory/platform v0.0.0-00010101000000-000000000000
	github.com/Artyom099/factory/shared v0.0.0-00010101000000-000000000000
	github.com/caarlos0/env/v11 v11.3.1
	github.com/gomodule/redigo v1.9.3
	github.com/google/uuid v1.6.0
	github.com/jackc/pgx/v5 v5.7.6
	github.com/joho/godotenv v1.5.1
	github.com/pkg/errors v0.9.1
	go.uber.org/zap v1.27.1
	golang.org/x/crypto v0.43.0
	google.golang.org/grpc v1.76.0
)

require (
	github.com/envoyproxy/protoc-gen-validate v1.2.1 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/puddle/v2 v2.2.2 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/net v0.46.0 // indirect
	golang.org/x/sync v0.17.0 // indirect
	golang.org/x/sys v0.37.0 // indirect
	golang.org/x/text v0.30.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250825161204-c5933d9347a5 // indirect
	google.golang.org/protobuf v1.36.10 // indirect
)
