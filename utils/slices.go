package utils

import "slices"

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

func Difference[T comparable](a, b []T) []T {
	diff := make([]T, 0)
	for _, aEntry := range a {
		if !slices.Contains(b, aEntry) {
			diff = append(diff, aEntry)
		}
	}

	return diff
}

func Map[T, R any](array []T, f func(T) R) []R {
	mappedArray := make([]R, len(array))
	for i, entry := range array {
		mappedArray[i] = f(entry)
	}

	return mappedArray
}
