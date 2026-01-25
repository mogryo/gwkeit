package additionpage

func (ap *AdditionPage) clearFields() {
	ap.body.SetText("", true)
	ap.title.SetText("", true)
	ap.description.SetText("", true)
	ap.urls.SetText("", true)
}
