package mysql

import (
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"

	"github.com/go-gorm/caches/v4"
	"github.com/starfork/stargo/store"
)

// Page fenye
// page 分页数，lmt每页限制条数
func Page(page, lmt uint32) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page <= 0 {
			page = 1
		}
		//var limit uint32
		//limit = 10
		if lmt > 100 {
			lmt = 100
		}
		if lmt <= 0 {
			lmt = 10
		}
		// var Arr = [6]uint32{5, 10, 20, 30, 50, 100}
		// for _, v := range Arr {
		// 	if lmt == v {
		// 		limit = v
		// 	}
		// }
		offset := (page - 1) * lmt
		//fmt.Println(lmt)

		return db.Offset(int(offset)).Limit(int(lmt))
	}
}

func int2time(stamp int64, format ...string) string {
	f := store.TFORMAT
	if len(format) > 0 {
		f = format[0]
	}
	if store.TZ1K {
		stamp /= 1e3
	}
	return time.Unix(stamp, 0).Format(f)
}

// Timezome time query
// tz["from"] 开始时间戳,tz["to"]截止时间戳
// field 时间字段名
func Timezome(tz map[string]int64, field string, format ...string) func(db *gorm.DB) *gorm.DB {

	return func(db *gorm.DB) *gorm.DB {
		if tz["from"] != 0 && tz["to"] != 0 {
			return db.Where(field+" BETWEEN ? AND ?", int2time(tz["from"], format...), int2time(tz["to"], format...))
		}
		if tz["from"] != 0 && tz["to"] == 0 {
			return db.Where(field+" >= ?", int2time(tz["from"], format...))
		}
		if tz["from"] == 0 && tz["to"] != 0 {
			return db.Where(field+" <= ?", int2time(tz["to"], format...))
		}
		return db
	}
}

func Distance(point string, dist uint32) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		l := strings.Replace(point, ",", " ", -1)
		if dist < 1000 || dist > 10000 {
			dist = 1000
		}
		return db.Where("ST_Distance_Sphere(ST_GeomFromText(\"POINT("+l+")\"),location) < ?", dist)
	}
}

// Like("field","value%")==> Where field LIKE value%
func Like(column, pattern string) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("? LIKE ?", column, pattern)
	}
}

func Cache(key string) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		v, ok := db.InstanceGet("gorm:caches")
		if !ok {
			// 没找到插件就原样返回
			return db
		}

		cachePlugin := v.(*caches.Caches)

		fmt.Println(cachePlugin)

		return db
	}
}
