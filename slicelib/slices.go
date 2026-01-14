package slicelib

func TakeLast[T any](list []T, n int) []T {
	if len(list) > n {
		copyArray := make([]T, n)
		copy(copyArray, (list)[len(list)-n:])
		return copyArray
	}

	copyArray := make([]T, len(list))
	copy(copyArray, list)
	return copyArray
}

func Concat[T any](list ...[]T) []T {
	var totalLen int
	for _, s := range list {
		totalLen += len(s)
	}

	aggregatedSlices := make([]T, totalLen)
	var i int
	for _, s := range list {
		i += copy(aggregatedSlices[i:], s)
	}
	return aggregatedSlices
}
