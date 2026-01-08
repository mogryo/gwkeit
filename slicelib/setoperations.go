package slicelib

import "slices"

func Difference[T comparable](a, b []T) []T {
	diff := make([]T, 0)
	for _, aEntry := range a {
		if !slices.Contains(b, aEntry) {
			diff = append(diff, aEntry)
		}
	}

	return diff
}

func DifferenceGetA[T any, K comparable](a []T, b []K, get func(entry T) K) []T {
	diff := make([]T, 0)
	for _, aEntry := range a {
		if !slices.Contains(b, get(aEntry)) {
			diff = append(diff, aEntry)
		}
	}

	return diff
}

func DifferenceGetB[T comparable, K any](a []T, b []K, get func(entry K) T) []T {
	diff := make([]T, 0)
	parsedB := Map(b, get)
	for _, aEntry := range a {
		if !slices.Contains(parsedB, aEntry) {
			diff = append(diff, aEntry)
		}
	}

	return diff
}
