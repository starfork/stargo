module github.com/starfork/stargo/logger/zap

go 1.26.4

require (
	github.com/starfork/stargo v1.1.4
	go.uber.org/zap v1.28.0
)

require (
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	go.opentelemetry.io/otel v1.45.0 // indirect
	go.opentelemetry.io/otel/trace v1.45.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
)

replace github.com/starfork/stargo => ../../
