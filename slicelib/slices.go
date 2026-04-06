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

// SubSlice
/**
Returns a slice of the given slice, starting at the given index and ending at the given index.
This is a copy of the Go stdlib's Slice function, but with the bounds clamped to the length of the slice.
If param to is less than from, the returned slice will be empty.
If param to exceeds the length of the slice, to will be set to the length of the slice.

@param s The slice to return a subslice of.

@param from The index to start the subslice at.

@param to The index to end the subslice at.

@return A subslice of the given slice.
*/
func SubSlice[T any](s []T, from, to int) []T {
	initialSliceLength := len(s)
	start := from
	end := to

	// Clamp bounds to valid range.
	if start < 0 {
		start = 0
	}
	if end < 0 {
		end = 0
	}
	if start > initialSliceLength {
		start = initialSliceLength
	}
	if end > initialSliceLength {
		end = initialSliceLength
	}

	// If the range is invalid, return empty slice.
	if start >= end {
		return []T{}
	}

	// Return a copy so it is a new slice/array-backed slice.
	result := make([]T, end-start)
	copy(result, s[start:end])
	return result
}
