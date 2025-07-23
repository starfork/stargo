package validator

import (
	"bytes"
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"
	"github.com/starfork/stargo/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Unary Interceptor
func Unary() grpc.UnaryServerInterceptor {
	var (
		validate = validator.New()
		uni      = ut.New(zh.New())
		trans, _ = uni.GetTranslator("zh")
	)
	validate.RegisterTagNameFunc(tagNameFunc_vLabel)
	///validate.RegisterTagNameFunc(tagNameFunc_vFor)
	validate.RegisterTranslation("money", trans, registerFn_Money, translationFn_Money)

	err := zh_translations.RegisterDefaultTranslations(validate, trans)
	if err != nil {
		panic(err)
	}
	validate.RegisterValidation("money", ValidateMoney)

	return func(ctx context.Context, req any, in *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (resp any, err error) {

		vfields := []string{}
		method := api.MetaMethod(ctx)

		if method != "" {
			tp := reflect.TypeOf(req).Elem()
			for i := 0; i < tp.NumField(); i++ {
				f := tp.Field(i)
				if strings.Contains(strings.ToLower(f.Tag.Get("vexcept")), method) {
					vfields = append(vfields, f.Name)
				}
			}
		}

		err = validate.StructExcept(req, vfields...)

		if err != nil {
			if tErrs, ok := err.(validator.ValidationErrors); !ok {
				return resp, status.New(codes.Unknown, fmt.Sprintf("error%s", err)).Err()
			} else {
				translations := tErrs.Translate(trans)
				var buf bytes.Buffer
				for _, s2 := range translations {
					buf.WriteString(s2)
				}
				return resp, status.New(codes.InvalidArgument, buf.String()).Err()
			}
		}
		return handler(ctx, req)
	}
}
