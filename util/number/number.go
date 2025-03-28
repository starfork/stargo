package number

import (
	"crypto/rand"
	"math"
	"math/big"
)

// 返回一个指定范围的随机数
func RangeRand(min, max int64) (int64, error) {
	if min > max {
		//return 0, errors.New("the min is greater than max")
		min, max = max, min
	}

	if min < 0 {
		f64Min := math.Abs(float64(min))
		i64Min := int64(f64Min)
		result, _ := rand.Int(rand.Reader, big.NewInt(max+1+i64Min))

		return result.Int64() - i64Min, nil
	}

	result, _ := rand.Int(rand.Reader, big.NewInt(max-min+1))
	return min + result.Int64(), nil

}

func Pager(page, l uint32) (offset, lmt uint32) {
	lmt = l
	if page <= 0 {
		page = 1
	}
	if l > 100 {
		lmt = 100
	}
	if l <= 0 {
		lmt = 10
	}
	offset = (page - 1) * lmt

	return offset, lmt

}
