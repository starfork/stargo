package server

import (
	"google.golang.org/grpc"
)

var (
	ENV_DEBUG      = "debug"
	ENV_PRODUCTION = "production"
)

// 公共配置模板
type Config struct {
	Env  string
	Name string
	Addr string
	ID   int
	//ServerName string //服务名称--4-11改。通过app启动设置
	//RpcServer *RpcServer

	UnaryInterceptor  []grpc.UnaryServerInterceptor
	StreamInterceptor []grpc.StreamServerInterceptor

	ServerOpts []grpc.ServerOption
}

type RpcServer struct {
	Entry string
	Name  string
	Host  string
	Port  string
	Auth  string //[keyfilepath]:[key]:
}

var DefaultConfig = &Config{
	//RpcServer:         &RpcServer{},
	UnaryInterceptor:  []grpc.UnaryServerInterceptor{},
	StreamInterceptor: []grpc.StreamServerInterceptor{},
	ServerOpts:        []grpc.ServerOption{},
}
