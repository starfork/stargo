package config

import (
	"fmt"
	"reflect"
	"testing"
)

func TestGetStep(t *testing.T) {
	c := &Config{}

	val := `200|3
500|5`
	c.SetVal("test", val)
	v := c.GetStep(203, "test")
	fmt.Printf("v:%+v\n", v)
}

func TestConvertKey(t *testing.T) {
	var key interface{}
	var abc uint32 = 112312
	key = abc

	rf := reflect.TypeOf(abc)
	fmt.Println(rf.Kind())

	fmt.Println(key)
}
