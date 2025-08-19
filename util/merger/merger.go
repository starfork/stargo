package merger

type Key[T any] func(item T) any

func Merge[A any, B any](rsA []A, rsB []B, keyA Key[A], keyB Key[B], merger func(A, B)) {
	bMap := make(map[any]B)
	for _, b := range rsB {
		bMap[keyB(b)] = b
	}
	for i := range rsA {
		if b, exists := bMap[keyA(rsA[i])]; exists {
			merger(rsA[i], b)
		}
	}
}
func Reduce[T any, R any](input []T, mapper func(T) R) []R {
	result := make([]R, len(input))
	for i, item := range input {
		result[i] = mapper(item)
	}
	return result
}

func MergeAppend[A any, B any](
	rsA []A,
	rsB []B,
	keyA func(A) any,
	keyB func(B) any,
	appender func(A, B),
) {
	bMap := make(map[any][]B)
	for _, b := range rsB {
		k := keyB(b)
		bMap[k] = append(bMap[k], b) // 支持一对多
	}
	for _, a := range rsA {
		if bs, ok := bMap[keyA(a)]; ok {
			for _, b := range bs {
				appender(a, b) // 注意 a、b 不做拷贝，直接传递
			}
		}
	}
}
