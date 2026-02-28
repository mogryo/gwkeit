package langdetector

import (
	"go/parser"
	"go/token"
	"strings"
)

func IsGo(src string) bool {
	src = strings.TrimSpace(src)
	if src == "" {
		return false
	}

	fileSet := token.NewFileSet()

	// 1) Try parsing as a full Go file (works when src already has `package ...`).
	if _, err := parser.ParseFile(fileSet, "input.go", src, parser.AllErrors); err == nil {
		return true
	}

	// 2) Try parsing as top-level declarations by prepending a package clause.
	// This handles snippets like: `func main() {}` which are illegal inside a function body.
	withPkg := "package p\n" + src + "\n"
	if _, err := parser.ParseFile(fileSet, "decls.go", withPkg, parser.AllErrors); err == nil {
		return true
	}

	// 3) Try parsing as statements/expressions by wrapping into a function body.
	wrapped := "package p\nfunc _(){\n" + src + "\n}\n"
	if _, err := parser.ParseFile(fileSet, "snippet.go", wrapped, parser.AllErrors); err == nil {
		return true
	}

	return false
}
