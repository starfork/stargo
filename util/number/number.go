package number

import (
	"crypto/rand"
	"errors"
	"math"
	"math/big"
)

// 返回一个指定范围的随机数
func RangeRand(min, max int64) (int64, error) {
	if min > max {
		return 0, errors.New("the min is greater than max")
	}

	if min < 0 {
		f64Min := math.Abs(float64(min))
		i64Min := int64(f64Min)
		result, _ := rand.Int(rand.Reader, big.NewInt(max+1+i64Min))

		return result.Int64() - i64Min, nil
	} else {
		result, _ := rand.Int(rand.Reader, big.NewInt(max-min+1))
		return min + result.Int64(), nil
	}
}
