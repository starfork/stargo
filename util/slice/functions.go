package slice

func IsEmpty[T any](slice []T) bool {
	return slice == nil
}
