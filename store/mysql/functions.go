package mysql

import (
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
func Timezome(tz map[string]int64, field string) func(db *gorm.DB) *gorm.DB {

	return func(db *gorm.DB) *gorm.DB {
		//fmt.Println(tz["from"] / 1000)
		if tz["from"] != 0 && tz["to"] != 0 {

			return db.Where(field+" BETWEEN ? AND ?", int2time(tz["from"]), int2time(tz["to"]))
		}
		if tz["from"] != 0 && tz["to"] == 0 {
			return db.Where(field+" >= ?", int2time(tz["from"]))
		}
		if tz["from"] == 0 && tz["to"] != 0 {
			return db.Where(field+" <= ?", int2time(tz["to"]))
		}
		return db

	}
}

// TimeAdded
// tz["from"] 开始时间戳,tz["to"]截止时间戳
// field 时间字段名
// ts n为空 m为本月 d为当天
func TimeAdded(tz map[string]int64, field string, ts string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if ts == "n" {
			if tz["from"] != 0 && tz["to"] != 0 {
				return db.Where(field+" BETWEEN ? AND ?", int2time(tz["from"]), int2time(tz["to"]))
			}
			if tz["from"] != 0 && tz["to"] == 0 {
				return db.Where(field+" >= ?", int2time(tz["from"]))
			}
			if tz["from"] == 0 && tz["to"] != 0 {
				return db.Where(field+" <= ?", int2time(tz["to"]))
			}
		} else {
			if tz["from"] == 0 && tz["to"] == 0 {
				timeNow := time.Now()
				//本月
				if ts == "m" {
					from := timeNow.Unix() //当前时间戳
					//月初时间
					to := time.Date(timeNow.Year(), timeNow.Month(), 1, 0, 0, 0, 0, timeNow.Location()).Unix()
					return db.Where(field+" BETWEEN ? AND ?", int2time(from), int2time(to))
				}
				//当天
				if ts == "d" {
					//当前开始时间
					from := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.Now().Location()).Unix()
					to := timeNow.Unix() //当前时间戳
					return db.Where(field+" BETWEEN ? AND ?", int2time(from), int2time(to))
				}

			}
		}

		return db

	}
}

func Now() string {
	var cstSh, _ = time.LoadLocation(TIME_LOCATION)
	format := "2006-01-02 15:04:05"
	return time.Now().In(cstSh).Format(format)
}

func ParseTime(format, time_str string) (time.Time, error) {
	loc, err := time.LoadLocation(TIME_LOCATION)
	if err != nil {
		return time.Time{}, err
	}
	return time.ParseInLocation(format, time_str, loc)

}

func int2time(stamp int64) string {
	return time.Unix(stamp, 0).Format("2006-01-02 15:04:05")
}

// func formatMoney(money float64) float64 {
// 	value, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", money), 64)
// 	return value
// }
