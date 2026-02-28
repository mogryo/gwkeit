package configuration

type Language string

const (
	Kotlin     Language = "Kotlin"
	Python     Language = "Python"
	TypeScript Language = "TypeScript"
	Ruby       Language = "Ruby"
	Go         Language = "Go"
	Text       Language = "Text"
)

var LanguagesStrings = []string{Text.String(), Go.String(), Kotlin.String(), Python.String(), Ruby.String(), TypeScript.String()}

func (l Language) String() string { return string(l) }
