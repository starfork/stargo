package validator

//先按照 https://github.com/favadi/protoc-go-inject-tag 给对应的struct创建tag
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
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
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

	return func(ctx context.Context, req interface{}, in *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (resp interface{}, err error) {

		vfields := []string{}
		method := ""
		if md, ok := metadata.FromIncomingContext(ctx); ok {
			m := md.Get("g-method")
			if len(md.Get("g-method")) > 0 {
				method = strings.ToLower(m[0])
			}

		}
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
