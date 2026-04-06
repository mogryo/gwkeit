package editpage

import (
	"slices"

	"github.com/gwkeit/configuration"
)

func (ep *EditPage) loadSnippet(snippetId int64) {
	ep.snippetId = snippetId
	snippet := ep.tools.Repo.FindSnippet(ep.tools.Ctx, snippetId)

	ep.body.SetText(snippet.Body, false)
	ep.title.SetText(snippet.Title, false)
	ep.description.SetText(snippet.Description, false)
	ep.urls.SetText(snippet.Url, false)
	if snippet.Language.Valid {
		index := slices.Index(configuration.LanguagesStrings, snippet.Language.String)
		if index < 0 {
			index = 0
		}
		ep.setLanguageOptionProgrammatically(index)
	} else {
		ep.setLanguageOptionProgrammatically(0)
	}

	ep.updateCodePreview()
}

func (ep *EditPage) setLanguageOptionProgrammatically(index int) {
	previousValue := ep.isLangSelectFuncSuppressed.Load()
	ep.isLangSelectFuncSuppressed.Store(true)
	ep.language.SetCurrentOption(index)
	ep.isLangSelectFuncSuppressed.Store(previousValue)
}

func (ep *EditPage) updateCodePreview() {
	_, textOption := ep.language.GetCurrentOption()

	if textOption != "" {
		selectedLanguage := configuration.Language(textOption)
		ep.codePreview.SetText(ep.body.GetText(), selectedLanguage)
	}
}
