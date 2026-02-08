package configuration

type PageName string

const (
	AdditionPage    PageName = "Addition"
	AllSnippetsPage PageName = "AllSnippets"
	EditPage        PageName = "Edit"
	SearchPage      PageName = "Search"
	ShortcutModal   PageName = "Modal"
)

func (p PageName) String() string { return string(p) }
