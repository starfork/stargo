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
