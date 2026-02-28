package langdetector

import (
	"regexp"
	"strings"
)

func IsPython(src string) bool {
	src = strings.ReplaceAll(src, "\r\n", "\n")
	src = strings.ReplaceAll(src, "\r", "\n")

	lines := strings.Split(src, "\n")
	if len(lines) == 0 {
		return false
	}

	ident := `[A-Za-z_][A-Za-z0-9_]*`
	indent := `[ \t]*`
	commentOrBlank := `^\s*(#.*)?$`

	// Common Python statements (must include ':' where Python requires it).
	reDef := regexp.MustCompile(`^` + indent + `def\s+` + ident + `\s*\(.*\)\s*:\s*(#.*)?$`)
	reClass := regexp.MustCompile(`^` + indent + `class\s+` + ident + `(\s*\(.*\))?\s*:\s*(#.*)?$`)

	reIf := regexp.MustCompile(`^` + indent + `if\s+.+:\s*(#.*)?$`)
	reElif := regexp.MustCompile(`^` + indent + `elif\s+.+:\s*(#.*)?$`)
	reElse := regexp.MustCompile(`^` + indent + `else\s*:\s*(#.*)?$`)

	reFor := regexp.MustCompile(`^` + indent + `for\s+.+\s+in\s+.+:\s*(#.*)?$`)
	reWhile := regexp.MustCompile(`^` + indent + `while\s+.+:\s*(#.*)?$`)

	reTry := regexp.MustCompile(`^` + indent + `try\s*:\s*(#.*)?$`)
	reExcept := regexp.MustCompile(`^` + indent + `except(\s+.+)?\s*:\s*(#.*)?$`)
	reFinally := regexp.MustCompile(`^` + indent + `finally\s*:\s*(#.*)?$`)

	// `with ...:` (single line)
	reWith := regexp.MustCompile(`^` + indent + `with\s+.+:\s*(#.*)?$`)
	// `with open(` / `with something(` (multiline start, no colon yet)
	reWithStart := regexp.MustCompile(`^` + indent + `with\s+.+\(\s*(#.*)?$`)

	reImport := regexp.MustCompile(`^` + indent + `import\s+` + ident + `(\s+as\s+` + ident + `)?(\s*,\s*` + ident + `(\s+as\s+` + ident + `)?)*\s*(#.*)?$`)
	reFrom := regexp.MustCompile(`^` + indent + `from\s+` + ident + `(\.` + ident + `)*\s+import\s+(.+)\s*(#.*)?$`)

	reReturn := regexp.MustCompile(`^` + indent + `return(\s+.+)?\s*(#.*)?$`)
	reRaise := regexp.MustCompile(`^` + indent + `raise(\s+.+)?\s*(#.*)?$`)
	reAssert := regexp.MustCompile(`^` + indent + `assert\s+.+\s*(#.*)?$`)

	rePass := regexp.MustCompile(`^` + indent + `pass\s*(#.*)?$`)
	reBreak := regexp.MustCompile(`^` + indent + `break\s*(#.*)?$`)
	reContinue := regexp.MustCompile(`^` + indent + `continue\s*(#.*)?$`)

	reGlobal := regexp.MustCompile(`^` + indent + `global\s+` + ident + `(\s*,\s*` + ident + `)*\s*(#.*)?$`)
	reNonlocal := regexp.MustCompile(`^` + indent + `nonlocal\s+` + ident + `(\s*,\s*` + ident + `)*\s*(#.*)?$`)

	// Assignment / augmented assignment (permissive)
	reAssign := regexp.MustCompile(`^` + indent + ident + `(\s*,\s*` + ident + `)*\s*([+\-*/%&|^]|<<|>>)?=\s*.+\s*(#.*)?$`)

	// Expression / call line (permissive, but excludes JS/Ruby-ish braces/semicolons)
	reExpr := regexp.MustCompile(`^` + indent + `[^{};]+?\s*(#.*)?$`)

	reBlankOrComment := regexp.MustCompile(commentOrBlank)

	// Ruby-specific tells to reject early (covers your example).
	reRubyEnd := regexp.MustCompile(`^\s*end\s*(#.*)?$`)
	reRubyInstanceVar := regexp.MustCompile(`^\s*@` + ident)
	reRubyPredicateCall := regexp.MustCompile(`\b` + ident + `\?\b`)
	reRubyRange := regexp.MustCompile(`\.\.`)      // e.g. from..to
	reRubySymbol := regexp.MustCompile(`\s:\w+\b`) // e.g. where(id: ...)
	reRubyKeyword := regexp.MustCompile(`^\s*(unless|elsif|then|begin|rescue)\b`)
	reRubyScope := regexp.MustCompile(`::`)
	reRubyBlockArgs := regexp.MustCompile(`\|\s*` + ident)

	// If a line begins with a Python block keyword, it MUST match the corresponding
	// block-line regex UNLESS we're inside a parenthesis/bracket continuation.
	reStartsWithPyBlockKeyword := regexp.MustCompile(`^\s*(def|class|if|elif|else|for|while|try|except|finally|with)\b`)

	// naive single-line string removal for bracket counting
	reSingleQuoted := regexp.MustCompile(`'(?:\\.|[^'\\])*'`)
	reDoubleQuoted := regexp.MustCompile(`"(?:\\.|[^"\\])*"`)

	anyNonBlank := false
	continuationDepth := 0 // open ( [ { minus close ) ] }
	continuedByBackslash := false

	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		anyNonBlank = true

		trim := strings.TrimSpace(line)

		if strings.HasPrefix(trim, "function ") || strings.Contains(trim, "=>") { // JS/TS hints
			return false
		}

		// Ruby rejection heuristics (line-level)
		if reRubyEnd.MatchString(line) ||
			reRubyInstanceVar.MatchString(line) ||
			reRubyPredicateCall.MatchString(line) ||
			reRubyRange.MatchString(line) ||
			reRubySymbol.MatchString(line) ||
			reRubyKeyword.MatchString(line) ||
			reRubyScope.MatchString(line) ||
			reRubyBlockArgs.MatchString(line) {
			return false
		}

		inContinuation := continuationDepth > 0 || continuedByBackslash

		// If it starts with a Python block keyword and we're NOT in a continuation,
		// require a proper Python block line (prevents Ruby's `if cond` form).
		if reStartsWithPyBlockKeyword.MatchString(line) && !inContinuation {
			if !(reDef.MatchString(line) ||
				reClass.MatchString(line) ||
				reIf.MatchString(line) || reElif.MatchString(line) || reElse.MatchString(line) ||
				reFor.MatchString(line) || reWhile.MatchString(line) ||
				reTry.MatchString(line) || reExcept.MatchString(line) || reFinally.MatchString(line) ||
				reWith.MatchString(line) ||
				reWithStart.MatchString(line)) {
				return false
			}
		} else {
			// Non-block-keyword start OR continuation line: validate with broader rules.
			if !(reBlankOrComment.MatchString(line) ||
				reImport.MatchString(line) || reFrom.MatchString(line) ||
				reReturn.MatchString(line) || reRaise.MatchString(line) || reAssert.MatchString(line) ||
				rePass.MatchString(line) || reBreak.MatchString(line) || reContinue.MatchString(line) ||
				reGlobal.MatchString(line) || reNonlocal.MatchString(line) ||
				reAssign.MatchString(line) ||
				reExpr.MatchString(line)) {
				return false
			}
		}

		// Update continuation tracking for next line.
		// Remove simple quoted strings to avoid counting brackets inside them.
		countLine := line
		countLine = reSingleQuoted.ReplaceAllString(countLine, `''`)
		countLine = reDoubleQuoted.ReplaceAllString(countLine, `""`)

		// Drop comments (naively). Good enough for your use-case here.
		if idx := strings.IndexByte(countLine, '#'); idx >= 0 {
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

		trimRight := strings.TrimRight(line, " \t")
		continuedByBackslash = strings.HasSuffix(trimRight, `\`)
	}

	return anyNonBlank
}
