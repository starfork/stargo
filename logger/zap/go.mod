module github.com/starfork/stargo/logger/zap

go 1.26.4

require (
	github.com/starfork/stargo v0.0.0
	go.uber.org/zap v1.28.0
)

require go.uber.org/multierr v1.10.0 // indirect

replace github.com/starfork/stargo => ../../
