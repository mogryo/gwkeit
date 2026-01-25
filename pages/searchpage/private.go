package searchpage

func (sp *SearchPage) showSnippet(snippetId int64) {
	snippet := sp.globalDeps.Repo.FindSnippet(sp.globalDeps.Ctx, snippetId)

	sp.title.SetText(snippet.Title, true)
	sp.body.SetText(snippet.Body, true)
	sp.description.SetText(snippet.Description, true)
	sp.urls.SetText(snippet.Url, true)
}
