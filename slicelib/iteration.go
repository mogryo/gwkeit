package slicelib

func Map[T, R any](array []T, f func(T) R) []R {
	mappedArray := make([]R, len(array))
	for i, entry := range array {
		mappedArray[i] = f(entry)
	}

	return mappedArray
}

func Filter[T any](array []T, f func(T) bool) []T {
	filteredArray := make([]T, 0)
	for _, entry := range array {
		if f(entry) {
			filteredArray = append(filteredArray, entry)
		}
	}

	return filteredArray
}

func Unique[T comparable](list []T) []T {
	seenElements := make(map[T]struct{})
	var uniqueList []T

	for _, item := range list {
		if _, exists := seenElements[item]; !exists {
			seenElements[item] = struct{}{}
			uniqueList = append(uniqueList, item)
		}
	}

	return uniqueList
}

func UniqueGet[T any, K comparable](list []T, get func(e T) K) []T {
	seenElements := make(map[K]struct{})
	var uniqueList []T

	for _, item := range list {
		identifier := get(item)
		if _, exists := seenElements[identifier]; !exists {
			seenElements[identifier] = struct{}{}
			uniqueList = append(uniqueList, item)
		}
	}

	return uniqueList
}
