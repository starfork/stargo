package custom

import (
	"fmt"
	"testing"
)

func TestEncode(t *testing.T) {
	rs, err := Encode(Key, []byte("app_id=stargo"))
	fmt.Println(string(rs), err)
}

func TestDecode(t *testing.T) {
	rs, err := Decode("nvK7GdnGaocZ0mJBTkkkr3JJ4J0hKLoAQwsbYwsDHw+RmXsZmVwtzkKK7Huk0dk=", string(Key))
	fmt.Println(string(rs), err)

}
