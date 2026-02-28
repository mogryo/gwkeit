package additionpage

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
