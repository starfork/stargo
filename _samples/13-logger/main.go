package main

import (
	"github.com/starfork/stargo"
	"github.com/starfork/stargo/config"
	// slog 是 Go 标准库实现, 无需额外 import（已在根模块自动注册）
	// slog is Go stdlib, no extra import needed (auto-registered in root module)
	_ "github.com/starfork/stargo/logger/zap" // zap 高性能日志驱动 / zap high-perf log driver
	pb "github.com/starfork/stargo/samples/proto/sample"
)

func main() {
	// 切换日志驱动: 修改 config.yaml 中的 log.driver 字段
	// Switch log driver: change log.driver in config.yaml
	//
	//   ""     → 默认 console 日志 / default console logger
	//   "slog" → Go 标准库结构日志 / Go stdlib structured logging
	//   "zap"  → Uber zap 高性能日志 / Uber zap high-performance logger
	//
	//   slog 不需要额外 import（已在根模块中），zap 需要：
	//   slog needs no import (in root module), zap requires:
	//   import _ "github.com/starfork/stargo/logger/zap"
	//
	//   $ go run .                              # default console
	//   $ go run . config.slog.yaml            # slog
	//   $ go run . config.zap.yaml             # zap

	conf, _ := config.LoadConfig()
	app := stargo.New("logger-demo", conf)
	h := NewHandler(app.Logger())

	pb.RegisterSampleServiceServer(app.RpcServer(), h)
	app.Run(&pb.SampleService_ServiceDesc, h)
}
