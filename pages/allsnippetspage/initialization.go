package allsnippetspage

import (
	"strconv"

	"github.com/gdamore/tcell/v2"
	"github.com/gwkeit/uibuilder"
	"github.com/rivo/tview"
)

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
			asp.populateTable(asp.globalDeps.Ctx, asp.currentPage, asp.pageSize)
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
			asp.populateTable(asp.globalDeps.Ctx, asp.currentPage, asp.pageSize)
			return true
		}

		return false
	})
	asp.pageSizeInput.SetText(strconv.FormatInt(DefaultPageSize, 10))
	asp.pageSize = DefaultPageSize
}

func (asp *AllSnippetsPage) initTablePagination() *tview.Flex {
	tablePaginationFlex := tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(uibuilder.NewWidget("Total items", asp.totalSnippetAmountView), 0, 1, false).
		AddItem(uibuilder.NewWidget("Total pages", asp.totalPagesView), 0, 1, false).
		AddItem(uibuilder.NewWidget("[ctrl+i] Items per page: ", asp.pageSizeInput), 0, 2, false).
		AddItem(uibuilder.NewWidget("[ctrl+p] Page: ", asp.currentPageInput), 0, 2, false)
	tablePaginationFlex.SetBorderPadding(0, 0, 0, 0).SetBackgroundColor(tcell.ColorDefault)

	return tablePaginationFlex
}

func (asp *AllSnippetsPage) initGridLayout(snippetDataFlex *tview.Flex, tablePaginationFlex *tview.Flex) {
	asp.grid = tview.NewGrid().
		SetRows(14, 0, 3).
		SetColumns(0, 50).
		SetBorders(false).
		AddItem(uibuilder.NewWidget("[ctrl+t] Snippets:", asp.table), 0, 0, 2, 1, 0, 100, false).
		AddItem(uibuilder.NewWidget("Logs:", asp.logs.View), 0, 1, 1, 1, 0, 100, false).
		AddItem(snippetDataFlex, 1, 1, 2, 1, 0, 100, false).
		AddItem(tablePaginationFlex, 2, 0, 1, 1, 0, 100, false)
	asp.grid.SetBackgroundColor(tcell.ColorDefault)
}

func (asp *AllSnippetsPage) initFrame() {
	asp.Frame = tview.NewFrame(asp.grid).
		SetBorders(0, 0, 0, 0, 0, 0).
		AddText("[::b]All snippets[::-]", true, tview.AlignCenter, tcell.ColorWhite)
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
			snippet := asp.globalDeps.Repo.FindSnippet(asp.globalDeps.Ctx, asp.snippets[row-1].ID)
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

		switch event.Key() {
		case tcell.KeyCtrlP:
			asp.globalDeps.App.SetFocus(asp.currentPageInput)
			resultEvent = nil
		case tcell.KeyCtrlI:
			asp.globalDeps.App.SetFocus(asp.pageSizeInput)
			resultEvent = nil
		case tcell.KeyCtrlT:
			if asp.table.GetRowCount() > 1 {
				asp.globalDeps.App.SetFocus(asp.table)
				asp.table.SetSelectable(true, false)
				asp.table.Select(1, 0)
			}
			resultEvent = nil
		case tcell.KeyCtrlE:
			if asp.selectedSnippetId > -1 {
				asp.globalDeps.GoToEditPage(asp.selectedSnippetId)
			} else {
				asp.logs.AddErrorLogs([]string{"No snippet selected."})
			}
			resultEvent = nil
		}

		return resultEvent
	})
}
