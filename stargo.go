package stargo

import (
	"sync"
	"time"

	"github.com/starfork/stargo/broker"
	"github.com/starfork/stargo/logger"
	"github.com/starfork/stargo/naming"
	"github.com/starfork/stargo/server"
	"github.com/starfork/stargo/store"
)

// App App
type App struct {
	opts   *Options
	server *server.Server
	logger logger.Logger

	store    map[string]store.Store
	broker   broker.Broker
	registry naming.Registry

	once sync.Once
	//client *client.Client
}

func New(opt ...Option) *App {

	opts := DefaultOptions()
	for _, o := range opt {
		o(opts)
	}

	conf := opts.Config
	time.LoadLocation(opts.Timezone)
	conf.Timezome = opts.Timezone

	app := &App{
		opts: opts,
		//conf:   conf,
		store: make(map[string]store.Store),
	}
	app.once.Do(func() {
		app.logger = logger.DefaultLogger
		app.server = server.New()
	})

	// //注册reflection
	// if conf.server.Env != ENV_PRODUCTION {
	// 	app.logger.Debugf("env:" + conf.Env)
	// 	reflection.Register(app.server)
	// }

	// for k, v := range conf.Store {
	// 	app.Store(k, v)
	// }

	return app
}

// Run   server
func (s *App) Run() {

	s.server.Run()

}

// Stop server
func (s *App) Stop() {
	s.stopStargo()
	s.server.Stop()
	s.registry.UnRegister(s.server.Service())
}
func (s *App) stopStargo() {
	if s.registry != nil {
		s.logger.Fatalf("UnRegister: [%s]\n", s.opts.Name)
		s.registry.UnRegister(s.server.Service())
	}

	for _, st := range s.store {
		st.Close()
	}

	if s.broker != nil {
		s.broker.UnSubscribe()
	}
}

// Restart server
func (s *App) Restart() {
	s.stopStargo()
	s.server.Restart()
}
