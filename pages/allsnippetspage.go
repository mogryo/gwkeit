package pages

import (
	"context"
	"fmt"
	"math"
	"strconv"

	"github.com/gdamore/tcell/v2"
	"github.com/gwkeit/cond"
	"github.com/gwkeit/globaldeps"
	"github.com/gwkeit/gwkeitdb"
	"github.com/gwkeit/uibuilder"
	"github.com/gwkeit/widgets"
	"github.com/rivo/tview"
)

type AllSnippetsPage struct {
	grid                   *tview.Grid
	frame                  *tview.Frame
	title                  *tview.TextArea
	description            *tview.TextArea
	urls                   *tview.TextArea
	table                  *tview.Table
	globalDeps             *globaldeps.GlobalDependencies
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
}

const (
	MAX_PAGE_SIZE     = 30
	DEFAULT_PAGE_SIZE = 10
)

func NewAllSnippetsPage(globalDeps *globaldeps.GlobalDependencies, logs *widgets.LogsWidget) *AllSnippetsPage {
	asp := &AllSnippetsPage{
		globalDeps:             globalDeps,
		title:                  uibuilder.NewTextArea("", ""),
		description:            uibuilder.NewTextArea("", ""),
		urls:                   uibuilder.NewTextArea("", ""),
		totalSnippetAmountView: uibuilder.NewTextView("Total items: 0"),
		totalPagesView:         uibuilder.NewTextView("Total pages: 0"),
		logs:                   logs,
		pageSize:               DEFAULT_PAGE_SIZE,
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

func (asp *AllSnippetsPage) populateTable(
	ctx context.Context,
	page int64,
	pageSize int64,
) {
	asp.table.Clear()
	asp.selectedSnippetId = -1
	snippets, err := asp.globalDeps.Repo.FindSnippetsByPage(ctx, page, pageSize)
	if err != nil {
		asp.logs.AddErrorLogs([]string{"error while fetching snippets"})
	}
	asp.snippets = snippets

	asp.table.SetCell(
		0,
		0,
		tview.NewTableCell("Title").SetAlign(tview.AlignLeft).SetBackgroundColor(tcell.ColorDefault).SetAttributes(tcell.AttrBold),
	)
	asp.table.SetCell(
		0,
		1,
		tview.NewTableCell("Created At").SetAlign(tview.AlignLeft).SetBackgroundColor(tcell.ColorDefault).SetAttributes(tcell.AttrBold),
	)
	asp.table.SetCell(
		0,
		2,
		tview.NewTableCell("Updated At").SetAlign(tview.AlignLeft).SetBackgroundColor(tcell.ColorDefault).SetAttributes(tcell.AttrBold),
	)
	for index, snippet := range snippets {
		rowIndex := index + 1
		asp.table.SetCell(
			rowIndex,
			0,
			tview.NewTableCell(snippet.Title).
				SetAlign(tview.AlignLeft).
				SetBackgroundColor(tcell.ColorDefault).
				SetExpansion(1),
		)

		createdTimestamp := cond.IfElse(
			snippet.CreatedAt.Valid,
			snippet.CreatedAt.Time.Format("2006-01-02 15:04:05"),
			"",
		)
		asp.table.SetCell(
			rowIndex,
			1,
			tview.NewTableCell(createdTimestamp).SetAlign(tview.AlignLeft).SetBackgroundColor(tcell.ColorDefault),
		)

		updatedTimestamp := cond.IfElse(
			snippet.UpdatedAt.Valid,
			snippet.UpdatedAt.Time.Format("2006-01-02 15:04:05"),
			"",
		)
		asp.table.SetCell(
			rowIndex,
			2,
			tview.NewTableCell(updatedTimestamp).SetAlign(tview.AlignLeft).SetBackgroundColor(tcell.ColorDefault),
		)
	}

	snippetCount := asp.globalDeps.Repo.GetSnippetCount(asp.globalDeps.Ctx)
	asp.pagesAmount = int64(math.Ceil(float64(snippetCount) / float64(pageSize)))
	asp.pageSize = pageSize
	asp.totalSnippetAmountView.SetText(fmt.Sprintf("%d", snippetCount))
	asp.totalPagesView.SetText(fmt.Sprintf("%d", asp.pagesAmount))
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

func (asp *AllSnippetsPage) initSnippetDataFlex() *tview.Flex {
	return tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(uibuilder.NewWidget("Title:", asp.title), 0, 1, false).
		AddItem(uibuilder.NewWidget("Description:", asp.description), 0, 3, false).
		AddItem(uibuilder.NewWidget("URLs:", asp.urls), 0, 3, false)
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

		if value > 0 && value <= MAX_PAGE_SIZE {
			asp.populateTable(asp.globalDeps.Ctx, asp.currentPage, asp.pageSize)
			return true
		}

		return false
	})
	asp.pageSizeInput.SetText(strconv.FormatInt(DEFAULT_PAGE_SIZE, 10))
	asp.pageSize = DEFAULT_PAGE_SIZE
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
	asp.frame = tview.NewFrame(asp.grid).
		SetBorders(0, 0, 0, 0, 0, 0).
		AddText("[::b]All snippets[::-]", true, tview.AlignCenter, tcell.ColorWhite)
	asp.frame.SetBackgroundColor(tcell.ColorDefault)
}

func (asp *AllSnippetsPage) SwitchToSnippetListPage() {
	asp.globalDeps.Pages.SwitchToPage("AllSnippets")
	asp.populateTable(asp.globalDeps.Ctx, asp.currentPage, asp.pageSize)
	asp.globalDeps.App.SetFocus(asp.table)
}
