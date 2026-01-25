package editpage

func (ep *EditPage) loadSnippet(snippetId int64) {
	ep.snippetId = snippetId
	snippet := ep.globalDeps.Repo.FindSnippet(ep.globalDeps.Ctx, snippetId)

	ep.body.SetText(snippet.Body, true)
	ep.title.SetText(snippet.Title, true)
	ep.description.SetText(snippet.Description, true)
	ep.urls.SetText(snippet.Url, true)
}
