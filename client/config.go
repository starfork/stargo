package client

import "github.com/starfork/stargo/naming"

type Config struct {
	Endpoints []string
	Target    string
	App       string
	Token     string
	Timeout   int64

	Registry naming.Config
}
