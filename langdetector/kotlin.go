package langdetector

import (
	"regexp"
	"strings"
)

func IsKotlin(src string) bool {
	src = strings.ReplaceAll(src, "\r\n", "\n")
	src = strings.ReplaceAll(src, "\r", "\n")

	lines := strings.Split(src, "\n")
	if len(lines) == 0 {
		return false
	}

	ident := `[A-Za-z_][A-Za-z0-9_]*`
	indent := `[ \t]*`

	// Basic line forms
	reBlankOrComment := regexp.MustCompile(`^\s*(//.*)?$`)
	reBlockCommentStart := regexp.MustCompile(`^\s*/\*`)
	reBlockCommentEnd := regexp.MustCompile(`\*/`)

	rePackage := regexp.MustCompile(`^` + indent + `package\s+` + ident + `(\.` + ident + `)*\s*$`)
	reImport := regexp.MustCompile(`^` + indent + `import\s+` + ident + `(\.` + ident + `)*(\.\*)?(\s+as\s+` + ident + `)?\s*$`)

	reFun := regexp.MustCompile(`^` + indent + `fun\s+` + ident + `(\s*<[^>]+>)?\s*\([^)]*\)\s*(?::\s*[^=]+)?\s*(\{|=)?\s*(//.*)?$`)
	reClassOrObject := regexp.MustCompile(`^` + indent + `(data\s+)?(sealed\s+)?(open\s+)?(abstract\s+)?(class|interface|object)\s+` + ident + `(\s*<[^>]+>)?(\s*\([^)]*\))?(\s*:\s*.+)?\s*(\{)?\s*(//.*)?$`)

	reValVar := regexp.MustCompile(`^` + indent + `(val|var)\s+` + ident + `(\s*:\s*[^=]+)?(\s*=\s*.+)?\s*(//.*)?$`)

	reIf := regexp.MustCompile(`^` + indent + `if\s*\(.+\)\s*(\{|$|//.*$)`)
	reFor := regexp.MustCompile(`^` + indent + `for\s*\(.+\)\s*(\{|$|//.*$)`)
	reWhile := regexp.MustCompile(`^` + indent + `while\s*\(.+\)\s*(\{|$|//.*$)`)
	reWhen := regexp.MustCompile(`^` + indent + `when\s*(\(.+\))?\s*\{\s*(//.*)?$`)

	reTry := regexp.MustCompile(`^` + indent + `try\s*\{\s*(//.*)?$`)
	reCatch := regexp.MustCompile(`^` + indent + `catch\s*\(\s*` + ident + `\s*:\s*` + ident + `(\.` + ident + `)*\s*\)\s*\{\s*(//.*)?$`)
	reFinally := regexp.MustCompile(`^` + indent + `finally\s*\{\s*(//.*)?$`)

	reReturn := regexp.MustCompile(`^` + indent + `return(\s+.+)?\s*(//.*)?$`)
	reThrow := regexp.MustCompile(`^` + indent + `throw\s+.+\s*(//.*)?$`)
	reBreak := regexp.MustCompile(`^` + indent + `break(\s*@` + ident + `)?\s*(//.*)?$`)
	reContinue := regexp.MustCompile(`^` + indent + `continue(\s*@` + ident + `)?\s*(//.*)?$`)

	// Kotlin-safe "expression-ish" line: allow dots, calls, generics, lambdas, etc.
	// Explicitly reject some very non-Kotlin signatures.
	reExpr := regexp.MustCompile(`^` + indent + `[^;]+?\s*(//.*)?$`)

	// Quick rejections for other languages to reduce false positives.
	reRejectPython := regexp.MustCompile(`^\s*(def|class|elif|except|with)\b|\:\s*(#.*)?$`)
	reRejectRuby := regexp.MustCompile(`^\s*(end|elsif|unless|begin|rescue)\b|@` + ident + `\s*=|(^|[\s\(\[\{,]):` + ident + `\b`)
	reRejectGo := regexp.MustCompile(`^\s*func\b|:=`)
	reRejectJS := regexp.MustCompile(`^\s*function\b`)

	// naive single-line string removal for bracket counting
	reSingleQuoted := regexp.MustCompile(`'(?:\\.|[^'\\])*'`)
	reDoubleQuoted := regexp.MustCompile(`"(?:\\.|[^"\\])*"`)

	anyNonBlank := false
	inBlockComment := false
	continuationDepth := 0

	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		anyNonBlank = true

		// Handle /* ... */ block comments (very simply).
		if inBlockComment {
			if reBlockCommentEnd.MatchString(line) {
				inBlockComment = false
			}
			continue
		}
		if reBlockCommentStart.MatchString(line) && !reBlockCommentEnd.MatchString(line) {
			inBlockComment = true
			continue
		}

		trim := strings.TrimSpace(line)

		// Quick reject patterns from other languages
		if reRejectJS.MatchString(trim) || reRejectGo.MatchString(trim) || reRejectPython.MatchString(trim) || reRejectRuby.MatchString(trim) {
			return false
		}

		inContinuation := continuationDepth > 0

		// If we are in continuation, allow broader expression lines (still reject obvious foreign tokens above).
		if inContinuation {
			if reBlankOrComment.MatchString(line) || reExpr.MatchString(line) {
				// ok
			} else {
				return false
			}
		} else {
			// Not in continuation: require line to match at least one Kotlin-ish construct.
			if reBlankOrComment.MatchString(line) ||
				rePackage.MatchString(line) ||
				reImport.MatchString(line) ||
				reFun.MatchString(line) ||
				reClassOrObject.MatchString(line) ||
				reValVar.MatchString(line) ||
				reIf.MatchString(line) ||
				reFor.MatchString(line) ||
				reWhile.MatchString(line) ||
				reWhen.MatchString(line) ||
				reTry.MatchString(line) ||
				reCatch.MatchString(line) ||
				reFinally.MatchString(line) ||
				reReturn.MatchString(line) ||
				reThrow.MatchString(line) ||
				reBreak.MatchString(line) ||
				reContinue.MatchString(line) ||
				reExpr.MatchString(line) {
				// ok
			} else {
				return false
			}
		}

		// Update continuationDepth for next line by counting (), [], {} ignoring simple strings and // comments.
		countLine := line
		countLine = reSingleQuoted.ReplaceAllString(countLine, `''`)
		countLine = reDoubleQuoted.ReplaceAllString(countLine, `""`)

		// remove // comment part
		if idx := strings.Index(countLine, "//"); idx >= 0 {
			countLine = countLine[:idx]
		}

		for i := 0; i < len(countLine); i++ {
			switch countLine[i] {
			case '(', '[', '{':
				continuationDepth++
			case ')', ']', '}':
				if continuationDepth > 0 {
					continuationDepth--
				}
			}
		}
	}

	return anyNonBlank
}
