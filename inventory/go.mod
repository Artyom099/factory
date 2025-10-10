module github.com/Artyom099/factory/inventory

go 1.24.0

replace github.com/Artyom099/factory/shared => ../shared

require (
	github.com/Artyom099/factory/shared v0.0.0-00010101000000-000000000000
	google.golang.org/grpc v1.76.0
)

require (
	golang.org/x/net v0.44.0 // indirect
	golang.org/x/sys v0.36.0 // indirect
	golang.org/x/text v0.29.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250804133106-a7a43d27e69b // indirect
	google.golang.org/protobuf v1.36.10 // indirect
)
