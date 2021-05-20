package validator

//先按照 https://github.com/favadi/protoc-go-inject-tag 给对应的struct创建tag
import (
	"bytes"
	"context"
	"fmt"

	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

//Unary Interceptor
func Unary() grpc.UnaryServerInterceptor {
	var (
		validate = validator.New()
		uni      = ut.New(zh.New())
		trans, _ = uni.GetTranslator("zh")
	)
	err := zh_translations.RegisterDefaultTranslations(validate, trans)
	if err != nil {
		panic(err)
	}
	return func(ctx context.Context, req interface{}, in *grpc.UnaryServerInfo, h grpc.UnaryHandler) (resp interface{}, err error) {

		if err := validate.Struct(req); err != nil {
			if transErr, ok := err.(validator.ValidationErrors); ok {
				translations := transErr.Translate(trans)
				var buf bytes.Buffer
				for _, s2 := range translations {
					buf.WriteString(s2)
				}
				err = status.New(codes.InvalidArgument, buf.String()).Err()
				return resp, err
			}
			err = status.New(codes.Unknown, fmt.Sprintf("error%s", err)).Err()
			return resp, err
		}
		return h(ctx, req)
	}
}
