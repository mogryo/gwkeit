package langdetector

import (
	"regexp"
	"strings"
)

func IsRuby(src string) bool {
	src = strings.ReplaceAll(src, "\r\n", "\n")
	src = strings.ReplaceAll(src, "\r", "\n")

	lines := strings.Split(src, "\n")
	if len(lines) == 0 {
		return false
	}

	ident := `[A-Za-z_][A-Za-z0-9_]*`
	constIdent := `[A-Z][A-Za-z0-9_]*`
	indent := `[ \t]*`

	reBlankOrComment := regexp.MustCompile(`^\s*(#.*)?$`)

	// Common Ruby constructs
	reDef := regexp.MustCompile(`^` + indent + `def\s+` + ident + `([!?=])?(\s*\(.*\))?(\s+.*)?\s*(#.*)?$`)
	reClass := regexp.MustCompile(`^` + indent + `class\s+` + constIdent + `(\s*<\s*` + constIdent + `(\:\:` + constIdent + `)*)?\s*(#.*)?$`)
	reModule := regexp.MustCompile(`^` + indent + `module\s+` + constIdent + `(\:\:` + constIdent + `)*\s*(#.*)?$`)

	// Control flow / blocks (Ruby doesn't require ":" at end; relies on `end`)
	reIf := regexp.MustCompile(`^` + indent + `if\b\s+.+\s*(#.*)?$`)
	reElsif := regexp.MustCompile(`^` + indent + `elsif\b\s+.+\s*(#.*)?$`)
	reElse := regexp.MustCompile(`^` + indent + `else\b\s*(#.*)?$`)
	reUnless := regexp.MustCompile(`^` + indent + `unless\b\s+.+\s*(#.*)?$`)
	reWhile := regexp.MustCompile(`^` + indent + `while\b\s+.+\s*(#.*)?$`)
	reUntil := regexp.MustCompile(`^` + indent + `until\b\s+.+\s*(#.*)?$`)
	reFor := regexp.MustCompile(`^` + indent + `for\b\s+` + ident + `\s+in\s+.+\s*(#.*)?$`)
	reCase := regexp.MustCompile(`^` + indent + `case(\s+.+)?\s*(#.*)?$`)
	reWhen := regexp.MustCompile(`^` + indent + `when\b\s+.+\s*(#.*)?$`)

	reBegin := regexp.MustCompile(`^` + indent + `begin\b\s*(#.*)?$`)
	reRescue := regexp.MustCompile(`^` + indent + `rescue(\s+.+)?\s*(#.*)?$`)
	reEnsure := regexp.MustCompile(`^` + indent + `ensure\b\s*(#.*)?$`)
	reEnd := regexp.MustCompile(`^` + indent + `end\b\s*(#.*)?$`)

	// Common statements
	reRequire := regexp.MustCompile(`^` + indent + `require(_relative)?\s+.+\s*(#.*)?$`)
	reInclude := regexp.MustCompile(`^` + indent + `include\s+` + constIdent + `(\:\:` + constIdent + `)*\s*(#.*)?$`)
	reExtend := regexp.MustCompile(`^` + indent + `extend\s+` + constIdent + `(\:\:` + constIdent + `)*\s*(#.*)?$`)

	reReturn := regexp.MustCompile(`^` + indent + `return(\s+.+)?\s*(#.*)?$`)
	reRaise := regexp.MustCompile(`^` + indent + `raise(\s+.+)?\s*(#.*)?$`)
	reNext := regexp.MustCompile(`^` + indent + `next(\s+.+)?\s*(#.*)?$`)
	reBreak := regexp.MustCompile(`^` + indent + `break(\s+.+)?\s*(#.*)?$`)
	reRedo := regexp.MustCompile(`^` + indent + `redo\b\s*(#.*)?$`)
	reRetry := regexp.MustCompile(`^` + indent + `retry\b\s*(#.*)?$`)

	// Variables / assignments
	reInstanceVarAssign := regexp.MustCompile(`^` + indent + `@` + ident + `\s*=\s*.+\s*(#.*)?$`)
	reClassVarAssign := regexp.MustCompile(`^` + indent + `@@` + ident + `\s*=\s*.+\s*(#.*)?$`)
	reGlobalVarAssign := regexp.MustCompile(`^` + indent + `\$\w+\s*=\s*.+\s*(#.*)?$`)
	reAssign := regexp.MustCompile(`^` + indent + ident + `(\s*,\s*` + ident + `)*\s*=\s*.+\s*(#.*)?$`)

	// Hash label syntax and symbols (strong Ruby signals)
	reHashLabel := regexp.MustCompile(`\b` + ident + `:\s*`)
	reSymbol := regexp.MustCompile(`(^|[\s\(\[\{,])\:\w+`)
	reRange := regexp.MustCompile(`\.\.`)                                // from..to
	rePredicateCall := regexp.MustCompile(`\b` + ident + `\?(\b|\s|\()`) // foo?

	// Block syntax: do/end, or { |x| ... }
	reDo := regexp.MustCompile(`\bdo\b\s*(\|[^|]*\|)?\s*(#.*)?$`)
	reBlockArgs := regexp.MustCompile(`\{\s*\|[^|]+\|\s*`)

	// General "expression-ish" line (permissive)
	reExpr := regexp.MustCompile(`^` + indent + `.+\s*(#.*)?$`)

	// Quick rejects for other languages to reduce false positives.
	reRejectGo := regexp.MustCompile(`^\s*package\b|\bfunc\b|:=`)
	reRejectPython := regexp.MustCompile(`^\s*(def\s+` + ident + `\s*\(.*\)\s*:\s*$|class\s+` + ident + `.*:\s*$|elif\b|except\b|with\b)|#.*$`)
	reRejectKotlin := regexp.MustCompile(`^\s*(fun|val|var|when)\b`)
	reRejectTS := regexp.MustCompile(`^\s*(import|export|interface|type)\b|=>`)

	// naive single-line string removal for bracket counting
	reSingleQuoted := regexp.MustCompile(`'(?:\\.|[^'\\])*'`)
	reDoubleQuoted := regexp.MustCompile(`"(?:\\.|[^"\\])*"`)

	anyNonBlank := false
	continuationDepth := 0

	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		anyNonBlank = true

		trim := strings.TrimSpace(line)

		// Reject obvious other languages early (so reExpr doesn't accept everything).
		if reRejectGo.MatchString(trim) || reRejectKotlin.MatchString(trim) || reRejectTS.MatchString(trim) || reRejectPython.MatchString(trim) {
			return false
		}

		inContinuation := continuationDepth > 0

		// If not in continuation, require some Ruby-ish match.
		// If in continuation, allow broad expression lines.
		if !inContinuation {
			if reBlankOrComment.MatchString(line) ||
				reDef.MatchString(line) ||
				reClass.MatchString(line) ||
				reModule.MatchString(line) ||
				reIf.MatchString(line) || reElsif.MatchString(line) || reElse.MatchString(line) || reUnless.MatchString(line) ||
				reWhile.MatchString(line) || reUntil.MatchString(line) || reFor.MatchString(line) ||
				reCase.MatchString(line) || reWhen.MatchString(line) ||
				reBegin.MatchString(line) || reRescue.MatchString(line) || reEnsure.MatchString(line) || reEnd.MatchString(line) ||
				reRequire.MatchString(line) || reInclude.MatchString(line) || reExtend.MatchString(line) ||
				reReturn.MatchString(line) || reRaise.MatchString(line) ||
				reNext.MatchString(line) || reBreak.MatchString(line) || reRedo.MatchString(line) || reRetry.MatchString(line) ||
				reInstanceVarAssign.MatchString(line) || reClassVarAssign.MatchString(line) || reGlobalVarAssign.MatchString(line) ||
				reAssign.MatchString(line) ||
				reExpr.MatchString(line) {
				// ok
			} else {
				return false
			}
		} else {
			if !(reBlankOrComment.MatchString(line) || reExpr.MatchString(line)) {
				return false
			}
		}

		// Strong Ruby signals. If none ever appear, we conservatively return false
		// at the end to avoid "any text is Ruby".
		_ = reHashLabel
		_ = reSymbol
		_ = reRange
		_ = rePredicateCall
		_ = reDo
		_ = reBlockArgs

		// Update continuationDepth for next line by counting (), [], {} ignoring strings and # comments.
		countLine := line
		countLine = reSingleQuoted.ReplaceAllString(countLine, `''`)
		countLine = reDoubleQuoted.ReplaceAllString(countLine, `""`)

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
	}

	if !anyNonBlank {
		return false
	}

	// Require at least one "Ruby-ish" token somewhere to avoid excessive false positives.
	for _, line := range lines {
		l := strings.TrimSpace(line)
		if l == "" || strings.HasPrefix(l, "#") {
			continue
		}
		if reDef.MatchString(l) ||
			reClass.MatchString(l) ||
			reModule.MatchString(l) ||
			reEnd.MatchString(l) ||
			reInstanceVarAssign.MatchString(l) ||
			reHashLabel.MatchString(l) ||
			reSymbol.MatchString(l) ||
			reRange.MatchString(l) ||
			rePredicateCall.MatchString(l) ||
			reDo.MatchString(l) ||
			reBlockArgs.MatchString(l) ||
			strings.Contains(l, "::") {
			return true
		}
	}

	return false
}
