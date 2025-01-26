package merger

import (
	"fmt"
	"testing"
)

type Result struct {
	Uid       string
	UserName  string
	Age       int
	IsDeleted bool
	F1        string
}
type UserInfo struct {
	Uid   string
	Name  string
	Age   int
	Fuser string
}

type UserInfoList struct {
	UserInfoList []*UserInfo
}
type OriList struct {
	Data []*Result
}

var rsA = &OriList{
	Data: []*Result{
		{Uid: "1", UserName: ""},
		{Uid: "2", UserName: ""},
	},
}
var rB = &UserInfoList{
	UserInfoList: []*UserInfo{
		{Uid: "1", Name: "Name1", Fuser: "fuser1"},
		{Uid: "2", Name: "Name2", Fuser: "fuser2"},
	},
}

func TestMerge(t *testing.T) {

	Merge(rsA.Data, rB.UserInfoList,
		func(a *Result) string { return a.Uid },
		func(b *UserInfo) string { return b.Uid },
		func(a *Result, b *UserInfo) {
			a.UserName = b.Name
			a.F1 = b.Fuser
		})

	for _, v := range rsA.Data {
		fmt.Printf("%+v \n", v)
	}
}

func TestReduce(t *testing.T) {
	uids := Reduce(rsA.Data, func(r *Result) string { return r.Uid })
	fmt.Println(uids)
}
