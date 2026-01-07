package slicelib

func TakeLast[T any](array []T, n int) []T {
	if len(array) > n {
		copyArray := make([]T, n)
		copy(copyArray, (array)[len(array)-n:])
		return copyArray
	}

	copyArray := make([]T, len(array))
	copy(copyArray, array)
	return copyArray
}

func Concat[T any](slices ...[]T) []T {
	var totalLen int
	for _, s := range slices {
		totalLen += len(s)
	}

	aggregatedSlices := make([]T, totalLen)
	var i int
	for _, s := range slices {
		i += copy(aggregatedSlices[i:], s)
	}
	return aggregatedSlices
}
