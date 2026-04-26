package langdetector

import (
	"regexp"
	"strings"
)

func IsPostgreSQL(src string) bool {
	src = strings.ReplaceAll(src, "\r\n", "\n")
	src = strings.ReplaceAll(src, "\r", "\n")

	lines := strings.Split(src, "\n")

	ident := `[A-Za-z_][A-Za-z0-9_]*`
	indent := `[ \t]*`

	reBlankOrComment := regexp.MustCompile(`^\s*(--.*)?$`)
	reBlockCommentStart := regexp.MustCompile(`^\s*/\*`)
	reBlockCommentEnd := regexp.MustCompile(`\*/`)

	// DDL
	reCreateTable := regexp.MustCompile(`(?i)^` + indent + `CREATE\s+(TEMP(ORARY)?\s+)?TABLE\s+(IF\s+NOT\s+EXISTS\s+)?`)
	reCreateIndex := regexp.MustCompile(`(?i)^` + indent + `CREATE\s+(UNIQUE\s+)?INDEX\s+(CONCURRENTLY\s+)?(IF\s+NOT\s+EXISTS\s+)?`)
	reCreateView := regexp.MustCompile(`(?i)^` + indent + `CREATE\s+(OR\s+REPLACE\s+)?(TEMP(ORARY)?\s+)?VIEW\s+`)
	reCreateFunction := regexp.MustCompile(`(?i)^` + indent + `CREATE\s+(OR\s+REPLACE\s+)?FUNCTION\s+`)
	reCreateProcedure := regexp.MustCompile(`(?i)^` + indent + `CREATE\s+(OR\s+REPLACE\s+)?PROCEDURE\s+`)
	reCreateTrigger := regexp.MustCompile(`(?i)^` + indent + `CREATE\s+(OR\s+REPLACE\s+)?(CONSTRAINT\s+)?TRIGGER\s+`)
	reCreateSequence := regexp.MustCompile(`(?i)^` + indent + `CREATE\s+(TEMP(ORARY)?\s+)?SEQUENCE\s+`)
	reCreateSchema := regexp.MustCompile(`(?i)^` + indent + `CREATE\s+SCHEMA\s+`)
	reCreateExtension := regexp.MustCompile(`(?i)^` + indent + `CREATE\s+EXTENSION\s+`)
	reCreateType := regexp.MustCompile(`(?i)^` + indent + `CREATE\s+TYPE\s+`)
	reAlterTable := regexp.MustCompile(`(?i)^` + indent + `ALTER\s+TABLE\s+`)
	reAlterSequence := regexp.MustCompile(`(?i)^` + indent + `ALTER\s+SEQUENCE\s+`)
	reDropTable := regexp.MustCompile(`(?i)^` + indent + `DROP\s+TABLE\s+`)
	reDropIndex := regexp.MustCompile(`(?i)^` + indent + `DROP\s+INDEX\s+`)
	reDropView := regexp.MustCompile(`(?i)^` + indent + `DROP\s+VIEW\s+`)
	reDropFunction := regexp.MustCompile(`(?i)^` + indent + `DROP\s+FUNCTION\s+`)
	reDropSequence := regexp.MustCompile(`(?i)^` + indent + `DROP\s+SEQUENCE\s+`)
	reDropSchema := regexp.MustCompile(`(?i)^` + indent + `DROP\s+SCHEMA\s+`)
	reDropExtension := regexp.MustCompile(`(?i)^` + indent + `DROP\s+EXTENSION\s+`)

	// DML
	reSelect := regexp.MustCompile(`(?i)^` + indent + `SELECT\b`)
	reInsert := regexp.MustCompile(`(?i)^` + indent + `INSERT\s+INTO\b`)
	reUpdate := regexp.MustCompile(`(?i)^` + indent + `UPDATE\s+` + ident)
	reDelete := regexp.MustCompile(`(?i)^` + indent + `DELETE\s+FROM\b`)
	reWith := regexp.MustCompile(`(?i)^` + indent + `WITH\b`)

	// TCL / control
	reBegin := regexp.MustCompile(`(?i)^` + indent + `BEGIN\b`)
	reCommit := regexp.MustCompile(`(?i)^` + indent + `COMMIT\b`)
	reRollback := regexp.MustCompile(`(?i)^` + indent + `ROLLBACK\b`)
	reSavepoint := regexp.MustCompile(`(?i)^` + indent + `SAVEPOINT\b`)

	// PostgreSQL-specific statements
	reSet := regexp.MustCompile(`(?i)^` + indent + `SET\s+(LOCAL\s+|SESSION\s+)?` + ident)
	reShow := regexp.MustCompile(`(?i)^` + indent + `SHOW\b`)
	reGrant := regexp.MustCompile(`(?i)^` + indent + `GRANT\b`)
	reRevoke := regexp.MustCompile(`(?i)^` + indent + `REVOKE\b`)
	reCopyStmt := regexp.MustCompile(`(?i)^` + indent + `COPY\b`)
	reAnalyze := regexp.MustCompile(`(?i)^` + indent + `ANALYZE\b`)
	reExplain := regexp.MustCompile(`(?i)^` + indent + `EXPLAIN\b`)
	reVacuum := regexp.MustCompile(`(?i)^` + indent + `VACUUM\b`)
	reDoBlock := regexp.MustCompile(`(?i)^` + indent + `DO\b`)

	// PG-specific tokens detectable anywhere in a line
	reReturning := regexp.MustCompile(`(?i)\bRETURNING\b`)
	reDollarQuote := regexp.MustCompile(`\$[A-Za-z0-9_]*\$`)
	reLanguageLine := regexp.MustCompile(`(?i)\bLANGUAGE\s+` + ident)
	rePGTypes := regexp.MustCompile(`(?i)\b(SERIAL|BIGSERIAL|SMALLSERIAL|TABLESPACE|INHERITS|EXCLUDE)\b`)

	// SQLite-specific rejects
	reRejectSQLite := regexp.MustCompile(`(?i)\b(AUTOINCREMENT|WITHOUT\s+ROWID)\b`)
	reRejectSQLiteStmt := regexp.MustCompile(`(?i)^\s*(PRAGMA|ATTACH|DETACH|REINDEX)\b`)

	// Reject obvious non-SQL languages
	reRejectOther := regexp.MustCompile(`^\s*(func|def|class|val|var|fun|import|export|package)\b|^\s*[{]?\s*(=>|->)`)

	anyNonBlank := false
	inBlockComment := false
	hasSQLStatement := false
	hasPGSpecific := false

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

		// Reject non-SQL languages
		if reRejectOther.MatchString(trim) {
			return false
		}

		// Reject SQLite-specific
		if reRejectSQLite.MatchString(line) || reRejectSQLiteStmt.MatchString(line) {
			return false
		}

		// Detect PG-specific signals anywhere in the line
		if reReturning.MatchString(trim) || reDollarQuote.MatchString(trim) ||
			reLanguageLine.MatchString(trim) || rePGTypes.MatchString(trim) ||
			reCreateFunction.MatchString(trim) || reCreateProcedure.MatchString(trim) ||
			reCreateSequence.MatchString(trim) || reCreateSchema.MatchString(trim) ||
			reCreateExtension.MatchString(trim) || reCreateType.MatchString(trim) ||
			reAlterSequence.MatchString(trim) || reDropFunction.MatchString(trim) ||
			reDropSequence.MatchString(trim) || reDropSchema.MatchString(trim) ||
			reDropExtension.MatchString(trim) || reGrant.MatchString(trim) ||
			reRevoke.MatchString(trim) || reCopyStmt.MatchString(trim) ||
			reDoBlock.MatchString(trim) {
			hasPGSpecific = true
		}

		isSQLStatement := reCreateTable.MatchString(trim) || reCreateIndex.MatchString(trim) ||
			reCreateView.MatchString(trim) || reCreateFunction.MatchString(trim) ||
			reCreateProcedure.MatchString(trim) || reCreateTrigger.MatchString(trim) ||
			reCreateSequence.MatchString(trim) || reCreateSchema.MatchString(trim) ||
			reCreateExtension.MatchString(trim) || reCreateType.MatchString(trim) ||
			reAlterTable.MatchString(trim) || reAlterSequence.MatchString(trim) ||
			reDropTable.MatchString(trim) || reDropIndex.MatchString(trim) ||
			reDropView.MatchString(trim) || reDropFunction.MatchString(trim) ||
			reDropSequence.MatchString(trim) || reDropSchema.MatchString(trim) ||
			reDropExtension.MatchString(trim) ||
			reSelect.MatchString(trim) || reInsert.MatchString(trim) ||
			reUpdate.MatchString(trim) || reDelete.MatchString(trim) ||
			reWith.MatchString(trim) ||
			reBegin.MatchString(trim) || reCommit.MatchString(trim) ||
			reRollback.MatchString(trim) || reSavepoint.MatchString(trim) ||
			reSet.MatchString(trim) || reShow.MatchString(trim) ||
			reGrant.MatchString(trim) || reRevoke.MatchString(trim) ||
			reCopyStmt.MatchString(trim) || reAnalyze.MatchString(trim) ||
			reExplain.MatchString(trim) || reVacuum.MatchString(trim) ||
			reDoBlock.MatchString(trim) || reReturning.MatchString(trim) ||
			reDollarQuote.MatchString(trim) || reLanguageLine.MatchString(trim)

		if isSQLStatement {
			hasSQLStatement = true
			// A SELECT/INSERT/UPDATE/DELETE with no dialect-specific tokens is still
			// counted as PG-specific when no SQLite tokens are present.
			if reSelect.MatchString(trim) || reInsert.MatchString(trim) ||
				reUpdate.MatchString(trim) || reDelete.MatchString(trim) {
				hasPGSpecific = true
			}
		}
	}

	if !anyNonBlank || !hasSQLStatement {
		return false
	}

	return hasPGSpecific
}
