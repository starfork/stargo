package merger

import (
	"fmt"
	"testing"
)

type Result struct {
	Uid       uint32
	UserName  string
	Age       int
	IsDeleted bool
	F1        string
	Tt        map[string]any
	TagIds    []uint32
	Tags      []*UserTag
}
type UserInfo struct {
	Uid   uint32
	Name  string
	Age   int
	Fuser string
	UTt   map[string]any
}
type UserTag struct {
	Id   uint32
	Name string
}

type UserInfoList struct {
	UserInfoList []*UserInfo
}
type OriList struct {
	Data []*Result
}
type TagList struct {
	List []*UserTag
}

var rsA = &OriList{
	Data: []*Result{
		{Uid: 1, UserName: "", TagIds: []uint32{1, 2}},
		{Uid: 2, UserName: ""},
	},
}
var rB = &UserInfoList{
	UserInfoList: []*UserInfo{
		{Uid: 1, Name: "Name1", Fuser: "fuser1", UTt: map[string]any{
			"abc": 1,
		}},
		{Uid: 2, Name: "Name2", Fuser: "fuser2", UTt: map[string]any{
			"def": 1,
		}},
	},
}

var rsTags = &TagList{
	List: []*UserTag{
		{Id: 1, Name: "tag1"},
		{Id: 2, Name: "tag2"},
	},
}

func TestMerge(t *testing.T) {

	Merge(rsA.Data, rB.UserInfoList,
		func(a *Result) any { return a.Uid },
		func(b *UserInfo) any { return b.Uid },
		func(a *Result, b *UserInfo) {
			a.UserName = b.Name
			a.F1 = b.Fuser
			a.Tt = b.UTt
		})
	MergeFlat(rsA.Data, rsTags.List,
		func(a *Result) []uint32 { return a.TagIds },
		func(b *UserTag) uint32 { return b.Id },
		func(a *Result, b *UserTag) {
			a.Tags = append(a.Tags, b)
		})
	for _, v := range rsA.Data {
		fmt.Printf("%+v \n", v)
	}
}

func TestMergeFlat(t *testing.T) {

	Merge(rsA.Data, rB.UserInfoList,
		func(a *Result) any { return a.Uid },
		func(b *UserInfo) any { return b.Uid },
		func(a *Result, b *UserInfo) {
			a.UserName = b.Name
			a.F1 = b.Fuser
			a.Tt = b.UTt
		})

	for _, v := range rsA.Data {
		fmt.Printf("%+v \n", v)
	}
}

func TestReduce(t *testing.T) {
	uids := Reduce(rsA.Data, func(r *Result) any { return r.Uid })
	fmt.Println(uids)
}
