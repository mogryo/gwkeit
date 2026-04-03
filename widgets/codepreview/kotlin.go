package codepreview

import (
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/gwkeit/uibuilder"
)

func addKotlinDynamicColors(codeTheme *uibuilder.CodeThemeConfig, text string) string {
	var b strings.Builder
	b.Grow(len(text) + len(text)/8)

	keywords := map[string]struct{}{
		"class": {}, "object": {}, "interface": {}, "fun": {}, "val": {}, "var": {},
		"if": {}, "else": {}, "when": {}, "for": {}, "while": {}, "do": {},
		"return": {}, "break": {}, "continue": {}, "this": {}, "super": {}, "is": {},
		"in": {}, "as": {}, "typealias": {}, "package": {}, "import": {}, "try": {},
		"catch": {}, "finally": {}, "throw": {}, "typeof": {},
		"public": {}, "private": {}, "protected": {}, "internal": {}, "open": {},
		"final": {}, "abstract": {}, "override": {}, "const": {}, "lateinit": {},
		"data": {}, "sealed": {}, "enum": {}, "inner": {}, "annotation": {}, "companion": {},
		"suspend": {}, "operator": {}, "infix": {}, "inline": {}, "tailrec": {},
		"external": {}, "expect": {}, "actual": {},
		"true": {}, "false": {}, "null": {}, "field": {}, "it": {},
	}

	escape := func(s string) string {
		s = strings.ReplaceAll(s, "[", "[[")
		return s
	}

	isIdentStart := func(r rune) bool { return r == '_' || unicode.IsLetter(r) }
	isIdentPart := func(r rune) bool { return r == '_' || unicode.IsLetter(r) || unicode.IsDigit(r) }
	isDigit := func(r rune) bool { return r >= '0' && r <= '9' }

	i := 0
	for i < len(text) {
		r, size := utf8.DecodeRuneInString(text[i:])

		if r == '/' && i+1 < len(text) {
			next, _ := utf8.DecodeRuneInString(text[i+1:])
			if next == '/' {
				j := i + 2
				for j < len(text) && text[j] != '\n' {
					j++
				}
				b.WriteString(codeTheme.Comment)
				b.WriteString(escape(text[i:j]))
				b.WriteString(colorReset)
				i = j
				continue
			}
			if next == '*' {
				j := i + 2
				for j+1 < len(text) && !(text[j] == '*' && text[j+1] == '/') {
					j++
				}
				if j+1 < len(text) {
					j += 2
				} else {
					j = len(text)
				}
				b.WriteString(codeTheme.Comment)
				b.WriteString(escape(text[i:j]))
				b.WriteString(colorReset)
				i = j
				continue
			}
		}

		if r == '"' && strings.HasPrefix(text[i:], `"""`) {
			j := i + 3
			for j+2 < len(text) && !(text[j] == '"' && text[j+1] == '"' && text[j+2] == '"') {
				j++
			}
			if j+2 < len(text) {
				j += 3
			} else {
				j = len(text)
			}
			b.WriteString(codeTheme.String)
			b.WriteString(escape(text[i:j]))
			b.WriteString(colorReset)
			i = j
			continue
		}

		if r == '"' || r == '\'' {
			quote := r
			j := i + size
			for j < len(text) {
				rr, ss := utf8.DecodeRuneInString(text[j:])
				if rr == '\\' && j+ss < len(text) {
					j += ss
					_, ss2 := utf8.DecodeRuneInString(text[j:])
					j += ss2
					continue
				}
				if rr == quote {
					j += ss
					break
				}
				j += ss
			}
			b.WriteString(codeTheme.String)
			b.WriteString(escape(text[i:j]))
			b.WriteString(colorReset)
			i = j
			continue
		}

		if unicode.IsSpace(r) {
			b.WriteRune(r)
			i += size
			continue
		}

		if isIdentStart(r) {
			j := i + size
			for j < len(text) {
				rr, ss := utf8.DecodeRuneInString(text[j:])
				if !isIdentPart(rr) {
					break
				}
				j += ss
			}
			word := text[i:j]
			if _, ok := keywords[word]; ok {
				b.WriteString(codeTheme.Keyword)
				b.WriteString(word)
				b.WriteString(colorReset)
			} else {
				b.WriteString(codeTheme.Identifier)
				b.WriteString(escape(word))
				b.WriteString(colorReset)
			}
			i = j
			continue
		}

		if isDigit(r) {
			j := i + size
			for j < len(text) {
				rr, ss := utf8.DecodeRuneInString(text[j:])
				if !(unicode.IsDigit(rr) || rr == '.' || rr == 'x' || rr == 'X' || rr == 'e' || rr == 'E' || rr == 'p' || rr == 'P' || rr == '+' || rr == '-' || (rr >= 'a' && rr <= 'f') || (rr >= 'A' && rr <= 'F')) {
					break
				}
				j += ss
			}
			b.WriteString(codeTheme.Number)
			b.WriteString(escape(text[i:j]))
			b.WriteString(colorReset)
			i = j
			continue
		}

		b.WriteString(escape(string(r)))
		i += size
	}

	return b.String()
}
