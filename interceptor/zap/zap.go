package zap

import (
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
	"gopkg.in/natefinch/lumberjack.v2"
)

//Unary interceptor
func Unary() grpc.UnaryServerInterceptor {
	return grpc_zap.UnaryServerInterceptor(Interceptor())
}

//Interceptor 返回zap.logger实例(把日志写到文件中)
func Interceptor() *zap.Logger {
	w := zapcore.AddSync(&lumberjack.Logger{
		Filename:  "log/debug.log",
		MaxSize:   1024, //MB
		LocalTime: true,
	})

	config := zap.NewProductionEncoderConfig()
	config.EncodeTime = zapcore.ISO8601TimeEncoder
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(config),
		w,
		zap.NewAtomicLevel(),
	)

	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	grpc_zap.ReplaceGrpcLogger(logger)
	return logger
}

// // Interceptor 返回zap.logger实例(把日志输出到控制台)
// func Interceptor() *zap.Logger {
// 	logger, err := zap.NewDevelopment()
// 	if err != nil {
// 		log.Fatalf("failed to initialize zap logger: %v", err)
// 	}
// 	grpc_zap.ReplaceGrpcLogger(logger)
// 	return logger
// }
