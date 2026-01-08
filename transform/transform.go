package transform

import (
	"regexp"
	"strings"

	"github.com/gwkeit/gwkeitdb"
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

func FieldDescriptionToTagList(tags string) []string {
	reg, _ := regexp.Compile(`\s+`)
	s := reg.ReplaceAllString(strings.TrimSpace(tags), " ")
	if len(s) == 0 {
		return []string{}
	}
	s = strings.Trim(s, " ")
	return strings.Split(s, " ")
}

func TagListToFieldDescription(tags []gwkeitdb.Tag) string {
	return strings.Join(slicelib.Map(tags, func(tag gwkeitdb.Tag) string {
		return tag.Tag
	}), " ")
}

func UrlListToFieldUrls(urls []gwkeitdb.Url) string {
	return strings.Join(slicelib.Map(urls, func(url gwkeitdb.Url) string {
		return url.Url
	}), "\n")
}

func CleanupBody(body string) string {
	return strings.TrimSpace(body)
}

func CleanupTitle(title string) string {
	return strings.TrimSpace(title)
}
