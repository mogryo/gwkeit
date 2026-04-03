package codepreview

import (
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/gwkeit/uibuilder"
)

func addTypeScriptDynamicColors(codeTheme *uibuilder.CodeThemeConfig, text string) string {
	var b strings.Builder
	b.Grow(len(text) + len(text)/8)

	keywords := map[string]struct{}{
		"break": {}, "case": {}, "catch": {}, "class": {}, "const": {},
		"continue": {}, "debugger": {}, "default": {}, "delete": {}, "do": {},
		"else": {}, "enum": {}, "export": {}, "extends": {}, "false": {},
		"finally": {}, "for": {}, "function": {}, "if": {}, "import": {},
		"in": {}, "instanceof": {}, "new": {}, "null": {}, "return": {},
		"super": {}, "switch": {}, "this": {}, "throw": {}, "true": {},
		"try": {}, "typeof": {}, "var": {}, "void": {}, "while": {},
		"with": {}, "as": {}, "implements": {}, "interface": {}, "let": {},
		"package": {}, "private": {}, "protected": {}, "public": {}, "static": {},
		"yield": {}, "any": {}, "boolean": {}, "number": {}, "string": {},
		"symbol": {}, "type": {}, "from": {}, "of": {}, "async": {}, "await": {},
		"readonly": {}, "declare": {}, "namespace": {}, "abstract": {},
	}

	escape := func(s string) string {
		s = strings.ReplaceAll(s, "[", "[[")
		return s
	}

	isIdentStart := func(r rune) bool {
		return r == '_' || r == '$' || unicode.IsLetter(r)
	}
	isIdentPart := func(r rune) bool {
		return r == '_' || r == '$' || unicode.IsLetter(r) || unicode.IsDigit(r)
	}
	isDigit := func(r rune) bool {
		return r >= '0' && r <= '9'
	}

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

		if r == '"' || r == '\'' || r == '`' {
			quote := r
			j := i + size
			for j < len(text) {
				rr, ss := utf8.DecodeRuneInString(text[j:])
				if quote == '`' {
					if rr == '`' {
						j += ss
						break
					}
					j += ss
					continue
				}
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
				if !(unicode.IsDigit(rr) || rr == '.' || rr == '_' || rr == 'x' || rr == 'X' || rr == 'e' || rr == 'E' || rr == 'p' || rr == 'P' || rr == '+' || rr == '-' || (rr >= 'a' && rr <= 'f') || (rr >= 'A' && rr <= 'F')) {
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
