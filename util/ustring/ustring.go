package ustring

import (
	"math/rand"
	"strconv"
	"strings"
	"time"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type Ustring string

// New New
func NewRand(l int, t string) string {
	var str string
	switch t {
	case "lower":
		str = "qwertyupajshdfkznxbccmuytxgskljbsqpenxglsnllqsbtsg"
	case "upper":
		str = "QWERTYUPAJSHDFKZNXBCCMAQWJBKAYBFJYGHGKSHGAGPLANVBCRE"
	case "number":
		str = "01234567897875647512321361232130182282746154129484756"
	default:
		str = "97126123456789787564751232136123213182282746154129484756abcdefghjkl612321318228mnpqrstuvwxyzytwvx612321318228nagqpkjjsdfhbqwbs"
	}

	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < l; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

func RandNumber(l int) string {
	return NewRand(l, "number")
}

func U32String(ids []uint32) string {
	build := strings.Builder{}
	for i, v := range ids {
		build.WriteString(strconv.Itoa(int(v)))
		if i != len(ids)-1 {
			build.WriteString(",")
		}
	}
	return build.String()
}

func ToCamel(str string, tag ...language.Tag) string {
	s := strings.Split(str, "_")
	if len(s) == 1 {
		return str
	}
	var tmp string
	for _, v := range s {
		tmp += Title(strings.ToLower(v), tag...)
	}

	return tmp
}

// ucfirst
func Title(s string, tag ...language.Tag) string {
	t := language.English
	if len(tag) > 0 {
		t = tag[0]
	}
	return cases.Title(t).String(s)
}

// 如果str2不是空则返回str2，否则返回str1
func OrString(str1, str2 string) string {
	if str2 != "" {
		return str2
	}
	return str1
}

// 如果str2不是空则返回str2，否则返回str1
func Or(str1, str2 string) string {
	return OrString(str1, str2)
}

func Int(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return i
}

/** 驼峰转下划线 **/
func SnakeString(s string) string {
	num := len(s)
	data := make([]byte, 0, num*2)
	f := func(d byte) bool {
		return d >= 65 && d <= 90
	}
	if f(s[0]) {
		data = append(data, s[0]+32)
	} else {
		data = append(data, s[0])
	}
	for i := 1; i < num; i++ {
		d := s[i]
		if f(d) {
			data = append(data, '_', d+32)
		} else {
			data = append(data, d)
		}

	}
	return string(data[:])
}

// 下划线转驼峰
func CamelString(s string, ugly ...bool) string {
	l := len(s)
	data := make([]byte, 0, l)

	var conv = true
	// var ug = false
	// if len(ugly) > 0 {
	// 	ug = true
	// }
	for i := 0; i < l; i++ {
		d := s[i]
		if conv && d >= 'a' && d <= 'z' {
			d = d - 32
		}
		// else if d != '_' {
		// 	//对于aXcd_edf_1这种奇葩的“X”，是否转小写的x
		// 	// if d >= 65 && d <= 90 {
		// 	// 	d = d + 32
		// 	// }
		// 	data = append(data, d)
		// }

		conv = d == '_'
		if !conv {
			//对于aXcd_edf_1这种奇葩的“X”，是否转小写的x
			// if d >= 65 && d <= 90 && ug {
			// 	d = d + 32
			// }
			data = append(data, d)
		}

	}
	return string(data[:])
}
