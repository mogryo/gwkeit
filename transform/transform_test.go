package transform

import (
	"testing"

	"github.com/gwkeit/gwkeitdb"
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
	t.Run("should return empty list for empty description field", func(t *testing.T) {
		val := FieldDescriptionToTagList("")
		assert.Empty(t, val)
	})
	t.Run("should return single tag", func(t *testing.T) {
		val := FieldDescriptionToTagList("one")
		assert.Equal(t, []string{"one"}, val)
	})
	t.Run("should return multiple tags", func(t *testing.T) {
		val := FieldDescriptionToTagList("one  two three   ")
		assert.Equal(t, []string{"one", "two", "three"}, val)
	})
}

func Test_TagListToFieldDescription(t *testing.T) {
	t.Run("should return empty string for empty tag list", func(t *testing.T) {
		val := TagListToFieldDescription([]gwkeitdb.Tag{})
		assert.Empty(t, val)
	})
	t.Run("should return string with two tags", func(t *testing.T) {
		val := TagListToFieldDescription([]gwkeitdb.Tag{
			{
				ID:  0,
				Tag: "one",
			},
			{
				ID:  2,
				Tag: "two",
			},
		})
		assert.Equal(t, "one two", val)
	})
}

func Test_UrlListToFieldUrls(t *testing.T) {
	t.Run("should return empty list", func(t *testing.T) {
		val := UrlListToFieldUrls([]gwkeitdb.Url{})
		assert.Empty(t, val)
	})
	t.Run("should return list with single url", func(t *testing.T) {
		val := UrlListToFieldUrls([]gwkeitdb.Url{
			{
				ID:        0,
				Url:       "https://example.com",
				SnippetID: 0,
			},
		})
		assert.Equal(t, "https://example.com", val)
	})
}

func Test_CleanupBody(t *testing.T) {
	t.Run("should return empty string", func(t *testing.T) {
		val := CleanupBody("")
		assert.Empty(t, val)
	})
	t.Run("should return string without whitespaces in the beginning and end", func(t *testing.T) {
		val := CleanupBody("   test   ")
		assert.Equal(t, "test", val)
	})
}

func Test_CleanupTitle(t *testing.T) {
	t.Run("should return empty string", func(t *testing.T) {
		val := CleanupTitle("")
		assert.Empty(t, val)
	})
	t.Run("should return string without whitespaces in the beginning and end", func(t *testing.T) {
		val := CleanupTitle("   test   ")
		assert.Equal(t, "test", val)
	})
}
