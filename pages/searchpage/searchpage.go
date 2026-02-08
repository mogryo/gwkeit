package searchpage

import (
	"context"

	"github.com/gwkeit/apptools"
	"github.com/gwkeit/configuration"
	"github.com/gwkeit/gwkeitdb"
	"github.com/gwkeit/widgets"
	"github.com/rivo/tview"
)

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
	tools             *apptools.Tools
	logs              *widgets.LogsWidget
	pageConf          configuration.ISearchPageConf
}

func NewPage(
	tools *apptools.Tools,
	pageState configuration.ISearchPageConf,
	logs *widgets.LogsWidget,
) *SearchPage {
	searchPage := &SearchPage{
		selectedSnippetId: -1,
		tools:             tools,
		logs:              logs,
		pageConf:          pageState,
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
