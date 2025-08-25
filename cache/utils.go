package cache

import (
	"math"
)

const (
	MinUint32 uint32 = 0
	MinUint64 uint64 = 0
)

func Incr(originVal any) (any, error) {
	switch val := originVal.(type) {
	case int:
		tmp := val + 1
		if val > 0 && tmp < 0 {
			return nil, ErrIncrementOverflow
		}
		return tmp, nil
	case int32:
		if val == math.MaxInt32 {
			return nil, ErrIncrementOverflow
		}
		return val + 1, nil
	case int64:
		if val == math.MaxInt64 {
			return nil, ErrIncrementOverflow
		}
		return val + 1, nil
	case uint:
		tmp := val + 1
		if tmp < val {
			return nil, ErrIncrementOverflow
		}
		return tmp, nil
	case uint32:
		if val == math.MaxUint32 {
			return nil, ErrIncrementOverflow
		}
		return val + 1, nil
	case uint64:
		if val == math.MaxUint64 {
			return nil, ErrIncrementOverflow
		}
		return val + 1, nil
	default:
		return nil, ErrNotIntegerType
	}
}

func Decr(originVal any) (any, error) {
	switch val := originVal.(type) {
	case int:
		tmp := val - 1
		if val < 0 && tmp > 0 {
			return nil, ErrDecrementOverflow
		}
		return tmp, nil
	case int32:
		if val == math.MinInt32 {
			return nil, ErrDecrementOverflow
		}
		return val - 1, nil
	case int64:
		if val == math.MinInt64 {
			return nil, ErrDecrementOverflow
		}
		return val - 1, nil
	case uint:
		if val == 0 {
			return nil, ErrDecrementOverflow
		}
		return val - 1, nil
	case uint32:
		if val == MinUint32 {
			return nil, ErrDecrementOverflow
		}
		return val - 1, nil
	case uint64:
		if val == MinUint64 {
			return nil, ErrDecrementOverflow
		}
		return val - 1, nil
	default:
		return nil, ErrNotIntegerType
	}
}
