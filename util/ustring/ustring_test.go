package ustring

import (
	"fmt"
	"testing"
)

func TestSnakeString(t *testing.T) {
	fmt.Println(SnakeString("AbcDef1A"))
}

func TestCamelString(t *testing.T) {
	fmt.Println(CamelString("aAcd_edf_1", true))
}

func TestVs(t *testing.T) {
	fmt.Println('z')
	fmt.Println('a')
	fmt.Println('Z')
	fmt.Println('A')
	fmt.Println('_')

	fmt.Println('a' > 'A')
	fmt.Println('a' - 'A')
	fmt.Println('a' - 32)
	fmt.Println(string('a' - 32))
	fmt.Println(string(byte(92)))
}
