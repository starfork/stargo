package jwt_test

import (
	"fmt"
	"testing"

	"github.com/starfork/stargo/util/crypt/jwt"
)

func TestEncode(t *testing.T) {
	jwt, _ := jwt.New()
	fmt.Printf("jwt%+v", jwt)
}
