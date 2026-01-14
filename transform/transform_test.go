package transform

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_FieldUrlsToUrlList(t *testing.T) {
	t.Run("should return empty list if empty string", func(t *testing.T) {
		val := FieldUrlsToUrlList("")
		assert.Empty(t, val)
	})
	t.Run("should return empty list if string only with whitespace", func(t *testing.T) {
		val := FieldUrlsToUrlList("    ")
		assert.Empty(t, val)
	})
	t.Run("should return list with single url", func(t *testing.T) {
		val := FieldUrlsToUrlList("  https://randompage.io  ")
		assert.Equal(t, []string{"https://randompage.io"}, val)
	})
	t.Run("should return list with multiple url", func(t *testing.T) {
		val := FieldUrlsToUrlList("  https://randompage.io  \n https://second.io")
		assert.Equal(t, []string{"https://randompage.io", "https://second.io"}, val)
	})
}

func Test_FieldDescriptionToTagList(t *testing.T) {
	t.Run("should return empty list for empty fields", func(t *testing.T) {
		val := FieldTitleAndDescToTagList("", "")
		assert.Empty(t, val)
	})
	t.Run("should return single tag", func(t *testing.T) {
		val := FieldTitleAndDescToTagList("one", "  ")
		assert.Equal(t, []string{"one"}, val)
	})
	t.Run("should return multiple tags", func(t *testing.T) {
		val := FieldTitleAndDescToTagList("one  two three   ", "asd")
		assert.Equal(t, []string{"one", "two", "three", "asd"}, val)
	})
}

func Test_CleanupString(t *testing.T) {
	t.Run("should return empty string", func(t *testing.T) {
		val := CleanupString("")
		assert.Empty(t, val)
	})
	t.Run("should return string without whitespaces in the beginning and end", func(t *testing.T) {
		val := CleanupString("   test   ")
		assert.Equal(t, "test", val)
	})
}
