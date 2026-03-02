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

func Test_AlignTextLeft(t *testing.T) {
	t.Run("should return empty string", func(t *testing.T) {
		val, isParsed := AlignTextLeft("")
		assert.Empty(t, val)
		assert.False(t, isParsed)
	})
	t.Run("should align single line text to left", func(t *testing.T) {
		val, isParsed := AlignTextLeft("  test  ")
		assert.Equal(t, "test  ", val)
		assert.True(t, isParsed)
	})
	t.Run("should do nothing for single line with no whitespaces text", func(t *testing.T) {
		val, isParsed := AlignTextLeft("test")
		assert.Equal(t, "test", val)
		assert.False(t, isParsed)
	})
	t.Run("should align multiline text, both lines to left", func(t *testing.T) {
		val, isParsed := AlignTextLeft("  test\n  test2")
		assert.Equal(t, "test\ntest2", val)
		assert.True(t, isParsed)
	})
	t.Run("should align multiline text, both lines to left but leave one space", func(t *testing.T) {
		val, isParsed := AlignTextLeft("  test\n test2")
		assert.Equal(t, " test\ntest2", val)
		assert.True(t, isParsed)
	})
	t.Run("should do nothing for multiline text, if one line has to whitespace", func(t *testing.T) {
		val, isParsed := AlignTextLeft("test\n  test2")
		assert.Equal(t, "test\n  test2", val)
		assert.False(t, isParsed)
	})
	t.Run("", func(t *testing.T) {
		val, parsed := AlignTextLeft("\t\t\tsnippetDto := dto.NewSnippetFromFields(\n\t\t\t\tep.title.GetText(),\n\t\t\t\tep.body.GetText(),\n\t\t\t\tep.description.GetText(),\n\t\t\t\tep.urls.GetText(),\n\t\t\t\tselectedLanguage,\n\t\t\t)\n\t\t\tvalidationErrors := validator.ValidateSnippet(snippetDto)")
		assert.Equal(t, "snippetDto := dto.NewSnippetFromFields(\n\tep.title.GetText(),\n\tep.body.GetText(),\n\tep.description.GetText(),\n\tep.urls.GetText(),\n\tselectedLanguage,\n)\nvalidationErrors := validator.ValidateSnippet(snippetDto)", val)
		assert.True(t, parsed)
	})
}
