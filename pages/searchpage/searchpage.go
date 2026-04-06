package searchpage

import (
	"context"

	"github.com/gwkeit/apptools"
	"github.com/gwkeit/configuration"
	"github.com/gwkeit/gwkeitdb"
	"github.com/gwkeit/uibuilder"
	"github.com/gwkeit/widgets"
	"github.com/rivo/tview"
)

var (
	shortcutRunes = []rune{'1', '2', '3', '4', '5', '6', '7', '8', '9'}
	itemsPerPage  = int64(len(shortcutRunes))
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
	paginationBox     *tview.Flex
	totalFoundView    *tview.TextView
	currentPageView   *tview.TextView
	searchCallback    func(ctx context.Context, words []string) []gwkeitdb.Snippet
	selectedSnippetId int64
	tools             *apptools.Tools
	logs              *widgets.LogsWidget
	pageConf          configuration.ISearchPageConf
	themeName         uibuilder.AppThemeName
	foundSnippets     []gwkeitdb.Snippet
	totalFoundAmount  int64
	currentPage       int64
}

func NewPage(
	tools *apptools.Tools,
	pageState configuration.ISearchPageConf,
	logs *widgets.LogsWidget,
	themeName uibuilder.AppThemeName,
) *SearchPage {
	searchPage := &SearchPage{
		selectedSnippetId: -1,
		tools:             tools,
		logs:              logs,
		pageConf:          pageState,
		themeName:         themeName,
		totalFoundAmount:  0,
		currentPage:       1,
		foundSnippets:     []gwkeitdb.Snippet{},
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
