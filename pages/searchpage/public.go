package searchpage

func (sp *SearchPage) FocusSearchField() {
	sp.globalDeps.App.SetFocus(sp.searchField)
}

func (sp *SearchPage) SwitchToSearchPage() {
	if sp.selectedSnippetId > -1 {
		sp.showSnippet(sp.selectedSnippetId)
	}
	sp.globalDeps.Pages.SwitchToPage(PageName)
	sp.globalDeps.App.SetFocus(sp.searchField)
}
