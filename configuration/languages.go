package configuration

type Language string

const (
	Text       Language = "Text"
	Go         Language = "Go"
	Kotlin     Language = "Kotlin"
	PostgreSQL Language = "PostgreSQL"
	Python     Language = "Python"
	Ruby       Language = "Ruby"
	SQLite     Language = "SQLite"
	TypeScript Language = "TypeScript"
)

var LanguagesStrings = []string{
	Text.String(),
	Go.String(),
	Kotlin.String(),
	PostgreSQL.String(),
	Python.String(),
	Ruby.String(),
	SQLite.String(),
	TypeScript.String(),
}

func (l Language) String() string { return string(l) }
