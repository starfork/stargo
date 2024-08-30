package validator

import (
	"context"

	"google.golang.org/grpc"
)

func Unary() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, in *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (resp interface{}, err error) {

		//value := reflect.ValueOf(req).Elem()

		//ut, _ := request.GetMeta(ctx, "ut")
		//lang, _ := request.GetMeta(ctx, "lang")
		//fmt.Println("-----------lang" + lang)

		return handler(ctx, req)
	}
}
