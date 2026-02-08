package allsnippetspage

import (
	"github.com/gwkeit/apptools"
	"github.com/gwkeit/configuration"
	"github.com/gwkeit/gwkeitdb"
	"github.com/gwkeit/uibuilder"
	"github.com/gwkeit/widgets"
	"github.com/rivo/tview"
)

type AllSnippetsPage struct {
	grid                   *tview.Grid
	Frame                  *tview.Frame
	title                  *tview.TextArea
	description            *tview.TextArea
	urls                   *tview.TextArea
	table                  *tview.Table
	tools                  *apptools.Tools
	logs                   *widgets.LogsWidget
	totalSnippetAmountView *tview.TextView
	totalPagesView         *tview.TextView
	pageSize               int64
	pagesAmount            int64
	currentPage            int64
	currentPageInput       *tview.InputField
	pageSizeInput          *tview.InputField
	snippets               []gwkeitdb.Snippet
	selectedSnippetId      int64
	pageConf               configuration.IAllSnippetsConf
}

const (
	MaxPageSize = 30
)

func NewPage(
	tools *apptools.Tools,
	logs *widgets.LogsWidget,
	pageConf configuration.IAllSnippetsConf,
) *AllSnippetsPage {
	asp := &AllSnippetsPage{
		tools:                  tools,
		title:                  uibuilder.NewTextArea("", ""),
		description:            uibuilder.NewTextArea("", ""),
		urls:                   uibuilder.NewTextArea("", ""),
		totalSnippetAmountView: uibuilder.NewTextView("Total items: 0"),
		totalPagesView:         uibuilder.NewTextView("Total pages: 0"),
		logs:                   logs,
		pageConf:               pageConf,
		pageSize:               pageConf.GetPageSize(),
		selectedSnippetId:      -1,
	}

	snippetDataFlex := asp.initSnippetDataFlex()
	asp.initTable()
	asp.initCurrentPageInput()
	asp.initPageSizeInput()
	tablePaginationFlex := asp.initTablePagination()
	asp.initGridLayout(snippetDataFlex, tablePaginationFlex)
	asp.initFrame()
	asp.initInputCapture()

	return asp
}
