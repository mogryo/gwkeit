package allsnippetspage

import (
	"context"
	"fmt"
	"math"

	"github.com/gdamore/tcell/v2"
	"github.com/gwkeit/cond"
	"github.com/rivo/tview"
)

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
