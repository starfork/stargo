package mysql

import (
	"strings"
	"time"

	"gorm.io/gorm"
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

// Timezome time qquery
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

func Now() string {
	var cstSh, _ = time.LoadLocation(TIME_LOCATION)
	return time.Now().In(cstSh).Format(TFORMAT)
}

func ParseTime(format, time_str string) (time.Time, error) {
	loc, err := time.LoadLocation(TIME_LOCATION)
	if err != nil {
		return time.Time{}, err
	}
	return time.ParseInLocation(format, time_str, loc)

}

func int2time(stamp int64, format ...string) string {
	f := TFORMAT
	if len(format) > 0 {
		f = format[0]
	}
	return time.Unix(stamp, 0).Format(f)
}

// func formatMoney(money float64) float64 {
// 	value, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", money), 64)
// 	return value
// }

func Distance(point string, dist uint32) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		l := strings.Replace(point, ",", " ", -1)
		if dist < 1000 || dist > 10000 {
			dist = 1000
		}
		return db.Where("ST_Distance_Sphere(ST_GeomFromText(\"POINT("+l+")\"),location) < ?", dist)
	}
}
