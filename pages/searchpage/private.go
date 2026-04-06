package searchpage

import (
	"strconv"

	"github.com/gwkeit/slicelib"
	"github.com/gwkeit/utils"
)

func (sp *SearchPage) showSnippet(snippetId int64) {
	snippet := sp.tools.Repo.FindSnippet(sp.tools.Ctx, snippetId)

	sp.title.SetText(snippet.Title, true)
	sp.body.SetText(snippet.Body, true)
	sp.description.SetText(snippet.Description, true)
	sp.urls.SetText(snippet.Url, true)
}

// setResultListPage
/**
Set the result list page to the given page number.
*/
func (sp *SearchPage) setResultListPage(newPage int64) {
	sp.currentPage = newPage
	if len(sp.foundSnippets) == 0 {
		sp.currentPageView.SetText("0 of 0")
	} else {
		sp.currentPageView.SetText(
			strconv.FormatInt(newPage, 10) + " of " + strconv.FormatInt((sp.totalFoundAmount-1)/itemsPerPage+1, 10),
		)
	}
	items := slicelib.SubSlice(sp.foundSnippets, int((sp.currentPage-1)*itemsPerPage), int(itemsPerPage)*int(sp.currentPage))

	sp.resultList.Clear()
	for i, snippet := range items {
		sp.resultList.AddItem(
			snippet.Title,
			strconv.FormatInt(snippet.ID, 10),
			utils.IfElse(i < len(shortcutRunes), shortcutRunes[i], 0),
			nil,
		)
	}
}

// nextResultPage
/**
Move to the next result page if there are more results available.
*/
func (sp *SearchPage) nextResultPage() {
	if sp.totalFoundAmount > itemsPerPage*sp.currentPage {
		sp.setResultListPage(sp.currentPage + 1)
	}
}

// previousResultPage
/**
Move to the previous result page if there is such.
*/
func (sp *SearchPage) previousResultPage() {
	if sp.currentPage > 1 {
		sp.setResultListPage(sp.currentPage - 1)
	}
}
