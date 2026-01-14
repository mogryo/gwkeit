package dto

import "github.com/gwkeit/transform"

type Snippet struct {
	Body        string
	Title       string
	Description string
	UrlText     string
	UrlList     []string
	Tags        []string
}

func NewSnippetFromFields(
	titleField string,
	bodyField string,
	descriptionField string,
	urlsField string,
) *Snippet {
	return &Snippet{
		Title:       transform.CleanupString(titleField),
		Body:        transform.CleanupString(bodyField),
		Tags:        transform.FieldTitleAndDescToTagList(titleField, descriptionField),
		UrlList:     transform.FieldUrlsToUrlList(urlsField),
		UrlText:     transform.CleanupString(urlsField),
		Description: transform.CleanupString(descriptionField),
	}
}
