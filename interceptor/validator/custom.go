package validator

import (
	"reflect"
	"regexp"
	"strconv"

	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

var tagNameFunc_vLabel = func(field reflect.StructField) string {
	//fmt.Println(field.Tag)
	label := field.Tag.Get("vlabel")
	if label == "" {
		return field.Name
	}
	return label
}

var registerFn_Money = func(ut ut.Translator) error {
	return ut.Add("money", "金额格式不正确", false) // see universal-translator for details
}

var translationFn_Money = func(ut ut.Translator, fe validator.FieldError) string {
	t, _ := ut.T("money", fe.Field())
	return t
}

func ValidateMoney(fl validator.FieldLevel) bool {
	pattern := regexp.MustCompile(`^\d+(\.\d{1,2})?$`)
	return pattern.MatchString(strconv.FormatFloat(fl.Field().Float(), 'f', -1, 64))
}
