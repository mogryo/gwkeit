package validator

import (
	"fmt"
	"net/url"

	"github.com/gwkeit/dto"
	"github.com/gwkeit/slicelib"
)

func ValidateTags(tags []string) []string {
	validationErrors := make([]string, 0)
	if len(tags) == 0 {
		validationErrors = append(validationErrors, "Tags cannot be empty.")
	}

	return validationErrors
}

func ValidateUrls(urls []string) []string {
	validationErrors := make([]string, 0)

	if len(urls) > 0 {
		for _, entry := range urls {
			_, err := url.ParseRequestURI(entry)
			if err != nil {
				validationErrors = append(validationErrors, fmt.Sprintf("'%s' is not a valid URL.", entry))
			}
		}
	}

	return validationErrors
}

func ValidateBody(body string) []string {
	validationErrors := make([]string, 0)

	if len(body) == 0 {
		validationErrors = append(validationErrors, "Body cannot be empty.")
	}

	return validationErrors
}

func ValidateTitle(title string) []string {
	validationErrors := make([]string, 0)

	if len(title) == 0 {
		validationErrors = append(validationErrors, "Title cannot be empty.")
	}

	return validationErrors
}

func ValidateSnippet(snippet *dto.Snippet) []string {
	tagErrors := ValidateTags(snippet.Tags)
	urlErrors := ValidateUrls(snippet.UrlList)
	bodyErrors := ValidateBody(snippet.Body)
	titleErrors := ValidateTitle(snippet.Title)

	return slicelib.Concat(bodyErrors, titleErrors, tagErrors, urlErrors)
}
