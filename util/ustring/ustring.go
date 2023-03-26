package ustring

import (
	"math/rand"
	"strconv"
	"strings"
	"time"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

//New New
func New(l int, t string) string {
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

func Number(l int) string {
	return New(l, "number")
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

func ToCamel(str string) string {
	s := strings.Split(str, "_")
	if len(s) == 1 {
		return str
	}
	var tmp string
	for _, v := range s {
		tmp += strings.Title(strings.ToLower(v))
	}

	return tmp
}

func BuildSN(app int32, sn uint64) string {
	return strconv.Itoa(int(app)) + strconv.FormatUint(sn, 10)
}
func BuildSNS(app int32, sn string) string {
	return strconv.Itoa(int(app)) + sn
}

func ParseSN(strsn string) (pa int32, sn uint64) {
	app, _ := strconv.Atoi(strsn[:3])
	sn, _ = strconv.ParseUint(strsn[3:], 10, 64)
	return int32(app), sn
}

func Title(s string) string {
	return cases.Title(language.English).String(s)
}

//如果str2不是空则返回str2，否则返回str1
func OrString(str1, str2 string) string {
	if str2 != "" {
		return str2
	}
	return str1
}