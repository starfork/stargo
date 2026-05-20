package main

import (
	"github.com/starfork/stargo"
	"github.com/starfork/stargo/config"
)

func main() {
	conf, _ := config.LoadConfig()
	app := stargo.New("client-demo", conf)

	// app.Client() returns a client connected via the configured resolver.
	// With etcd naming configured, it discovers the target service.
	c := app.Client()
	if c == nil {
		app.LogInfof("no resolver configured, client unavailable")
		return
	}

	// Connect to a target service.
	// The target format is: scheme:///org/service-name
	conn, err := c.NewClient("target-service")
	if err != nil {
		app.LogFatalf("create client: %v", err)
		return
	}
	defer conn.Close()

	// Use conn to create a gRPC client stub:
	// pb.NewTargetServiceClient(conn)

	app.LogInfof("connected to target-service via: %s", conn.Target())
}
