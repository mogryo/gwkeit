package allsnippetspage

import (
	"context"
	"fmt"
	"math"

	"github.com/gwkeit/uibuilder"
	"github.com/gwkeit/utils"
)

func (asp *AllSnippetsPage) populateTable(ctx context.Context) {
	asp.table.Clear()
	asp.selectedSnippetId = -1
	snippets, err := asp.tools.Repo.FindSnippetsByPage(ctx, asp.currentPage, asp.pageSize)
	if err != nil {
		asp.logs.AddErrorLogs([]string{"error while fetching snippets"})
	}
	asp.snippets = snippets

	asp.table.SetCell(
		0,
		0,
		uibuilder.NewTableCell(asp.themeName, "Title"),
	)
	asp.table.SetCell(
		0,
		1,
		uibuilder.NewTableCell(asp.themeName, "Created At"),
	)
	asp.table.SetCell(
		0,
		2,
		uibuilder.NewTableCell(asp.themeName, "Updated At"),
	)
	for index, snippet := range snippets {
		rowIndex := index + 1
		asp.table.SetCell(
			rowIndex,
			0,
			uibuilder.NewTableCell(asp.themeName, snippet.Title).SetExpansion(1),
		)

		createdTimestamp := utils.IfElse(
			snippet.CreatedAt.Valid,
			snippet.CreatedAt.Time.Format("2006-01-02 15:04:05"),
			"",
		)
		asp.table.SetCell(
			rowIndex,
			1,
			uibuilder.NewTableCell(asp.themeName, createdTimestamp),
		)

		updatedTimestamp := utils.IfElse(
			snippet.UpdatedAt.Valid,
			snippet.UpdatedAt.Time.Format("2006-01-02 15:04:05"),
			"",
		)
		asp.table.SetCell(
			rowIndex,
			2,
			uibuilder.NewTableCell(asp.themeName, updatedTimestamp),
		)
	}

	snippetCount := asp.tools.Repo.GetSnippetCount(asp.tools.Ctx)
	asp.pagesAmount = int64(math.Ceil(float64(snippetCount) / float64(asp.pageSize)))
	asp.totalSnippetAmountView.SetText(fmt.Sprintf("%d", snippetCount))
	asp.totalPagesView.SetText(fmt.Sprintf("%d", asp.pagesAmount))
}

func (asp *AllSnippetsPage) clearMetadataFields() {
	asp.description.SetText("", false)
	asp.title.SetText("", false)
	asp.urls.SetText("", false)
}

func (asp *AllSnippetsPage) focusTable() {
	if asp.table.GetRowCount() > 1 {
		asp.tools.Focus(asp.table)
		asp.table.SetSelectable(true, false)
		asp.table.Select(1, 0)
	}
}
