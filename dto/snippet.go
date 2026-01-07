package dto

import "github.com/gwkeit/transform"

type Snippet struct {
	Title string
	Body  string
	Tags  []string
	Urls  []string
}

func NewSnippetFromFields(titleField string, bodyField string, descriptionField string, urlsField string) *Snippet {
	return &Snippet{
		Title: transform.CleanupTitle(titleField),
		Body:  transform.CleanupBody(bodyField),
		Tags:  transform.FormDescriptionToTagList(descriptionField),
		Urls:  transform.FormUrlsToUrlList(urlsField),
	}
}
