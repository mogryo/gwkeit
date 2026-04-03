package codepreview

import (
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/gwkeit/uibuilder"
)

func addGoDynamicColors(codeTheme *uibuilder.CodeThemeConfig, text string) string {
	var b strings.Builder
	b.Grow(len(text) + len(text)/8)

	keywords := map[string]struct{}{
		"break": {}, "case": {}, "chan": {}, "const": {}, "continue": {},
		"default": {}, "defer": {}, "else": {}, "fallthrough": {}, "for": {},
		"func": {}, "go": {}, "goto": {}, "if": {}, "import": {}, "interface": {},
		"map": {}, "package": {}, "range": {}, "return": {}, "select": {},
		"struct": {}, "switch": {}, "type": {}, "var": {},
		"bool": {}, "byte": {}, "complex64": {}, "complex128": {}, "error": {},
		"float32": {}, "float64": {}, "int": {}, "int8": {}, "int16": {}, "int32": {},
		"int64": {}, "rune": {}, "string": {}, "uint": {}, "uint8": {}, "uint16": {},
		"uint32": {}, "uint64": {}, "uintptr": {},
		"true": {}, "false": {}, "nil": {},
		"iota": {}, "make": {}, "new": {}, "len": {}, "cap": {}, "append": {},
		"copy": {}, "close": {}, "delete": {}, "complex": {}, "real": {}, "imag": {},
		"panic": {}, "recover": {},
	}

	escape := func(s string) string {
		// Escape markup-sensitive characters, but keep normal text readable.
		s = strings.ReplaceAll(s, "[", "[[")
		return s
	}

	isIdentStart := func(r rune) bool {
		return r == '_' || unicode.IsLetter(r)
	}
	isIdentPart := func(r rune) bool {
		return r == '_' || unicode.IsLetter(r) || unicode.IsDigit(r)
	}
	isDigit := func(r rune) bool {
		return r >= '0' && r <= '9'
	}

	i := 0
	for i < len(text) {
		r, size := utf8.DecodeRuneInString(text[i:])

		// Comments
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

		// String / raw string literals
		if r == '"' || r == '`' || r == '\'' {
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

		// Whitespace
		if unicode.IsSpace(r) {
			b.WriteRune(r)
			i += size
			continue
		}

		// Identifiers / keywords
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

		// Numbers
		if isDigit(r) {
			j := i + size
			prev := r
			for j < len(text) {
				rr, ss := utf8.DecodeRuneInString(text[j:])
				if !(unicode.IsDigit(rr) || rr == '.' || rr == 'x' || rr == 'X' || rr == 'e' || rr == 'E' || rr == 'p' || rr == 'P' || rr == '+' || rr == '-' || (rr >= 'a' && rr <= 'f') || (rr >= 'A' && rr <= 'F')) {
					break
				}
				j += ss
				isHexAlpha := (rr >= 'a' && rr <= 'f') || (rr >= 'A' && rr <= 'F')
				if unicode.IsDigit(rr) || rr == '.' || rr == 'x' || rr == 'X' || rr == 'e' || rr == 'E' || rr == 'p' || rr == 'P' || isHexAlpha {
					j += ss
					prev = rr
					continue
				}
				if (rr == '+' || rr == '-') && (prev == 'e' || prev == 'E' || prev == 'p' || prev == 'P') {
					j += ss
					prev = rr
					continue
				}
				break
			}
			b.WriteString(codeTheme.Number)
			b.WriteString(escape(text[i:j]))
			b.WriteString(colorReset)
			i = j
			continue
		}

		// Everything else
		b.WriteString(escape(string(r)))
		i += size
	}

	return b.String()
}
