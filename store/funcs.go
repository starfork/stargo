package store

import (
	"time"
)

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

// func formatMoney(money float64) float64 {
// 	value, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", money), 64)
// 	return value
// }
