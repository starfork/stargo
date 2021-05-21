package jwt

import (
	"fmt"
	"testing"
)

func Test_Demo(t *testing.T) {
	claims := UserClaims{
		UID: 123,
	}
	jwt := New(
		Key("sfsdfsd"),
		Claims(claims),
	)
	c, _ := jwt.Encode()
	fmt.Println(c)
	jwt.SetToken(c)
	d, _ := jwt.Decode()

	fmt.Println(d["uid"])

	//fmt.Println(d["uid"])
}
