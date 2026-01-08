package slicelib

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Difference(t *testing.T) {
	t.Run("should return empty list, if two lists are equal", func(t *testing.T) {
		diff := Difference(
			[]int{1, 2, 3, 4, 5},
			[]int{1, 2, 3, 4, 5},
		)
		assert.Empty(t, diff)
	})
	t.Run("should return list of different elements", func(t *testing.T) {
		diff := Difference(
			[]int{1, 3},
			[]int{1, 2, 4, 5},
		)
		assert.Equal(t, []int{3}, diff)
	})
	t.Run("should return empty list, if list A is empty", func(t *testing.T) {
		diff := Difference(
			[]int{},
			[]int{1, 2, 4, 5},
		)
		assert.Empty(t, diff)
	})
	t.Run("should return whole list A, if list B is empty", func(t *testing.T) {
		diff := Difference(
			[]int{1, 2, 3, 4},
			[]int{},
		)
		assert.Equal(t, []int{1, 2, 3, 4}, diff)
	})
	t.Run("should return empty list, if both lists are empty", func(t *testing.T) {
		diff := Difference(
			[]int{},
			[]int{},
		)
		assert.Empty(t, diff)
	})
}

func Test_DifferenceGetA(t *testing.T) {
	t.Run("should return empty list, if two lists are equal", func(t *testing.T) {
		diff := DifferenceGetA(
			[]int{1, 2, 3, 4, 5},
			[]int{1, 2, 3, 4, 5},
			func(e int) int { return e },
		)
		assert.Empty(t, diff)
	})
	t.Run("should return empty list, if both lists are empty", func(t *testing.T) {
		diff := DifferenceGetA(
			[]int{},
			[]int{},
			func(e int) int { return e },
		)
		assert.Empty(t, diff)
	})
	t.Run("should return elements with type from list A", func(t *testing.T) {
		diff := DifferenceGetA(
			[]int{1, 2, 3, 4, 6},
			[]string{"1", "2", "3", "4", "5"},
			func(e int) string { return strconv.Itoa(e) },
		)
		assert.Equal(t, []int{6}, diff)
	})
}

func Test_DifferenceGetB(t *testing.T) {
	t.Run("should return empty list, if two lists are equal", func(t *testing.T) {
		diff := DifferenceGetB(
			[]int{1, 2, 3, 4, 5},
			[]int{1, 2, 3, 4, 5},
			func(e int) int { return e },
		)
		assert.Empty(t, diff)
	})
	t.Run("should return empty list, if both lists are empty", func(t *testing.T) {
		diff := DifferenceGetB(
			[]int{},
			[]int{},
			func(e int) int { return e },
		)
		assert.Empty(t, diff)
	})
	t.Run("should return elements with type from list A", func(t *testing.T) {
		diff := DifferenceGetB(
			[]int64{1, 2, 3, 4, 6},
			[]string{"1", "2", "3", "4", "5"},
			func(e string) int64 {
				val, _ := strconv.ParseInt(e, 10, 64)
				return val
			},
		)
		assert.Equal(t, []int64{6}, diff)
	})
}
