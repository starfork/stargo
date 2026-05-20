package main

import (
	"os"

	"github.com/starfork/stargo"
	"github.com/starfork/stargo/config"
)

func main() {
	// Log level can be set via STARGO_LOG_LEVEL env var.
	// Supported: trace, debug, info, warn, error, fatal
	os.Setenv("STARGO_LOG_LEVEL", "debug")

	conf, _ := config.LoadConfig()
	app := stargo.New("logger-demo", conf)

	// Use the app-level convenience methods
	app.LogInfof("server starting up")
	app.LogDebugf("debug detail: %s", "some value")
	app.LogWarnf("this is a warning")
	app.LogErrorf("something went wrong: %v", "connection refused")

	// Or access the logger directly
	l := app.Logger()
	l.Infof("direct logger call")

	// Output:
	// 2025/01/01 12:00:00 [INFO] server starting up
	// 2025/01/01 12:00:00 [DEBUG] debug detail: some value
	// 2025/01/01 12:00:00 [WARN] this is a warning
	// 2025/01/01 12:00:00 [ERROR] something went wrong: connection refused
	// 2025/01/01 12:00:00 [INFO] direct logger call
}
