package custom

import (
	"fmt"
	"testing"
)

var TestKey = []byte("UneLqzFTLwjoch8lYJybbMU4urWbhzm6")

func TestEncode(t *testing.T) {
	rs, err := Encode(TestKey, []byte("app_id=stargo"))
	fmt.Println(string(rs), err)
}

func TestDecode(t *testing.T) {
	rs, err := Decode("GE2kaDGssEPEpNRwYuu--xWMn12pBdg-2iqMQjFaL2zmXR8efX91M5F2QwA28EHh_P6dgVtZ6QFEG5sL_Y04GKJrOjM=", string(TestKey))
	fmt.Println(string(rs), err)

}
