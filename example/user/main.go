package main

import (
	"user/internal/server"
	pb "user/pkg/pb"

	"github.com/starfork/stargo"
)

func main() {
	//c, _ := config.LoadConfig()
	//conf := c.GetServerConfig()
	//初始化数据库、redis,日志等

	app := stargo.New()

	pb.RegisterUserHandlerServer(app.Server(), server.New())
	app.Run()

}
