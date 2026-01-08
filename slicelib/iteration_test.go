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
