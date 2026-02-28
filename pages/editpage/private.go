package editpage

import (
	"slices"

	"github.com/gwkeit/configuration"
)

func (ep *EditPage) loadSnippet(snippetId int64) {
	ep.snippetId = snippetId
	snippet := ep.tools.Repo.FindSnippet(ep.tools.Ctx, snippetId)

	ep.body.SetText(snippet.Body, true)
	ep.title.SetText(snippet.Title, true)
	ep.description.SetText(snippet.Description, true)
	ep.urls.SetText(snippet.Url, true)
	if snippet.Language.Valid {
		index := slices.Index(configuration.LanguagesStrings, snippet.Language.String)
		ep.setLanguageOptionProgrammatically(index)
	} else {
		ep.setLanguageOptionProgrammatically(0)
	}
}

func (ep *EditPage) setLanguageOptionProgrammatically(index int) {
	previousValue := ep.isLangSelectFuncSuppressed.Load()
	ep.isLangSelectFuncSuppressed.Store(true)
	ep.language.SetCurrentOption(index)
	ep.isLangSelectFuncSuppressed.Store(previousValue)
}
