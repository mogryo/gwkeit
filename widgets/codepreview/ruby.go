package codepreview

import (
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/gwkeit/uibuilder"
)

func addRubyDynamicColors(codeTheme *uibuilder.CodeThemeConfig, text string) string {
	var b strings.Builder
	b.Grow(len(text) + len(text)/8)

	keywords := map[string]struct{}{
		"BEGIN": {}, "END": {}, "alias": {}, "and": {}, "begin": {}, "break": {},
		"case": {}, "class": {}, "def": {}, "defined?": {}, "do": {}, "else": {},
		"elsif": {}, "end": {}, "ensure": {}, "false": {}, "for": {}, "if": {},
		"in": {}, "module": {}, "next": {}, "nil": {}, "not": {}, "or": {},
		"redo": {}, "rescue": {}, "retry": {}, "return": {}, "self": {}, "super": {},
		"then": {}, "true": {}, "undef": {}, "unless": {}, "until": {}, "when": {},
		"while": {}, "yield": {},
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

		if r == '#' {
			j := i + size
			for j < len(text) && text[j] != '\n' {
				j++
			}
			b.WriteString(codeTheme.Comment)
			b.WriteString(escape(text[i:j]))
			b.WriteString(colorReset)
			i = j
			continue
		}

		if r == '=' && (i == 0 || text[i-1] == '\n') && strings.HasPrefix(text[i:], "=begin") {
			j := i + 2
			for j < len(text) {
				if text[j] == '\n' && j+1 < len(text) {
					lineStart := j + 1
					k := lineStart
					for k < len(text) && (text[k] == ' ' || text[k] == '\t') {
						k++
					}
					if k < len(text) && strings.HasPrefix(text[k:], "=end") {
						j = k + len("=end")
						break
					}
				}
				j++
			}
			if j > len(text) {
				j = len(text)
			}
			b.WriteString(codeTheme.Comment)
			b.WriteString(escape(text[i:j]))
			b.WriteString(colorReset)
			i = j
			continue
		}

		if r == '"' || r == '\'' || r == '`' {
			quote := r
			j := i + size
			for j < len(text) {
				rr, ss := utf8.DecodeRuneInString(text[j:])
				if rr == '\\' && j+ss < len(text) && quote != '`' {
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

		if r == ':' && i+size < len(text) {
			next, nextSize := utf8.DecodeRuneInString(text[i+size:])
			if isIdentStart(next) {
				j := i + size + nextSize
				for j < len(text) {
					rr, ss := utf8.DecodeRuneInString(text[j:])
					if !isIdentPart(rr) {
						break
					}
					j += ss
				}
				b.WriteString(codeTheme.Keyword)
				b.WriteString(escape(text[i:j]))
				b.WriteString(colorReset)
				i = j
				continue
			}
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
