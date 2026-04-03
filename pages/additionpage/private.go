package additionpage

import "github.com/gwkeit/configuration"

func (ap *AdditionPage) clearFields() {
	ap.body.SetText("", true)
	ap.title.SetText("", true)
	ap.description.SetText("", true)
	ap.urls.SetText("", true)
}

func (ap *AdditionPage) setLanguageOptionProgrammatically(index int) {
	previousValue := ap.isLangSelectFuncSuppressed.Load()
	ap.isLangSelectFuncSuppressed.Store(true)
	ap.language.SetCurrentOption(index)
	ap.isLangSelectFuncSuppressed.Store(previousValue)
}

func (ap *AdditionPage) updateCodePreview() {
	_, textOption := ap.language.GetCurrentOption()

	if textOption != "" {
		selectedLanguage := configuration.Language(textOption)
		ap.codePreview.SetText(ap.body.GetText(), selectedLanguage)
	}
}
