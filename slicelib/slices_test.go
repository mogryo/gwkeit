package slicelib

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_TakeLast(t *testing.T) {
	t.Run("should return empty list, if input list is empty", func(t *testing.T) {
		res := TakeLast([]int{}, 1)
		assert.Empty(t, res)
	})
	t.Run("should return all elements from list, if provided N is bigger than list length", func(t *testing.T) {
		res := TakeLast([]int{1, 2}, 5)
		assert.Equal(t, []int{1, 2}, res)
	})
	t.Run("should return empty list, if provided N is zero", func(t *testing.T) {
		res := TakeLast([]int{1, 2}, 0)
		assert.Empty(t, res)
	})
	t.Run("should return single one element", func(t *testing.T) {
		res := TakeLast([]int{1, 2}, 1)
		assert.Equal(t, []int{2}, res)
	})
}

func Test_Concat(t *testing.T) {
	t.Run("should return empty list, if both lists are empty", func(t *testing.T) {
		res := Concat([]int{}, []int{})
		assert.Empty(t, res)
	})
	t.Run("should return the provided one list", func(t *testing.T) {
		res := Concat([]int{1, 2})
		assert.Equal(t, []int{1, 2}, res)
	})
	t.Run("should return concatenation of three lists", func(t *testing.T) {
		res := Concat([]int{1, 2}, []int{3, 4}, []int{5, 6})
		assert.Equal(t, []int{1, 2, 3, 4, 5, 6}, res)
	})
}
