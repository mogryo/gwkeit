package transform

import (
	"regexp"
	"strings"

	"github.com/gwkeit/slicelib"
)

func FieldUrlsToUrlList(urls string) []string {
	s := strings.TrimSpace(urls)

	if len(s) == 0 {
		return []string{}
	}
	splitUrls := strings.Split(strings.TrimSpace(s), "\n")
	return slicelib.Map(splitUrls, strings.TrimSpace)
}

func FieldTitleAndDescToTagList(title, description string) []string {
	reg, _ := regexp.Compile(`\s+`)
	trimmedTitle := reg.ReplaceAllString(strings.TrimSpace(title), " ")
	trimmedDescription := reg.ReplaceAllString(strings.TrimSpace(description), " ")

	var splitTitle []string
	if len(trimmedTitle) > 0 {
		splitTitle = strings.Split(strings.Trim(trimmedTitle, " "), " ")
	} else {
		splitTitle = []string{}
	}

	var splitDescription []string
	if len(trimmedDescription) > 0 {
		splitDescription = strings.Split(strings.Trim(trimmedDescription, " "), " ")
	} else {
		splitDescription = []string{}
	}
	joinedTags := slicelib.Concat(splitTitle, splitDescription)
	lowercaseTags := slicelib.Map(joinedTags, strings.ToLower)

	return slicelib.Unique(lowercaseTags)
}

func CleanupString(value string) string {
	return strings.TrimSpace(value)
}
