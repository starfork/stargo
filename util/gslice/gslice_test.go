package gslice_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/starfork/stargo/util/gslice"
)

func TestContainsFilter(t *testing.T) {
	noAuth := gslice.New([]string{"/v1/passport/", "/v1/public/"})
	path := "/v1/passport1/login"
	rs := noAuth.ContainsFilter(func(key string) bool {
		return strings.Contains(path, key)
	})
	fmt.Println(rs)
}

func TestAsSlice(t *testing.T) {
	str := ""
	s := gslice.New(strings.Split(str, ","))
	fmt.Println(s)
	if !s.Contains("abc") {
		s = append(s, "abc")
	}
	if !s.Contains("abc") {
		s = append(s, "abc")
	}
	ns := gslice.New(strings.Split(strings.Join(s, ","), ","))
	fmt.Println(ns)
	if !ns.Contains("def") {
		ns = append(ns, "def")
	}
	if !ns.Contains("def") {
		ns = append(ns, "def")
	}
	fmt.Println(len(ns))

	fmt.Println(strings.Join(ns, ","))
}

func TestIntersect(t *testing.T) {
	a := gslice.New([]string{"/v1/passport/", "/v1/public/"})
	b := gslice.New([]string{"/v1/passport/", "/v1/public/1"})
	fmt.Println(a.Intersect(b))
}

func TestDiff(t *testing.T) {
	a := gslice.New([]string{"/v1/passport/", "/v1/public/", "/v1/public/2"})
	b := gslice.New([]string{"/v1/passport/", "/v1/public/1"})
	fmt.Println(a.Diff(b))
}
func TestUnion(t *testing.T) {
	a := gslice.New([]string{"/v1/passport/", "/v1/public/", "/v1/public/2"})
	b := gslice.New([]string{"/v1/passport/", "/v1/public/1"})
	fmt.Println(a.Union(b))
}

func TestCalc(t *testing.T) {
	a := gslice.New([]string{}) //老权限
	//删除2，新增7，8之后
	b := gslice.New([]string{"/v1/estate/estates$$GET", "/v1/service/display/services$$GET", "/v1/service/evaluate$$GET", "/v1/service/service$$PUT", "/v1/user/auth/rule$$PUT", "/v1/user/log/login$$GET", "/v1/user/friends$$GET", "/v1/user/auth/roles$$GET", "/v1/user/certifications$$GET"})
	inter := a.Intersect(b)
	fmt.Printf("交集:%+v", inter)

	//需要删除的
	fmt.Printf("需要删除的节点:%+v", a.Diff(inter))
	add := b.Diff(inter)
	fmt.Printf("需要增加的节点:%+v", add)
	for k, v := range add {
		fmt.Printf("add-v,%d,%+v\n", k, v)
	}
}
