package main

import (
	"log"
	"os"

	"github.com/starfork/stargo"
	"github.com/starfork/stargo/config"
)

func main() {
	// Log level can be set via STARGO_LOG_LEVEL env var
	// Supported: trace, debug, info, warn, error, fatal
	os.Setenv("STARGO_LOG_LEVEL", "debug")

	conf, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	app := stargo.New("basic-service", conf)

	// Your gRPC service would be registered here:
	// pb.RegisterYourServiceServer(app.RpcServer(), &handler{})

	// app.Run(desc, impl) should be called with your service descriptor and implementation.
	// For this demo we skip the actual gRPC service registration.
	app.LogInfof("basic-service initialized, ready to register gRPC services")
	_ = app
}
