package jwt_test

import (
	"fmt"
	"public/pkg/util/crypt/jwt"
	"testing"
)

func TestEncode(t *testing.T) {
	jwt, _ := jwt.New()
	fmt.Printf("jwt%+v", jwt)
}
