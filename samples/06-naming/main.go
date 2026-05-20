package main

import (
	"github.com/starfork/stargo"
	"github.com/starfork/stargo/config"
)

func main() {
	conf, _ := config.LoadConfig()
	app := stargo.New("naming-demo", conf)

	// When YAML config has a "registry" section with scheme "etcd",
	// the service is automatically registered with etcd on app.Run()
	// and deregistered on app.Stop().

	r := app.Registry()
	if r == nil {
		app.LogInfof("registry not configured, skipping demo")
		return
	}
	app.LogInfof("registry scheme: %s", r.Scheme())

	rr := app.Resolver()
	if rr == nil {
		app.LogInfof("resolver not configured, skipping demo")
		return
	}
	app.LogInfof("resolver scheme: %s", rr.Scheme())

	// The service is discoverable by other services via the etcd resolver.
	app.LogInfof("naming-demo ready for service discovery")
}
