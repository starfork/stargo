package parser

import (
	"strconv"
)

//ParserIdCard 解析身份证号。返回男女。1:男，2:女
func ParserIdCard(id string) (bd string, gd string) {
	bd = id[6:14]
	num, _ := strconv.Atoi(id[16:18])
	if num/2 == 1 || num == 0 {
		gd = "2"
	} else {
		gd = "1"
	}

	return bd, gd
}
