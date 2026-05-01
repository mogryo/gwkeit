package langdetector

import (
	"regexp"
	"strings"
)

func IsSQLite(src string) bool {
	src = strings.ReplaceAll(src, "\r\n", "\n")
	src = strings.ReplaceAll(src, "\r", "\n")

	lines := strings.Split(src, "\n")

	ident := `[A-Za-z_][A-Za-z0-9_]*`
	indent := `[ \t]*`

	reBlankOrComment := regexp.MustCompile(`^\s*(--.*)?$`)
	reBlockCommentStart := regexp.MustCompile(`^\s*/\*`)
	reBlockCommentEnd := regexp.MustCompile(`\*/`)

	// DDL
	reCreateTable := regexp.MustCompile(`(?i)^` + indent + `CREATE\s+(TEMP(ORARY)?\s+)?TABLE\s+(IF\s+NOT\s+EXISTS\s+)?` + ident)
	reCreateIndex := regexp.MustCompile(`(?i)^` + indent + `CREATE\s+(UNIQUE\s+)?INDEX\s+(IF\s+NOT\s+EXISTS\s+)?` + ident)
	reCreateView := regexp.MustCompile(`(?i)^` + indent + `CREATE\s+(TEMP(ORARY)?\s+)?VIEW\s+(IF\s+NOT\s+EXISTS\s+)?` + ident)
	reCreateTrigger := regexp.MustCompile(`(?i)^` + indent + `CREATE\s+(TEMP(ORARY)?\s+)?TRIGGER\s+(IF\s+NOT\s+EXISTS\s+)?` + ident)
	reAlterTable := regexp.MustCompile(`(?i)^` + indent + `ALTER\s+TABLE\s+` + ident)
	reDropTable := regexp.MustCompile(`(?i)^` + indent + `DROP\s+TABLE\s+(IF\s+EXISTS\s+)?` + ident)
	reDropIndex := regexp.MustCompile(`(?i)^` + indent + `DROP\s+INDEX\s+(IF\s+EXISTS\s+)?` + ident)
	reDropView := regexp.MustCompile(`(?i)^` + indent + `DROP\s+VIEW\s+(IF\s+EXISTS\s+)?` + ident)

	// DML
	reSelect := regexp.MustCompile(`(?i)^` + indent + `SELECT\b`)
	reInsert := regexp.MustCompile(`(?i)^` + indent + `INSERT\s+(OR\s+(REPLACE|IGNORE|ABORT|FAIL|ROLLBACK)\s+)?INTO\b`)
	reUpdate := regexp.MustCompile(`(?i)^` + indent + `UPDATE\s+(OR\s+(REPLACE|IGNORE|ABORT|FAIL|ROLLBACK)\s+)?` + ident)
	reDelete := regexp.MustCompile(`(?i)^` + indent + `DELETE\s+FROM\b`)
	reReplace := regexp.MustCompile(`(?i)^` + indent + `REPLACE\s+INTO\b`)
	reWith := regexp.MustCompile(`(?i)^` + indent + `WITH\b`)

	// TCL
	reBegin := regexp.MustCompile(`(?i)^` + indent + `BEGIN\b`)
	reCommit := regexp.MustCompile(`(?i)^` + indent + `COMMIT\b`)
	reRollback := regexp.MustCompile(`(?i)^` + indent + `ROLLBACK\b`)
	reSavepoint := regexp.MustCompile(`(?i)^` + indent + `SAVEPOINT\b`)

	// SQLite-specific
	rePragma := regexp.MustCompile(`(?i)^` + indent + `PRAGMA\b`)
	reAttach := regexp.MustCompile(`(?i)^` + indent + `ATTACH\s+(DATABASE\s+)?\S+\s+AS\b`)
	reDetach := regexp.MustCompile(`(?i)^` + indent + `DETACH\s+(DATABASE\s+)?\S+`)
	reVacuum := regexp.MustCompile(`(?i)^` + indent + `VACUUM\b`)
	reAnalyze := regexp.MustCompile(`(?i)^` + indent + `ANALYZE\b`)
	reReindex := regexp.MustCompile(`(?i)^` + indent + `REINDEX\b`)
	reAutoincrement := regexp.MustCompile(`(?i)\bAUTOINCREMENT\b`)
	reWithoutRowid := regexp.MustCompile(`(?i)\bWITHOUT\s+ROWID\b`)

	// Inline SQL continuation / expression lines (column definitions, WHERE clauses, etc.)
	reInline := regexp.MustCompile(`(?i)^` + indent + `(FROM|WHERE|JOIN|LEFT\s+JOIN|INNER\s+JOIN|OUTER\s+JOIN|ON|SET|VALUES|GROUP\s+BY|ORDER\s+BY|HAVING|LIMIT|OFFSET|UNION|INTERSECT|EXCEPT|AND|OR|NOT|AS|CASE|WHEN|THEN|ELSE|END|IS|IN|LIKE|BETWEEN|NULL|PRIMARY\s+KEY|FOREIGN\s+KEY|REFERENCES|UNIQUE|CHECK|DEFAULT|NOT\s+NULL|AUTOINCREMENT|WITHOUT\s+ROWID)\b`)

	anyNonBlank := false
	inBlockComment := false
	hasSQLStatement := false
	hasSQLiteSpecific := false

	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		anyNonBlank = true

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

		if reBlankOrComment.MatchString(line) {
			continue
		}

		trim := strings.TrimSpace(line)

		// Check for SQLite-specific constructs
		if rePragma.MatchString(trim) ||
			reAttach.MatchString(trim) ||
			reDetach.MatchString(trim) ||
			reAutoincrement.MatchString(trim) ||
			reWithoutRowid.MatchString(trim) {
			hasSQLiteSpecific = true
		}

		isSQLStatement := reCreateTable.MatchString(trim) || reCreateIndex.MatchString(trim) ||
			reCreateView.MatchString(trim) || reCreateTrigger.MatchString(trim) ||
			reAlterTable.MatchString(trim) ||
			reDropTable.MatchString(trim) || reDropIndex.MatchString(trim) || reDropView.MatchString(trim) ||
			reSelect.MatchString(trim) || reInsert.MatchString(trim) ||
			reUpdate.MatchString(trim) || reDelete.MatchString(trim) ||
			reReplace.MatchString(trim) || reWith.MatchString(trim) ||
			reBegin.MatchString(trim) || reCommit.MatchString(trim) ||
			reRollback.MatchString(trim) || reSavepoint.MatchString(trim) ||
			rePragma.MatchString(trim) || reAttach.MatchString(trim) ||
			reDetach.MatchString(trim) || reVacuum.MatchString(trim) ||
			reAnalyze.MatchString(trim) || reReindex.MatchString(trim) ||
			reInline.MatchString(trim)

		if isSQLStatement {
			hasSQLStatement = true
			if reSelect.MatchString(trim) || reInsert.MatchString(trim) ||
				reUpdate.MatchString(trim) || reDelete.MatchString(trim) {
				hasSQLiteSpecific = true
			}
		}
	}

	if !anyNonBlank || !hasSQLStatement {
		return false
	}

	return hasSQLiteSpecific
}
