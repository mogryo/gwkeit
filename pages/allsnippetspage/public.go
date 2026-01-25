package allsnippetspage

func (asp *AllSnippetsPage) SwitchToSnippetListPage() {
	asp.globalDeps.Pages.SwitchToPage("AllSnippets")
	asp.populateTable(asp.globalDeps.Ctx, asp.currentPage, asp.pageSize)
	asp.globalDeps.App.SetFocus(asp.table)
}
