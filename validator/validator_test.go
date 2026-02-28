package validator

import (
	"testing"

	"github.com/gwkeit/dto"
	"github.com/stretchr/testify/assert"
)

func Test_ValidateTags(t *testing.T) {
	t.Run("should return empty list", func(t *testing.T) {
		validationResult := ValidateTags([]string{"tag1", "tag2"})
		assert.Empty(t, validationResult)
	})
	t.Run("should return error, if no tags provided", func(t *testing.T) {
		validationResult := ValidateTags([]string{})
		assert.Equal(t, "Tags cannot be empty.", validationResult[0])
	})
}

func Test_ValidateUrls(t *testing.T) {
	t.Run("should return empty list, if empty list provided", func(t *testing.T) {
		validationResult := ValidateUrls([]string{})
		assert.Empty(t, validationResult)
	})
	t.Run("should return empty list, if valid single url provided", func(t *testing.T) {
		validationResult := ValidateUrls([]string{"https://example.com"})
		assert.Empty(t, validationResult)
	})
	t.Run("should return error, if invalid url provided", func(t *testing.T) {
		validationResult := ValidateUrls([]string{"asd"})
		assert.Equal(t, "'asd' is not a valid URL.", validationResult[0])
	})
}

func Test_ValidateBody(t *testing.T) {
	t.Run("should return empty list, if body is not empty", func(t *testing.T) {
		validationResult := ValidateBody("body")
		assert.Empty(t, validationResult)
	})
	t.Run("should return error, if body is empty", func(t *testing.T) {
		validationResult := ValidateBody("")
		assert.Equal(t, "Body cannot be empty.", validationResult[0])
	})
}

func Test_ValidateTitle(t *testing.T) {
	t.Run("should return empty list, if title is not empty", func(t *testing.T) {
		validationResult := ValidateTitle("title")
		assert.Empty(t, validationResult)
	})
	t.Run("should return error, if title is empty", func(t *testing.T) {
		validationResult := ValidateTitle("")
		assert.Equal(t, "Title cannot be empty.", validationResult[0])
	})
}

func Test_ValidateSnippet(t *testing.T) {
	t.Run("should return errors for all fields", func(t *testing.T) {
		validationResult := ValidateSnippet(&dto.Snippet{
			Title:    "",
			Body:     "",
			UrlList:  []string{"asd"},
			Tags:     []string{},
			Language: "",
		})
		assert.Equal(t, 5, len(validationResult))
	})
}
