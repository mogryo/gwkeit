package transform

import (
	"regexp"
	"slices"
	"strings"
	"unicode/utf8"

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

// AlignTextLeft removes common leading indentation from a multi-line text block.
//
// It returns the modified text and a boolean flag indicating whether the text was
// parsed and changed (true) or whether no alignment was performed (false).
//
// No alignment is performed when the input is empty, does not start with a space
// or tab, or when there is no common leading indentation across non-empty lines.
func AlignTextLeft(text string) (string, bool) {
	if utf8.RuneCountInString(text) == 0 {
		return text, false
	}
	spacingChar := rune(text[0])
	if spacingChar != ' ' && spacingChar != '\t' {
		return text, false
	}

	lines := strings.Split(text, "\n")
	smallestWhitespaceCount := slices.Max(slicelib.Map(lines, func(line string) int { return utf8.RuneCountInString(line) }))
	for _, line := range lines {
		if trimmerString := strings.TrimSpace(line); utf8.RuneCountInString(trimmerString) == 0 {
			continue
		}

		whitespaceCount := 0
		for _, char := range line {
			if char == spacingChar {
				whitespaceCount++
			} else {
				break
			}
		}

		if whitespaceCount == 0 {
			smallestWhitespaceCount = whitespaceCount
			break
		}
		if whitespaceCount < smallestWhitespaceCount {
			smallestWhitespaceCount = whitespaceCount
		}
	}

	if smallestWhitespaceCount == 0 {
		return text, false
	}

	trimmedLines := slicelib.Map(lines, func(line string) string {
		//if trimmerString := strings.TrimSpace(line); utf8.RuneCountInString(trimmerString) == 0 {
		//	return line
		//}
		return strings.TrimPrefix(line, strings.Repeat(string(spacingChar), smallestWhitespaceCount))
	})
	return strings.Join(trimmedLines, "\n"), true
}

func CleanupString(value string) string {
	return strings.TrimSpace(value)
}
