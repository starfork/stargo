package pm

import (
	"fmt"
	"testing"
)

func TestEncodeURL(t *testing.T) {
	data := Pm{
		"abc": "sdfsdf",
		"amt": 100.01,
		"ddd": "",
		"efg": 0,
	}
	fmt.Println(data.EncodeURL())
}
