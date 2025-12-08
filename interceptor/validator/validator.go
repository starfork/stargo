package validator

import (
	"bytes"
	"context"
	"fmt"
	"reflect"
	"strings"
	"sync"

	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"
	"github.com/starfork/stargo/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	validate *validator.Validate
	trans    ut.Translator
	once     sync.Once
)

func initValidator() {
	once.Do(func() {
		validate = validator.New()
		uni := ut.New(zh.New())
		trans, _ = uni.GetTranslator("zh")

		validate.RegisterTagNameFunc(tagNameFunc_vLabel)
		validate.RegisterValidation("money", ValidateMoney)
		_ = zh_translations.RegisterDefaultTranslations(validate, trans)
		validate.RegisterTranslation("money", trans, registerFn_Money, translationFn_Money)

		validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := fld.Tag.Get("mapstructure")
			if name == "-" {
				return ""
			}
			return name
		})
	})
}

// Unary Interceptor
func Unary() grpc.UnaryServerInterceptor {
	initValidator() // 确保全局 validator 初始化一次

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

		err = validate.StructExcept(req, vfields...) // <--- 只调用，不再注册规则

		if err != nil {
			if tErrs, ok := err.(validator.ValidationErrors); !ok {
				return resp, status.New(codes.InvalidArgument, fmt.Sprintf("error%s", err)).Err()
			} else {
				var buf bytes.Buffer
				for _, s2 := range tErrs.Translate(trans) {
					buf.WriteString(s2 + ",")
				}
				return resp, status.New(codes.InvalidArgument, buf.String()).Err()
			}
		}
		return handler(ctx, req)
	}
}
