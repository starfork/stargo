package stargo

import (
	"github.com/starfork/stargo/client"
	"github.com/starfork/stargo/queue/store"
	"github.com/starfork/stargo/server"
)

// App App
type App struct {
	server *server.Server
	client *client.Client
	store  *store.Store
}

func New(opt ...Option) *App {

	return &App{}
} //
