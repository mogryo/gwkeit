package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_IfElse(t *testing.T) {
	t.Run("should return first value if condition true", func(t *testing.T) {
		res := IfElse(true, 1, 2)
		assert.Equal(t, 1, res)
	})
	t.Run("should return second value if condition false", func(t *testing.T) {
		res := IfElse(false, 1, 2)
		assert.Equal(t, 2, res)
	})
}
