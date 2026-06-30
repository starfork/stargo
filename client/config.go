package client

import "github.com/starfork/stargo/naming"

type Config struct {
	Endpoints []string
	Target    string
	App       string
	Token     string

	Registry naming.Config
}
