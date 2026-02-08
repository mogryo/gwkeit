package allsnippetspage

import (
	"strconv"

	"github.com/gdamore/tcell/v2"
	"github.com/gwkeit/apptools"
	"github.com/gwkeit/configuration"
	"github.com/gwkeit/uibuilder"
	"github.com/rivo/tview"
)

var shortcutDescription = []apptools.ShortcutDescription{
	{"ctrl+P", "Focus page number field"},
	{"ctrl+I", "Focus page size field"},
	{"ctrl+T", "Focus table"},
	{"ctrl+E", "Edit selected snippet"},
}

func (asp *AllSnippetsPage) initCurrentPageInput() {
	asp.currentPageInput = uibuilder.NewInputField("", "").SetAcceptanceFunc(func(textToCheck string, lastChar rune) bool {
		isInteger := tview.InputFieldInteger(textToCheck, lastChar)

		if !isInteger {
			return false
		}
		value, err := strconv.ParseInt(textToCheck, 10, 64)
		if err != nil {
			panic(err)
		}

		if value > 0 && value <= asp.pagesAmount {
			asp.currentPage = value
			asp.populateTable(asp.tools.Ctx)
			return true
		}

		return false
	})
	asp.currentPageInput.SetText("1")
	asp.currentPage = 1
}

func (asp *AllSnippetsPage) initPageSizeInput() {
	asp.pageSizeInput = uibuilder.NewInputField("", "").SetAcceptanceFunc(func(textToCheck string, lastChar rune) bool {
		isInteger := tview.InputFieldInteger(textToCheck, lastChar)

		if !isInteger {
			return false
		}

		value, err := strconv.ParseInt(textToCheck, 10, 64)
		if err != nil {
			panic(err)
		}

		if value > 0 && value <= MaxPageSize {
			asp.pageSize = value
			asp.pageConf.SetPageSize(asp.pageSize)
			asp.populateTable(asp.tools.Ctx)
			return true
		}

		return false
	})
	asp.pageSizeInput.SetText(strconv.FormatInt(asp.pageSize, 10))
}

func (asp *AllSnippetsPage) initTablePagination() *tview.Flex {
	tablePaginationFlex := tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(uibuilder.NewWidget("Total items", asp.totalSnippetAmountView), 0, 1, false).
		AddItem(uibuilder.NewWidget("Total pages", asp.totalPagesView), 0, 1, false).
		AddItem(uibuilder.NewWidget("Items per page: ", asp.pageSizeInput), 0, 2, false).
		AddItem(uibuilder.NewWidget("Page: ", asp.currentPageInput), 0, 2, false)
	tablePaginationFlex.SetBorderPadding(0, 0, 0, 0).SetBackgroundColor(tcell.ColorDefault)

	return tablePaginationFlex
}

func (asp *AllSnippetsPage) initGridLayout(snippetDataFlex *tview.Flex, tablePaginationFlex *tview.Flex) {
	asp.grid = tview.NewGrid().
		SetRows(14, 0, 3).
		SetColumns(0, 50).
		SetBorders(false).
		AddItem(uibuilder.NewWidget("Snippets:", asp.table), 0, 0, 2, 1, 0, 100, false).
		AddItem(uibuilder.NewWidget("Logs:", asp.logs.View), 0, 1, 1, 1, 0, 100, false).
		AddItem(snippetDataFlex, 1, 1, 2, 1, 0, 100, false).
		AddItem(tablePaginationFlex, 2, 0, 1, 1, 0, 100, false)
	asp.grid.SetBackgroundColor(tcell.ColorDefault)
}

func (asp *AllSnippetsPage) initFrame() {
	asp.Frame = tview.NewFrame(asp.grid).
		SetBorders(0, 0, 0, 0, 0, 0).
		AddText("[::b]All snippets[::-]", true, tview.AlignCenter, tcell.ColorDefault)
	asp.Frame.SetBackgroundColor(tcell.ColorDefault)
}

func (asp *AllSnippetsPage) initTable() {
	asp.table = tview.NewTable().
		SetBorders(true)
	asp.table.SetBackgroundColor(tcell.ColorDefault)
	asp.table.SetFixed(1, 1).
		SetSelectedFunc(func(row int, column int) {
			if row == 0 {
				return
			}
			snippet := asp.tools.Repo.FindSnippet(asp.tools.Ctx, asp.snippets[row-1].ID)
			asp.selectedSnippetId = snippet.ID
			asp.title.SetText(snippet.Title, false)
			asp.description.SetText(snippet.Description, false)
			asp.urls.SetText(snippet.Url, false)
		}).
		SetBlurFunc(func() {
			asp.table.SetSelectable(false, false)
		})
}

func (asp *AllSnippetsPage) initSnippetDataFlex() *tview.Flex {
	return tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(uibuilder.NewWidget("Title:", asp.title), 0, 1, false).
		AddItem(uibuilder.NewWidget("Description:", asp.description), 0, 3, false).
		AddItem(uibuilder.NewWidget("URLs:", asp.urls), 0, 3, false)
}

func (asp *AllSnippetsPage) initInputCapture() {
	asp.grid.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		resultEvent := event

		if event.Rune() == '?' {
			asp.tools.GoToPage(configuration.ShortcutModal, shortcutDescription)
			resultEvent = nil
		}

		switch event.Key() {
		case tcell.KeyCtrlP:
			asp.tools.Focus(asp.currentPageInput)
			resultEvent = nil
		case tcell.KeyCtrlI:
			asp.tools.Focus(asp.pageSizeInput)
			resultEvent = nil
		case tcell.KeyCtrlT:
			if asp.table.GetRowCount() > 1 {
				asp.tools.Focus(asp.table)
				asp.table.SetSelectable(true, false)
				asp.table.Select(1, 0)
			}
			resultEvent = nil
		case tcell.KeyCtrlE:
			if asp.selectedSnippetId > -1 {
				asp.tools.GoToPage(configuration.EditPage, asp.selectedSnippetId)
			} else {
				asp.logs.AddErrorLogs([]string{"No snippet selected."})
			}
			resultEvent = nil
		}

		return resultEvent
	})
}
