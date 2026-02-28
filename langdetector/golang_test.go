package langdetector

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_IsGo(t *testing.T) {
	t.Run("should return true for golang function declaration with package name", func(t *testing.T) {
		src := "package main\nfunc main() {}\n"
		assert.True(t, IsGo(src))
	})
	t.Run("should return true for golang function declaration, no package name", func(t *testing.T) {
		src := "func main() {\n}\n"
		assert.True(t, IsGo(src))
	})
	t.Run("should return true for golang loop", func(t *testing.T) {
		src := "{\nfor _, i := range entries {\n}\n}"
		assert.True(t, IsGo(src))
	})
	t.Run("should return true for golang error check ", func(t *testing.T) {
		src := "if err != nil {\nreturn err\n}"
		assert.True(t, IsGo(src))
	})
	t.Run("should return true for golang type struct", func(t *testing.T) {
		src := "type User struct {\nName string\nAge int\n}"
		assert.True(t, IsGo(src))
	})
	t.Run("should return false for empty string", func(t *testing.T) {
		assert.False(t, IsGo(""))
	})
	t.Run("should return false for javascript function declaration", func(t *testing.T) {
		src := "function main() {}\n"
		assert.False(t, IsGo(src))
	})
	t.Run("should return false for python function declaration", func(t *testing.T) {
		src := "def main():\nreturn 1\n"
		assert.False(t, IsGo(src))
	})
}
