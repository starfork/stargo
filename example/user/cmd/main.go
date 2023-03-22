package main

import (
	"user/internal/server"
	pb "user/pkg/pb"

	"github.com/starfork/stargo"
	"github.com/starfork/stargo/config"
)

func main() {
	c, _ := config.LoadConfig()
	//conf := c.GetServerConfig()

	app := stargo.New(
		stargo.Config(c),
	)

	pb.RegisterUserHandlerServer(app.Server(), server.New())
	app.Run()

}
