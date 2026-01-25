package editpage

func (ep *EditPage) SwitchToEditPage(snippetId int64) {
	ep.loadSnippet(snippetId)
	ep.globalDeps.Pages.SwitchToPage("Edit")
	ep.globalDeps.App.SetFocus(ep.body)
}
