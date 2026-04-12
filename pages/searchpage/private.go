package searchpage

import (
	"strconv"

	"github.com/gwkeit/configuration"
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
	startIdx := int((sp.currentPage - 1) * itemsPerPage)
	endIdx := int(sp.currentPage * itemsPerPage)
	items := slicelib.SubSlice(sp.foundSnippets, startIdx, endIdx)

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

// showNextResultPage
/**
Move to the next result page if there are more results available.
*/
func (sp *SearchPage) showNextResultPage() {
	if sp.totalFoundAmount > itemsPerPage*sp.currentPage {
		sp.setResultListPage(sp.currentPage + 1)
	}
}

// showPreviousResultPage
/**
Move to the previous result page if there is such.
*/
func (sp *SearchPage) showPreviousResultPage() {
	if sp.currentPage > 1 {
		sp.setResultListPage(sp.currentPage - 1)
	}
}

// clearMetadataFields
/**
Clears the metadata fields.
*/
func (sp *SearchPage) clearMetadataFields() {
	sp.title.SetText("", true)
	sp.body.SetText("", true)
	sp.description.SetText("", true)
	sp.urls.SetText("", true)
	sp.language.SetCurrentOption(0)
}

func (sp *SearchPage) goToEditPage() {
	if sp.selectedSnippetId > -1 {
		sp.tools.GoToPage(configuration.EditPage, sp.selectedSnippetId)
	} else {
		sp.logs.AddErrorLogs([]string{"No snippet selected."})
	}
}
