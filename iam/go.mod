module github.com/Artyom099/factory/iam

go 1.24.0

replace github.com/Artyom099/factory/shared => ../shared

replace github.com/Artyom099/factory/platform => ../platform

require (
	github.com/Artyom099/factory/platform v0.0.0-00010101000000-000000000000
	go.uber.org/zap v1.27.1
)

require go.uber.org/multierr v1.11.0 // indirect
