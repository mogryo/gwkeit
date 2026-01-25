package searchpage

import (
	"context"

	"github.com/gwkeit/globaldeps"
	"github.com/gwkeit/gwkeitdb"
	"github.com/gwkeit/widgets"
	"github.com/rivo/tview"
)

const PageName = "Search"

var (
	shortcutRunes = []rune{'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j'}
)

type SearchPage struct {
	body              *tview.TextArea
	searchField       *tview.InputField
	resultList        *tview.List
	description       *tview.TextArea
	urls              *tview.TextArea
	title             *tview.TextArea
	grid              *tview.Grid
	Frame             *tview.Frame
	searchType        *tview.DropDown
	searchBox         *tview.Flex
	searchCallback    func(ctx context.Context, words []string) []gwkeitdb.Snippet
	selectedSnippetId int64
	globalDeps        *globaldeps.GlobalDependencies
	logs              *widgets.LogsWidget
}

func NewPage(
	globalDeps *globaldeps.GlobalDependencies,
	logs *widgets.LogsWidget,
) *SearchPage {
	searchPage := &SearchPage{
		selectedSnippetId: -1,
		globalDeps:        globalDeps,
		logs:              logs,
	}

	searchPage.initMetadataFields()
	searchPage.initBody()
	searchPage.initResultList()
	searchPage.initSearchField()
	searchPage.initGridLayout()
	searchPage.initInputCapture()
	searchPage.initFrame()

	return searchPage
}
