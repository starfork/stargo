package gslice

// 各种类型的slice
type Number interface {
	int | float32 | float64 | string | uint32 | uint64
}
type Slice[T Number] []T

func New[T Number](a []T) Slice[T] {
	return a
}

// 包含，f过滤函数
func (s Slice[T]) Contains(key T, f ...func(k T) bool) bool {
	var fn func(k T) bool

	if len(f) == 0 {
		fn = func(v T) bool {
			return key == v
		}
	} else {
		fn = f[0]
	}

	for _, v := range s {
		if fn(v) {
			return true
		}
	}
	return false
}

// 默认返回.不包含k，则返回v
func (s Slice[T]) Default(k, v T) T {
	if s.Contains(k) {
		return k
	}
	return v
}

// 求和
func (s Slice[T]) Sum() T {
	var sum T
	for _, v := range s {
		sum += v
	}
	return sum
}

// 取最大值。
func (s Slice[T]) Max() T {
	var max T = s[0]
	for _, item := range s {
		if item > max {
			max = item
		}
	}
	return max
}

// 取一个.非随机取一个
func (s Slice[T]) One(index ...int) T {
	max := len(s)
	if max == 1 {
		return s[0]
	}
	var idx int = 0
	if len(index) > 0 {
		idx = index[0]
	}
	if idx > max {
		idx = max - 1
	}
	return s[idx : idx+1][0]
}

// 过滤
func (s Slice[T]) Filter(fn func(T) bool) []T {
	var res []T
	for _, item := range s {
		if fn(item) {
			res = append(res, item)
		}
	}
	return res
}

// Tail 获取切片尾部元素
// dv: 空切片默认值
func (s Slice[T]) Tail(dv ...T) T {
	if s == nil && len(dv) > 0 {
		return dv[0]
	}
	return s[len(s)-1]
}

// 交集
func (s Slice[T]) Intersect(b Slice[T]) Slice[T] {
	inter := make([]T, 0)
	mp := make(map[T]bool)
	for _, sa := range s {
		if _, ok := mp[sa]; !ok {
			mp[sa] = true
		}
	}
	for _, sb := range b {
		if _, ok := mp[sb]; ok {
			inter = append(inter, sb)
		}
	}
	return inter
}

// 差集
func (s Slice[T]) Diff(b Slice[T]) Slice[T] {
	diff := make([]T, 0)
	mp := make(map[T]bool)
	for _, sa := range s {
		if _, ok := mp[sa]; !ok {
			mp[sa] = true
		}
	}
	for _, sb := range b {
		if ok := mp[sb]; ok {
			delete(mp, sb)
		}
	}
	for k := range mp {
		diff = append(diff, k)
	}

	return diff
}

// 并集
func (s Slice[T]) Union(b Slice[T]) Slice[T] {
	union := make([]T, 0)
	mp := make(map[T]bool)
	s = append(s, b...)
	for _, v := range s {
		if ok := mp[v]; ok {
			continue
		}
		mp[v] = true
		union = append(union, v)
	}

	return union
}
