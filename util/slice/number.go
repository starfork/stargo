package slice

// 数字类型特有的
type Number interface {
	int | float32 | float64 | string | uint32 | uint64
}

type SliceNumber[T Number] []T

func NewNumber[T Number](a []T) SliceNumber[T] {
	return a
}

// 求和
func (s SliceNumber[T]) Sum() T {
	var sum T
	for _, v := range s {
		sum += v
	}
	return sum
}

// 取最大值。
func (s SliceNumber[T]) Max() T {
	var max T = s[0]
	for _, item := range s {
		if item > max {
			max = item
		}
	}
	return max
}
