package slicelib

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Map(t *testing.T) {
	t.Run("should return doubled integers", func(t *testing.T) {
		list := []int{1, 2, 3, 4}

		mappedList := Map(list, func(i int) int {
			return i * 2
		})
		assert.Equal(t, []int{2, 4, 6, 8}, mappedList)
	})
	t.Run("should return list of strings", func(t *testing.T) {
		list := []int{1, 2, 3, 4}

		mappedList := Map(list, func(i int) string {
			return strconv.Itoa(i)
		})
		assert.Equal(t, []string{"1", "2", "3", "4"}, mappedList)
	})
}

func Test_Unique(t *testing.T) {
	t.Run("should return empty list, if input list is empty", func(t *testing.T) {
		res := Unique([]int{})
		assert.Empty(t, res)
	})
	t.Run("should return unique elements from list", func(t *testing.T) {
		res := Unique([]int{1, 2, 3, 2, 1})
		assert.Equal(t, []int{1, 2, 3}, res)
	})
	t.Run("should return same elements, if no duplicates", func(t *testing.T) {
		res := Unique([]int{1, 2, 3})
		assert.Equal(t, []int{1, 2, 3}, res)
	})
}

func Test_UniqueGet(t *testing.T) {
	t.Run("should return empty list, if input list is empty", func(t *testing.T) {
		res := UniqueGet([]int{}, func(e int) int { return e })
		assert.Empty(t, res)
	})
	t.Run("should return unique elements from list", func(t *testing.T) {
		res := UniqueGet(
			[]string{"one", "two", "three", "four", "five"},
			func(e string) uint8 { return e[0] },
		)
		assert.Equal(t, []string{"one", "two", "four"}, res)
	})
	t.Run("should return same elements, if no duplicates", func(t *testing.T) {
		res := UniqueGet([]string{"one", "two", "four"}, func(e string) uint8 { return e[0] })
		assert.Equal(t, []string{"one", "two", "four"}, res)
	})
}
