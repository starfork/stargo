package path

import (
	"fmt"
	"strconv"
	"strings"
	"testing"
)

const ConvertBase = 36 //转换机制，一旦程序运行，则不可修改
const MaxLevel = 10    //最大层级

func TestPath(t *testing.T) {

	p := strconv.FormatUint(117779567824864525, ConvertBase)
	fmt.Println(p)
	//var v uint64 = 117779567824864525
	fmt.Println(strconv.FormatUint(uint64(1295), ConvertBase))
	var v = ^uint64(0)
	fmt.Println(v)
	s := "D" + strconv.FormatUint(uint64(1800), ConvertBase) + strconv.FormatUint(v, ConvertBase)
	//最大长度17
	fmt.Println(strings.ToUpper(s))

}

func TestClength(t *testing.T) {
	i := 100
	for {

		s := strconv.FormatUint(uint64(i), ConvertBase)

		//fmt.Println(len(s))
		i += 1

		if len(s) > 2 {
			fmt.Println(i, s)
			break
		}

	}

}
