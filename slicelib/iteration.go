package slicelib

func Map[T, R any](array []T, f func(T) R) []R {
	mappedArray := make([]R, len(array))
	for i, entry := range array {
		mappedArray[i] = f(entry)
	}

	return mappedArray
}
