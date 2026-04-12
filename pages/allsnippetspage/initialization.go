package allsnippetspage

import (
	"slices"
	"strconv"

	"github.com/gdamore/tcell/v2"
	"github.com/gwkeit/apptools"
	"github.com/gwkeit/configuration"
	"github.com/gwkeit/slicelib"
	"github.com/gwkeit/uibuilder"
	"github.com/rivo/tview"
)

var shortcutDescription = []apptools.ShortcutDescription{
	{"ctrl+P", "Focus page number field"},
	{"ctrl+I", "Focus page size field"},
	{"ctrl+T", "Focus table"},
	{"ctrl+E/Enter", "Edit selected snippet"},
	{"}", "Show next table page"},
	{"{", "Show previous table page"},
}

func (asp *AllSnippetsPage) initCurrentPageInput() {
	asp.currentPageInput = uibuilder.NewInputField(asp.themeName, "", "").SetAcceptanceFunc(func(textToCheck string, lastChar rune) bool {
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
	asp.pageSizeInput = uibuilder.NewInputField(asp.themeName, "", "").SetAcceptanceFunc(func(textToCheck string, lastChar rune) bool {
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
		AddItem(uibuilder.NewWidget(asp.themeName, "Total items", asp.totalSnippetAmountView), 0, 1, false).
		AddItem(uibuilder.NewWidget(asp.themeName, "Total pages", asp.totalPagesView), 0, 1, false).
		AddItem(uibuilder.NewWidget(asp.themeName, "Items per page: ", asp.pageSizeInput), 0, 2, false).
		AddItem(uibuilder.NewWidget(asp.themeName, "Page: ", asp.currentPageInput), 0, 2, false)
	tablePaginationFlex.SetBorderPadding(0, 0, 0, 0).SetBackgroundColor(tcell.ColorDefault)

	return tablePaginationFlex
}

func (asp *AllSnippetsPage) initGridLayout(snippetDataFlex *tview.Flex, tablePaginationFlex *tview.Flex) {
	asp.grid = tview.NewGrid().
		SetRows(14, 0, 3).
		SetColumns(0, 50).
		SetBorders(false).
		AddItem(uibuilder.NewWidget(asp.themeName, "Snippets:", asp.table), 0, 0, 2, 1, 0, 100, false).
		AddItem(uibuilder.NewWidget(asp.themeName, "Logs:", asp.logs.View), 0, 1, 1, 1, 0, 100, false).
		AddItem(snippetDataFlex, 1, 1, 2, 1, 0, 100, false).
		AddItem(tablePaginationFlex, 2, 0, 1, 1, 0, 100, false)
	asp.grid.SetBackgroundColor(tcell.ColorDefault)
}

func (asp *AllSnippetsPage) initFrame() {
	asp.Frame = uibuilder.NewPageFrame(asp.themeName, asp.grid, "All snippets")
}

func (asp *AllSnippetsPage) initTable() {
	asp.table = uibuilder.NewTable(asp.themeName, 1, 1)
	asp.table.SetSelectionChangedFunc(func(row, column int) {
		if row == 0 {
			asp.selectedSnippetId = -1
			asp.clearMetadataFields()
			return
		}
		snippet := asp.tools.Repo.FindSnippet(asp.tools.Ctx, asp.snippets[row-1].ID)
		asp.selectedSnippetId = snippet.ID
		asp.title.SetText(snippet.Title, false)
		asp.description.SetText(snippet.Description, false)
		asp.urls.SetText(snippet.Url, false)
		if snippet.Language.Valid {
			idx := slices.Index(configuration.LanguagesStrings, snippet.Language.String)
			asp.language.SetCurrentOption(idx + 1)
		} else {
			asp.language.SetCurrentOption(0)
		}
	})
	asp.table.SetSelectedFunc(func(row int, column int) {
		if asp.selectedSnippetId > -1 {
			asp.tools.GoToPage(configuration.EditPage, asp.selectedSnippetId)
		} else {
			asp.logs.AddErrorLogs([]string{"No snippet selected."})
		}
	}).
		SetBlurFunc(func() {
			asp.table.SetSelectable(false, false)
		})
}

func (asp *AllSnippetsPage) initSnippetDataFlex() *tview.Flex {
	asp.language = uibuilder.NewDropDown(asp.themeName, "")
	asp.language.SetOptions(slicelib.Concat([]string{""}, configuration.LanguagesStrings), nil)
	asp.language.SetCurrentOption(0)
	asp.language.SetDisabled(true)

	return tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(uibuilder.NewWidget(asp.themeName, "Title:", asp.title), 0, 1, false).
		AddItem(uibuilder.NewWidget(asp.themeName, "Description:", asp.description), 0, 3, false).
		AddItem(uibuilder.NewWidget(asp.themeName, "URLs:", asp.urls), 0, 2, false).
		AddItem(uibuilder.NewWidget(asp.themeName, "Language:", asp.language), 3, 1, false)
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
			asp.focusTable()
			resultEvent = nil
		case tcell.KeyCtrlE:
			if asp.selectedSnippetId > -1 {
				asp.tools.GoToPage(configuration.EditPage, asp.selectedSnippetId)
			} else {
				asp.logs.AddErrorLogs([]string{"No snippet selected."})
			}
			resultEvent = nil
		}

		if event.Rune() == '}' {
			asp.showNextTablePage()
			resultEvent = nil
		}
		if event.Rune() == '{' {
			asp.showPreviousTablePage()
			resultEvent = nil
		}

		return resultEvent
	})
}
