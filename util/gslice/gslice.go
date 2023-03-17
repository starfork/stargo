package gslice

// 各种类型的slice
type Number interface {
	int | float32 | float64 | string | uint32
}
type Slice[T Number] []T

func New[T Number](a []T) Slice[T] {
	return a
}

func (s Slice[T]) Contains(key T) bool {
	for _, v := range s {
		if v == key {
			return true
		}
	}
	return false
}

//通过函数过滤确认是否包含

func (s Slice[T]) ContainsFilter(f func(key T) bool) bool {
	for _, v := range s {
		if f(v) {
			return true
		}
	}
	return false
}

func (s Slice[T]) Default(k, v T) T {
	if s.Contains(k) {
		return k
	}
	return v
}

func (s Slice[T]) Sum() T {
	var sum T
	for _, v := range s {
		sum += v
	}
	return sum
}

// 取一个
func (s Slice[T]) One(idx int) T {
	return s[idx : idx+1][0]
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
