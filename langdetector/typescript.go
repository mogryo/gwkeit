package langdetector

import (
	"regexp"
	"strings"
)

func IsTypeScript(src string) bool {
	src = strings.ReplaceAll(src, "\r\n", "\n")
	src = strings.ReplaceAll(src, "\r", "\n")

	lines := strings.Split(src, "\n")
	if len(lines) == 0 {
		return false
	}

	ident := `[A-Za-z_\$][A-Za-z0-9_\$]*`
	indent := `[ \t]*`

	reBlank := regexp.MustCompile(`^\s*$`)
	reLineComment := regexp.MustCompile(`^\s*//.*$`)
	reBlockCommentStart := regexp.MustCompile(`^\s*/\*`)
	reBlockCommentEnd := regexp.MustCompile(`\*/`)

	// Module-ish
	reImport := regexp.MustCompile(`^` + indent + `import\b.*$`)
	reExport := regexp.MustCompile(`^` + indent + `export\b.*$`)
	reType := regexp.MustCompile(`^` + indent + `type\s+` + ident + `\s*=.*;?\s*$`)
	reInterface := regexp.MustCompile(`^` + indent + `interface\s+` + ident + `\b.*\{?\s*$`)
	reEnum := regexp.MustCompile(`^` + indent + `enum\s+` + ident + `\b.*\{?\s*$`)
	reNamespace := regexp.MustCompile(`^` + indent + `(declare\s+)?namespace\s+` + ident + `\b.*\{?\s*$`)

	// Declarations / statements
	reClass := regexp.MustCompile(`^` + indent + `(export\s+)?(abstract\s+)?class\s+` + ident + `\b.*\{?\s*$`)
	reFn := regexp.MustCompile(`^` + indent + `(export\s+)?(async\s+)?function\s+` + ident + `\s*\(.*\)\s*(?::\s*[^ {]+.*)?\s*\{?\s*$`)
	reArrowFnAssign := regexp.MustCompile(`^` + indent + `(export\s+)?(const|let|var)\s+` + ident + `(\s*:\s*[^=]+)?\s*=\s*(async\s+)?\(?[^=]*\)?\s*=>.*$`)
	reVar := regexp.MustCompile(`^` + indent + `(export\s+)?(const|let|var)\s+` + ident + `(\s*:\s*[^=;]+)?(\s*=\s*.+)?;?\s*$`)

	reIf := regexp.MustCompile(`^` + indent + `if\s*\(.+\)\s*\{?\s*$`)
	reElse := regexp.MustCompile(`^` + indent + `else\s*\{?\s*$`)
	reFor := regexp.MustCompile(`^` + indent + `for\s*\(.+\)\s*\{?\s*$`)
	reWhile := regexp.MustCompile(`^` + indent + `while\s*\(.+\)\s*\{?\s*$`)
	reSwitch := regexp.MustCompile(`^` + indent + `switch\s*\(.+\)\s*\{?\s*$`)
	reCase := regexp.MustCompile(`^` + indent + `(case\s+.+:|default\s*:)\s*$`)
	reTry := regexp.MustCompile(`^` + indent + `try\s*\{?\s*$`)
	reCatch := regexp.MustCompile(`^` + indent + `catch\s*\(.*\)\s*\{?\s*$`)
	reFinally := regexp.MustCompile(`^` + indent + `finally\s*\{?\s*$`)

	reReturn := regexp.MustCompile(`^` + indent + `return(\s+.+)?;?\s*$`)
	reThrow := regexp.MustCompile(`^` + indent + `throw\s+.+;?\s*$`)
	reBreak := regexp.MustCompile(`^` + indent + `break(\s+` + ident + `)?;?\s*$`)
	reContinue := regexp.MustCompile(`^` + indent + `continue(\s+` + ident + `)?;?\s*$`)

	// General expression line (allow ; optional; disallow obvious foreign tokens).
	reExpr := regexp.MustCompile(`^` + indent + `.+;?\s*$`)

	// Quick rejects for other languages to reduce false positives.
	reRejectPython := regexp.MustCompile(`^\s*(def|class|elif|except|with)\b|:\s*(#.*)?$`)
	reRejectRuby := regexp.MustCompile(`^\s*(end|elsif|unless|begin|rescue)\b|@` + ident + `|\.\.|:\w+\b`)
	reRejectGo := regexp.MustCompile(`^\s*package\b|\bfunc\b|:=`)
	// Kotlin hint: `fun name(` or `val name`
	reRejectKotlin := regexp.MustCompile(`^\s*(fun|val|var|when)\b`)

	// naive single-line string removal for bracket counting
	reSingleQuoted := regexp.MustCompile(`'(?:\\.|[^'\\])*'`)
	reDoubleQuoted := regexp.MustCompile(`"(?:\\.|[^"\\])*"`)
	reBacktick := regexp.MustCompile("`(?:\\\\.|[^`\\\\])*`")

	anyNonBlank := false
	inBlockComment := false
	continuationDepth := 0

	for _, line := range lines {
		if reBlank.MatchString(line) {
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

		// Quick reject patterns from other languages (to avoid "anything" matching reExpr).
		if reRejectGo.MatchString(trim) || reRejectPython.MatchString(trim) || reRejectRuby.MatchString(trim) || reRejectKotlin.MatchString(trim) {
			return false
		}

		inContinuation := continuationDepth > 0

		// Not in continuation: require at least one TS-ish form or a plausible statement/expression.
		if !inContinuation {
			// Also quickly reject lines that look like Python block headers.
			if strings.HasSuffix(trim, ":") && strings.HasPrefix(trim, "if ") {
				return false
			}

			if reLineComment.MatchString(line) ||
				reImport.MatchString(line) || reExport.MatchString(line) ||
				reType.MatchString(line) || reInterface.MatchString(line) || reEnum.MatchString(line) || reNamespace.MatchString(line) ||
				reClass.MatchString(line) || reFn.MatchString(line) || reArrowFnAssign.MatchString(line) || reVar.MatchString(line) ||
				reIf.MatchString(line) || reElse.MatchString(line) || reFor.MatchString(line) || reWhile.MatchString(line) ||
				reSwitch.MatchString(line) || reCase.MatchString(line) ||
				reTry.MatchString(line) || reCatch.MatchString(line) || reFinally.MatchString(line) ||
				reReturn.MatchString(line) || reThrow.MatchString(line) || reBreak.MatchString(line) || reContinue.MatchString(line) ||
				reExpr.MatchString(line) {
				// ok
			} else {
				return false
			}
		} else {
			// In continuation: accept most non-empty lines that aren't obviously foreign.
			if !(reLineComment.MatchString(line) || reExpr.MatchString(line)) {
				return false
			}
		}

		// Update continuationDepth for next line by counting (), [], {} ignoring simple strings and // comments.
		countLine := line
		countLine = reSingleQuoted.ReplaceAllString(countLine, `''`)
		countLine = reDoubleQuoted.ReplaceAllString(countLine, `""`)
		countLine = reBacktick.ReplaceAllString(countLine, "``")

		// remove // comment part (naively; ignores // inside strings which we removed above)
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
