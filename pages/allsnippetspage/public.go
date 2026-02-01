package allsnippetspage

func (asp *AllSnippetsPage) SwitchToSnippetListPage() {
	asp.globalDeps.Pages.SwitchToPage("AllSnippets")
	asp.populateTable(asp.globalDeps.Ctx)
	asp.globalDeps.App.SetFocus(asp.table)
}
